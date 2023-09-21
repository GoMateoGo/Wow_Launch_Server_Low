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

	//5. 获取当前Server 的链接管理器
	GetConnMgr() IConnManager

	//注册[当创建连接后] Hook方法
	SetAfterStartConn(func(connection IConnection))
	//注册[当断开连接前] Hook方法
	SetBeforeStopConn(func(connection IConnection))
	//调用[当创建连接后] Hook方法
	CallAfterStartConn(connection IConnection)
	//调用[当断开连接前] Hook方法
	CallBeforeStopConn(connection IConnection)
}
