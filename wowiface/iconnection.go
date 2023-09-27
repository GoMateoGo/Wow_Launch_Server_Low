package wowiface

import "net"

// 定义链接模块的抽象层
type IConnection interface {
	// 启动链接 让当前的链接准备开始工作
	Start()

	// 停止链接 结束当前链接的工作
	Stop()

	// 获取当前链接绑定的socket connect
	GetTCPConnection() *net.TCPConn

	// 获取当前链接模块的链接id
	GetConnId() uint32

	//设置客户端Mac地址
	SetConnMac(mac string)

	//设置客户端系统
	SetConnOs(os string)

	//获取客户端Mac地址
	GetConnMac() string

	//获取客户端系统
	GetConnOs() string

	// 获取远程客户端的 tcp状态 ip和port
	RemoteAddr() net.Addr

	// 发送数据, 将数据发送给远程的客户端
	SendMsg(uint32, []byte) error

	//设置链接属性
	SetProperty(key string, value interface{})

	//获取链接属性
	GetProperty(key string) (interface{}, error)

	//移除链接属性
	RemoveProperty(key string)
}

// 定义一个处理链接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
