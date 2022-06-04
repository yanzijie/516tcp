package inface

type ServerInterface interface {
	StatServer() // StatServer 启动服务器
	StopServer() // StopServer 停止服务器
	RunServer()  // RunServer 运行服务器
}
