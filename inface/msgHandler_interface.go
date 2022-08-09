package inface

// 消息路由管理模块

type MsgHandlerInterface interface {
	DoMsgHandler(req RequestInterface)              // 根据msgId分发路由
	AddRouter(msgId uint32, router RouterInterface) // 注册路由到对应的map表
}
