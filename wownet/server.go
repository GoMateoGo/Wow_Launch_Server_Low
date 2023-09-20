package wownet

import (
	"errors"
	"fmt"
	"gitee.com/mrmateoliu/wow_launch.git/wowiface"
	"net"
)

type Server struct {
	// 服务器名称
	Name string
	// 服务器版本
	SVersion string
	//服务器IP
	Ip string
	//服务器端口
	Port int
}

// 定义当前客户端链接的所绑定的handle api(目前这个handle是写死的, 以后优化应该由业务端自定义handle方法)
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	//回显业务
	fmt.Println("处理回调客户端回显操作...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		if err != nil {
			fmt.Println("回写给客户端出现错误:", err)
			return errors.New("CallBackToClient Error")
		}
	}
	return nil
}

// 1.启动服务器
func (s *Server) Start() {
	fmt.Printf("服务器启动.. 地址:%s, 端口:%d\n", s.Ip, s.Port)

	go func() {
		// 1.获取一个TCP的addr
		addr, err := net.ResolveTCPAddr(s.SVersion, fmt.Sprintf("%s:%d", s.Ip, s.Port))
		if err != nil {
			//TODO 增加日志
			fmt.Println("服务器启动错误:", err)
			return
		}
		// 2.监听服务器地址
		listener, err := net.ListenTCP(s.SVersion, addr)
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
			dealConn := NewConnection(conn, cid, CallBackToClient)
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

// 初始化Server模块方法
func NewServer(name string) wowiface.IServer {
	s := &Server{
		Name:     name,
		SVersion: "tcp4",
		Ip:       "0.0.0.0",
		Port:     8999,
	}

	return s
}
