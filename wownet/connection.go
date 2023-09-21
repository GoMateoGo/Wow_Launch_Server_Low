package wownet

import (
	"errors"
	"fmt"
	"gitee.com/mrmateoliu/wow_launch.git/utils"
	"gitee.com/mrmateoliu/wow_launch.git/wowiface"
	"io"
	"net"
	"sync"
)

/*
链接模块
*/
type Connection struct {
	//当前Conn隶属于哪个Server
	TcpServer wowiface.IServer

	// 当前链接的socket TCP套接字
	Conn *net.TCPConn

	// 链接的id
	ConnId uint32

	//当前的链接状态
	IsClosed bool

	//告知当前链接已经退出的/停止 channel(由Reader告知Write退出)
	ExitChan chan bool

	//无缓冲的管道,用于读/写Goroutine之间的消息通信
	msgChan chan []byte

	//消息的管理MsgId 和对应的处理业务Api关系
	MsgHandle wowiface.IMsgHandle

	//连接属性的集合
	property map[string]interface{}

	//保护连接属性的锁
	propertyLock sync.RWMutex
}

// 初始化链接模块的方法
func NewConnection(server wowiface.IServer, conn *net.TCPConn, connId uint32, msgHandle wowiface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer: server,
		Conn:      conn,
		ConnId:    connId,
		MsgHandle: msgHandle,
		IsClosed:  false,
		msgChan:   make(chan []byte),
		ExitChan:  make(chan bool, 1),
		property:  make(map[string]interface{}),
	}

	//将conn加入到ConnManager中
	c.TcpServer.GetConnMgr().Add(c)

	return c
}

// 链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("[读 Goroutine 运行中...]")
	defer fmt.Println("[读Goroutine退出] 当前链接id:", c.ConnId, "..,远端地址:", c.Conn.RemoteAddr().String())
	defer c.Stop()

	for {
		//读取客户端的数据到buf中,目前最大512字节
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("读取客户端数据失败", err)
		//	continue

		// 创建一个拆包解包的对象
		dp := NewDataPack()
		//读取客户端的Msg Head二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("读取 信息头 错误:", err)
			break
		}

		//拆包,得到MsgId 和 MsgDataLen 放在一个Msg对象中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("[拆包头]信息错误", err)
			break
		}

		//根据DataLen读取,再次读取包体Data, 放在Msg.Data属性中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("读取 [消息包体] 错误", err)
			}
		}

		//设置消息
		msg.SetData(data)

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		//如果已经开启工作池机制
		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经开启了工作池机制,将消息发送给worker工作池处理即可
			c.MsgHandle.SendMsgToTaskQueue(&req)
		} else {

			//从路由中 找到注册绑定的conn对应的router调用
			//根据绑定好的msgId 找到对应的处理Api业务 执行
			go c.MsgHandle.DoMsgHandler(&req)
		}
	}
}

// 写消息的Goroutine, 专门发送给客户端消息的方法
func (c *Connection) StartWriter() {
	fmt.Println("[写 Goroutine 运行中...]")
	defer fmt.Println("[写Goroutine退出] 当前链接id:", c.ConnId, "..,远端地址:", c.Conn.RemoteAddr().String())
	//阻塞,等待channel的消息,然后进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("发送消息数据出错:", err)
				return
			}
		case <-c.ExitChan:
			//代表Reader已退出,此时Write也要退出
			return
		}
	}
}

// 启动链接 让当前的链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("链接成功, 当前链接id:", c.ConnId)

	//启动从当前连接的读数据业务
	go c.StartReader()

	//启动从当前连接写数据的业务
	go c.StartWriter()

	//按照开发者传递进来的 创建链接之后需要调用的处理业务, 执行对应的Hook函数
	c.TcpServer.CallAfterStartConn(c)
}

// 停止链接 结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("当前链接已停止...停止id", c.ConnId)

	//如果当前链接已经关闭
	if c.IsClosed == true {
		return
	}
	c.IsClosed = true

	//调用开发者注册的关闭连接之前的Hook函数
	c.TcpServer.CallBeforeStopConn(c)

	//应该关闭socekt链接
	c.Conn.Close()

	//告知writer关闭
	c.ExitChan <- true

	//将当前链接从ConnMgr中摘除
	c.TcpServer.GetConnMgr().Remove(c)

	//关闭管道 回收资源
	close(c.ExitChan)
	close(c.msgChan)
}

// 获取当前链接绑定的socket connect
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前链接模块的链接id
func (c *Connection) GetConnId() uint32 {
	return c.ConnId
}

// 获取远程客户端的 tcp状态 ip和port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 提供一个SendMsg方法, 将要发送给客户端的数据进行封包,在发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.IsClosed == true {
		return errors.New("链接已关闭,无需发送消息")
	}
	//将data进行封包 |MsgDatalen|MsgId|MsgData|
	dp := NewDataPack()

	//进行包的组装(封包),最终格式为: |MsgDataLen|MsgId|Data|
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("数据包封装失败,产生错误:", err)
		return errors.New("数据包封装失败")
	}

	//将数据回写给客户端
	c.msgChan <- binaryMsg

	return nil
}

// 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	//设置属性
	c.property[key] = value
}

// 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	//查询属性
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New(fmt.Sprintf("当前链接没有这个属性:%s", key))
	}
}

// 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
