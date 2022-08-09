package process

import (
	"github.com/yanzijie/516tcp/inface"
	"github.com/yanzijie/516tcp/utils"
)

type MsgHandlerProcess struct {
	// msgID与router路由对应关系map
	ApiMap map[uint32]inface.RouterInterface
}

func NewMsgHandlerProcess() *MsgHandlerProcess {
	return &MsgHandlerProcess{
		ApiMap: make(map[uint32]inface.RouterInterface),
	}
}

func (m *MsgHandlerProcess) DoMsgHandler(req inface.RequestInterface) {
	//1. 判断msgId
	handler, ok := m.ApiMap[req.GetMsgId()]
	if !ok {
		utils.Log.Error(" this msgId = %d is not register", req.GetMsgId())
		return
	}
	//2. 根据msgId调用对应的router
	handler.PreHandle(req)
	handler.Handle(req)
	handler.PostHandle(req)
}

func (m *MsgHandlerProcess) AddRouter(msgId uint32, router inface.RouterInterface) {
	// 校验msgID
	_, ok := m.ApiMap[msgId]
	if ok {
		// msgID已经注册
		utils.Log.Error(" this msgId = %d is register", msgId)
		return
	}
	m.ApiMap[msgId] = router
	utils.Log.Info(" msgID = %d register success!", msgId)
}
