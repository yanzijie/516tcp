package main

import (
	"fmt"
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

type ThreeRouter struct {
	process.BaseRouter
}

var msgIdOne uint32 = 1
var msgIdTwo uint32 = 2
var msgIdThree uint32 = 3

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

func (p *ThreeRouter) Handle(req inface.RequestInterface) {
	utils.Log.Info(" router Handle")
	//读取数据，已经读好了，在req里面
	utils.Log.Info(" receive msgId = %d, dataLen = %d, data = %s",
		req.GetMsgId(), len(req.GetMsgData()), string(req.GetMsgData()))

	//回复客户端
	err := req.GetConnection().Send(msgIdTwo, []byte("fuck!!fuck!! this is three"))
	if err != nil {
		utils.Log.Error(" Send to client error:%s", err.Error())
	}
}

// 链接创建之后执行
func connStart(conn inface.ConnectionInterface) {
	fmt.Println(" ===> conn start")
	err := conn.Send(202, []byte("conn start"))
	if err != nil {
		fmt.Println("Send error : ", err.Error())
	}

	// 给链接设置一些信息
	conn.SetProperty("name", "lao6")
	conn.SetProperty("nickName", "6666")
	conn.SetProperty("age", 10001)
}

// 链接断开之前执行
func connStop(conn inface.ConnectionInterface) {
	fmt.Println(" ===> connID = ", conn.GetConnID(), " is stop")

	if name, err := conn.GetProperty("name"); err == nil {
		fmt.Println(" name is ==> ", name)
	}

	if nickName, err := conn.GetProperty("nickName"); err == nil {
		fmt.Println(" nickName is ==> ", nickName)
	}

	if age, err := conn.GetProperty("age"); err == nil {
		fmt.Println(" age is ==> ", age)
	}
}

func main() {
	utils.Log.Info("this is server....")
	// 基于512tcp创建一个server
	s := process.NewServerProcess()

	// 注册钩子函数
	s.SetOnConnStart(connStart)
	s.SetOnConnStop(connStop)

	// 添加自定义路由
	s.AddRouter(msgIdOne, &OneRouter{})
	s.AddRouter(msgIdTwo, &TwoRouter{})
	s.AddRouter(msgIdThree, &ThreeRouter{})
	// 启动
	s.RunServer()
}
