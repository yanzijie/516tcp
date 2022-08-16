package inface

// 消息路由管理模块

type MsgHandlerInterface interface {
	DoMsgHandler(req RequestInterface)                            // 根据msgId分发路由
	AddRouter(msgId uint32, router RouterInterface)               // 注册路由到对应的map表
	StartWorkerPool()                                             // 启动协程池(只能开启一次)
	StartOneWorker(workerID int, taskQueen chan RequestInterface) // 启动协程工作流程
	SendMsgToTaskQueen(request RequestInterface)                  // 接受请求，发送到消息队列里面
}
