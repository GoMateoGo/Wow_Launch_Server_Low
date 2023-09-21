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

func (s *PingRouter) Handle(request wowiface.IRequest) {
	fmt.Println("Call Router Handle")
	//先读取客户端的数据,在回写ping.ping.ping

	fmt.Println("接受到的消息:", request.GetMsgId())
	fmt.Println("接受到的包体:", string(request.GetData()))

	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloRouter struct {
	wownet.BaseRouter
}

func (s *HelloRouter) Handle(request wowiface.IRequest) {
	fmt.Println("call router HelloWow")
	fmt.Println("接受到的消息Id:", request.GetMsgId())
	fmt.Println("接受到的包体内容:", string(request.GetData()))

	err := request.GetConnection().SendMsg(201, []byte("hello..Wow"))
	if err != nil {
		fmt.Println(err)
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
