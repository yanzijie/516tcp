package main

import (
	"github.com/yanzijie/516tcp/inface"
	"github.com/yanzijie/516tcp/process"
	"github.com/yanzijie/516tcp/utils"
)

// OneRouter 自定义路由
type OneRouter struct {
	process.BaseRouter
}

type TwoRouter struct {
	process.BaseRouter
}

var msgIdOne uint32 = 1
var msgIdTwo uint32 = 2

// PreHandle 处理业务之前的方法（钩子）
//func (p *OneRouter) PreHandle(req inface.RequestInterface) {
//	utils.Log.Info(" router preHandle")
//}

// PostHandle 处理业务之后的方法（钩子）
//func (p *OneRouter) PostHandle(req inface.RequestInterface) {
//	utils.Log.Info(" router PostHandle")
//}

// Handle 处理业务的主要方法
func (p *OneRouter) Handle(req inface.RequestInterface) {
	utils.Log.Info(" router Handle")
	//读取数据，已经读好了，在req里面
	utils.Log.Info(" receive msgId = %d, dataLen = %d, data = %s",
		req.GetMsgId(), len(req.GetMsgData()), string(req.GetMsgData()))

	//回复客户端
	err := req.GetConnection().Send(msgIdOne, []byte("fuck!!fuck!! this is one"))
	if err != nil {
		utils.Log.Error(" Send to client error:%s", err.Error())
	}
}

func (p *TwoRouter) Handle(req inface.RequestInterface) {
	utils.Log.Info(" router Handle")
	//读取数据，已经读好了，在req里面
	utils.Log.Info(" receive msgId = %d, dataLen = %d, data = %s",
		req.GetMsgId(), len(req.GetMsgData()), string(req.GetMsgData()))

	//回复客户端
	err := req.GetConnection().Send(msgIdTwo, []byte("fuck!!fuck!! this is two"))
	if err != nil {
		utils.Log.Error(" Send to client error:%s", err.Error())
	}
}

func main() {
	utils.Log.Info("this is server....")
	//1.基于512tcp创建一个server
	s := process.NewServerProcess()
	//2.添加自定义路由
	s.AddRouter(msgIdOne, &OneRouter{})
	s.AddRouter(msgIdTwo, &TwoRouter{})
	//3.启动
	s.RunServer()
}
