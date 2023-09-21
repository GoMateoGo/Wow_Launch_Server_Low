package wownet

import (
	"fmt"
	"gitee.com/mrmateoliu/wow_launch.git/utils"
	"gitee.com/mrmateoliu/wow_launch.git/wowiface"
	"net"
)

type Server struct {
	// 服务器名称
	Name string
	// 服务器版本
	IpVersion string
	//服务器IP
	Ip string
	//服务器端口
	Port int
	//当前server的消息管理模块, 用来绑定MsgId和对应的处理业务Api关系
	MsgHandler wowiface.IMsgHandle
}

// 1.启动服务器
func (s *Server) Start() {
	fmt.Printf("[配置信息]:\n 服务器名称:%s\n Ip地址:%s\n 端口号:%d\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf(" 版本:%s\n 最大链接数量:%d\n 最大包尺寸:%d\n", utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)
	fmt.Printf("服务器启动.. 地址:%s, 端口:%d\n", s.Ip, s.Port)

	go func() {
		// 0. 开启消息队列及Worker工作池
		s.MsgHandler.StartWorkerPool()

		// 1.获取一个TCP的addr
		addr, err := net.ResolveTCPAddr(s.IpVersion, fmt.Sprintf("%s:%d", s.Ip, s.Port))
		if err != nil {
			//TODO 增加日志
			fmt.Println("服务器启动错误:", err)
			return
		}
		// 2.监听服务器地址
		listener, err := net.ListenTCP(s.IpVersion, addr)
		if err != nil {
			//TODO 日志
			fmt.Println("监听服务器地址错误:", err)
			return
		}
		fmt.Printf("服务器启动成功: 服务器名:%s\n", s.Name)
		var cid uint32
		cid = 0

		// 3.阻塞等待客户端链接,处理客户端链接业务(读/写)
		for {
			//如果客户端链接过来会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				//TODO 日志
				fmt.Printf("链接客户端错误:%s\n", err)
				continue
			}

			// 将该处理新链接的业务方法 和 conn机型绑定,得到链接模块
			dealConn := NewConnection(conn, cid, s.MsgHandler)
			cid++

			//启动当前的连接业务处理
			go dealConn.Start()
		}
	}()
}

// 2.停止服务器
func (s *Server) Stop() {
	//TODO 释放资源等等.. 状态,资源,开辟的链接.
}

// 3.运行服务器
func (s *Server) Server() {
	s.Start()

	//TODO 其他业务...
	//需要阻塞状态
	select {}
}

// 4.添加一个路由方法
func (s *Server) AddRouter(msgId uint32, router wowiface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
	fmt.Println("添加 Router 成功")
}

// 初始化Server模块方法
func NewServer(name string) wowiface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IpVersion:  "tcp4",
		Ip:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
	}

	return s
}
