package main

import "github.com/yanzijie/516tcp/process"

func main() {
	//1.基于512tcp创建一个server
	s := process.NewServerProcess("512tcp-v0.1")
	//2.启动
	s.RunServer()
}
