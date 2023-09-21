package wownet

import (
	"fmt"
	"gitee.com/mrmateoliu/wow_launch.git/utils"
	"gitee.com/mrmateoliu/wow_launch.git/wowiface"
)

/*
	消息处理模块的实现
*/

type MsgHandle struct {
	//存放每个MsgId 所对应的处理方法
	APis map[uint32]wowiface.IRouter

	//负责Worker取任务的 [消息队列]
	TaskQueue []chan wowiface.IRequest

	//业务 [工作池数量]
	WorkerPoolSize uint32
}

// 初始化/创建MsgHandle方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		APis:           make(map[uint32]wowiface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, //从全局配置文件中读取
		TaskQueue:      make([]chan wowiface.IRequest, utils.GlobalObject.MaxWorkerTaskLen),
	}
}

// 调度/执行 对应的 Router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request wowiface.IRequest) {
	//1. 从Request中找到msgId
	handle, ok := mh.APis[request.GetMsgId()]
	if !ok {
		fmt.Println("没有在Api容器中找到对应的[注册消息Id]:", request.GetMsgId(), "请注册后再试")
		return
	}
	//2. 根据MsgId 调度对应的router业务即可
	handle.BeforeHandle(request)
	handle.Handle(request)
	handle.AfterHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgId uint32, router wowiface.IRouter) {
	// 1. 判断当前的MsgId绑定的Api处理方法是否存在
	if _, ok := mh.APis[msgId]; ok {
		//id已经注册,存在的
		panic(fmt.Sprintf("重复的消息路由Api已注册:,消息Id:%d", msgId))
	}
	// 2. 添加Msg与API的绑定管理
	mh.APis[msgId] = router
	fmt.Println("添加APi消息Id:", msgId, "成功!")
}
