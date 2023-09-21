package main

import (
	"fmt"
	"gitee.com/mrmateoliu/wow_launch.git/wowiface"
	"gitee.com/mrmateoliu/wow_launch.git/wownet"
)

// ping测试, 自定义路由
type PingRouter struct {
	wownet.BaseRouter
}

// Test beforeHandle
func (s *PingRouter) BeforeHandle(request wowiface.IRequest) {
	fmt.Println("Call Router BeforeHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("befor ping\n"))
	if err != nil {
		fmt.Println("call back BeforePing err", err)
	}
}

// Test Handle
func (s *PingRouter) Handle(request wowiface.IRequest) {
	fmt.Println("Call Router Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping ping ping\n"))
	if err != nil {
		fmt.Println("call back Ping Ping Ping err", err)
	}
}

// Test AfterHandle
func (s *PingRouter) AfterHandle(request wowiface.IRequest) {
	fmt.Println("Call Router AfterHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping\n"))
	if err != nil {
		fmt.Println("call back AfterHandle err", err)
	}
}

func main() {
	// 1. 创建server句柄
	s := wownet.NewServer("wow-launch")
	// 2. 给当前框架添加自定义router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	// 3. 启动服务器
	s.Server()
}
