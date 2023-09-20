package wowiface

/*
路由抽象接口
路由里的数据都是IRequest处理
*/
type IRouter interface {
	//处理conn业务之前的Hook方法
	BeforeHandle(request IRequest)
	//处理conn业务的主Hook方法
	Handle(request IRequest)
	//处理conn业务之后的Hook方法
	AfterHandle(request IRequest)
}
