package wownet

import (
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
		listenner, err := net.ListenTCP(s.SVersion, addr)
		if err != nil {
			//TODO 日志
			fmt.Println("监听服务器地址错误:", err)
			return
		}
		fmt.Printf("服务器启动成功: 服务器名:%s\n", s.Name)

		// 3.阻塞等待客户端链接,处理客户端链接业务(读/写)
		for {
			//如果客户端链接过来会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				//TODO 日志
				fmt.Printf("链接客户端错误:%s\n", err)
				continue
			}

			// 客户端已经建立成功
			go func() {
				for {
					//创建一个512字节切片
					buf := make([]byte, 512)
					//读取客户端发送过来的数据
					cnt, err := conn.Read(buf)
					if err != nil {
						//TODO 日志
						fmt.Println("读取客户端数据错误:", err)
						continue
					}

					fmt.Println("接收客户端数据:", string(buf))

					//回显
					if _, err = conn.Write(buf[:cnt]); err != nil {
						fmt.Println("回显出问题了..", err)
					}
				}
			}()
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
