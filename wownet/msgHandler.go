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
		TaskQueue:      make([]chan wowiface.IRequest, utils.GlobalObject.WorkerPoolSize),
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
	if utils.GlobalObject.Develop {
		fmt.Println("添加APi消息Id:", msgId, "成功!")
	}
}

// 启动一个Worker工作池 (开启工作池动作只能发生一次,最多只能有一个工作池)
func (mh *MsgHandle) StartWorkerPool() {
	//根据workerPoolSize 分别开启Worker, 每个Worker用一个Go承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		// 1.给当前的worker对应的channel消息队列 开辟空间, 第0个worker就用第0个channel
		mh.TaskQueue[i] = make(chan wowiface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 2.启动当前的worker, 阻塞等待消息从channel中传递进来
		go mh.startOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个Worker工作流程
func (mh *MsgHandle) startOneWorker(workerId int, taskQueue chan wowiface.IRequest) {
	if utils.GlobalObject.Develop {
		fmt.Println("当前工作Id:", workerId, "已启动.")
	}
	//不断的阻塞等待对应的队列消息
	for {
		select {
		//如果有消息过来,出列的是一个客户端的Request,执行当前Request所绑定的业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// 将消息交给TaskQueue, 由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request wowiface.IRequest) {
	// 1. 将消息平均分配给不同的worker
	//根据客户端建立的ConnId进行分配
	//----基本的平均分配的轮询法则
	workerId := request.GetConnection().GetConnId() % mh.WorkerPoolSize
	if utils.GlobalObject.Develop {
		fmt.Println("当前链接Id:", request.GetConnection().GetConnId(),
			"消息请求Id:", request.GetMsgId(),
			"所在工作池Id", workerId)
	}
	//2. 将消息发送给对应的worker的taskQueue即可.
	mh.TaskQueue[workerId] <- request
}
