package wowiface

/*
	消息管理抽象层
*/

type IMsgHandle interface {
	//调度/执行 对应的 Router消息处理方法
	DoMsgHandler(request IRequest)
	//为消息添加具体的处理逻辑
	AddRouter(msgId uint32, router IRouter)
}
