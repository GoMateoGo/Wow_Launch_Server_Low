package wowiface

type IServer interface {
	//1. 启动服务器
	Start()

	//2. 停止服务器
	Stop()

	//3. 运行服务器
	Server()
}
