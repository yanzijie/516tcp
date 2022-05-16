package inface

type ServerInface interface {
	// 启动服务器
	StatServer()
	// 停止服务器
	StopServer()
	// 运行服务器
	RunServer()
}
