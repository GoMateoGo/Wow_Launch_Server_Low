package wowiface

type IServer interface {
	//1. 启动服务器
	Start()

	//2. 停止服务器
	Stop()

	//3. 运行服务器
	Server()

	//4.路由功能:给当前的服务注册一个路由方法,供客户端的链接处理使用
	AddRouter(msgId uint32, router IRouter)
}
