package main

import (
	"github.com/yanzijie/516tcp/process"
	"github.com/yanzijie/516tcp/utils"
)

func main() {
	utils.Log.Info("this is server....")
	//1.基于512tcp创建一个server
	s := process.NewServerProcess("512tcp-v0.1")
	//2.启动
	s.RunServer()
}
