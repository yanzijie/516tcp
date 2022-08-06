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
func (p *TestRouter) PreHandle(req inface.RequestInterface) {
	utils.Log.Info(" router preHandle")
	// 回写客户端
	_, err := req.GetConnection().GetTCPSocketConn().Write([]byte("router preHandle...\n"))
	if err != nil {
		utils.Log.Error(" router preHandle error:%s", err.Error())
	}
}

// Handle 处理业务的主要方法
func (p *TestRouter) Handle(req inface.RequestInterface) {
	utils.Log.Info(" router Handle")
	// 回写客户端
	_, err := req.GetConnection().GetTCPSocketConn().Write([]byte("router Handle...\n"))
	if err != nil {
		utils.Log.Error(" router Handle error:%s", err.Error())
	}
}

// PostHandle 处理业务之后的方法（钩子）
func (p *TestRouter) PostHandle(req inface.RequestInterface) {
	utils.Log.Info(" router PostHandle")
	// 回写客户端
	_, err := req.GetConnection().GetTCPSocketConn().Write([]byte("router PostHandle...\n"))
	if err != nil {
		utils.Log.Error(" router PostHandle error:%s", err.Error())
	}
}

func main() {
	utils.Log.Info("this is server....")
	//1.基于512tcp创建一个server
	s := process.NewServerProcess("512tcp-v0.3")
	//2.添加自定义路由
	r := &TestRouter{}
	s.AddRouter(r)
	//3.启动
	s.RunServer()
}
