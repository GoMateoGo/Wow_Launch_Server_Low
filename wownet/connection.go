package wownet

import (
	"fmt"
	"gitee.com/mrmateoliu/wow_launch.git/wowiface"
	"net"
)

/*
链接模块
*/
type Connection struct {
	// 当前链接的socket TCP套接字
	Conn *net.TCPConn

	// 链接的id
	ConnId uint32

	//当前的链接状态
	IsClosed bool

	//当前链接所绑定的处理业务方法api
	HandleApi wowiface.HandleFunc

	//告知当前链接已经退出的/停止 channel
	ExitChan chan bool
}

// 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connId uint32, callback_api wowiface.HandleFunc) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnId:    connId,
		HandleApi: callback_api,
		IsClosed:  false,
		ExitChan:  make(chan bool, 1),
	}

	return c
}

// 链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("读取 goroutine 数据中....")
	defer fmt.Println("当前链接id:", c.ConnId, "读写退出..,远端地址:", c.Conn.RemoteAddr().String())
	defer c.Stop()

	for {
		//读取客户端的数据到buf中,目前最大512字节
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("读取客户端数据失败", err)
			continue
		}

		//调用当前链接所绑定的handleAPI
		if err := c.HandleApi(c.Conn, buf, cnt); err != nil {
			fmt.Println("绑定HandleAPI失败,链接id", c.ConnId, "错误:", err)
			break
		}
	}
}

// 启动链接 让当前的链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("链接成功, 当前链接id:", c.ConnId)

	//启动从当前连接的读数据业务
	go c.StartReader()

	//TODO 启动从当前连接写数据的业务

}

// 停止链接 结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("当前链接已停止...停止id", c.ConnId)

	//如果当前链接已经关闭
	if c.IsClosed == true {
		return
	}
	c.IsClosed = true

	//应该关闭socekt链接
	c.Conn.Close()

	//关闭管道 回收资源
	close(c.ExitChan)
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

// 发送数据, 将数据发送给远程的客户端
func (c *Connection) Send(data []byte) error {
	return nil
}
