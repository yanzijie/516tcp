package main

import (
	"github.com/yanzijie/516tcp/inface"
	"github.com/yanzijie/516tcp/process"
	"github.com/yanzijie/516tcp/utils"
)

// TestRouter 自定义路由
type TestRouter struct {
	process.BaseRouter
}

// PreHandle 处理业务之前的方法（钩子）
//func (p *TestRouter) PreHandle(req inface.RequestInterface) {
//	utils.Log.Info(" router preHandle")
//}

// Handle 处理业务的主要方法
func (p *TestRouter) Handle(req inface.RequestInterface) {
	utils.Log.Info(" router Handle")
	//读取数据，已经读好了，在req里面
	utils.Log.Info(" receive msgId = %d, dataLen = %d, data = %s",
		req.GetMsgId(), len(req.GetMsgData()), string(req.GetMsgData()))

	//回复客户端
	err := req.GetConnection().Send(1, []byte("fuck!!fuck!!"))
	if err != nil {
		utils.Log.Error(" Send to client error:%s", err.Error())
	}
}

// PostHandle 处理业务之后的方法（钩子）
//func (p *TestRouter) PostHandle(req inface.RequestInterface) {
//	utils.Log.Info(" router PostHandle")
//}

func main() {
	utils.Log.Info("this is server....")
	//1.基于512tcp创建一个server
	s := process.NewServerProcess()
	//2.添加自定义路由
	s.AddRouter(&TestRouter{})
	//3.启动
	s.RunServer()
}
