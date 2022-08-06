package inface

type ServerInterface interface {
	StatServer()                      // StatServer 启动服务器
	StopServer()                      // StopServer 停止服务器
	RunServer()                       // RunServer 运行服务器
	AddRouter(router RouterInterface) // AddRouter 给当前服务注册路由方法，用来处理客户端的链接
}
