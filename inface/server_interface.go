package inface

type ServerInterface interface {
	StatServer()                                    // StatServer 启动服务器
	StopServer()                                    // StopServer 停止服务器
	RunServer()                                     // RunServer 运行服务器
	AddRouter(msgID uint32, router RouterInterface) // AddRouter 给当前服务注册路由方法，用来处理客户端的链接
	GetConnManager() ConnManagerInterface           // GetConnManager 获取链接管理模块接口
	SetOnConnStart(func(conn ConnectionInterface))  //注册OnConnStart
	SetOnConnStop(func(conn ConnectionInterface))   //注册OnConnStop
	CallOnConnStart(conn ConnectionInterface)       //调用OnConnStart
	CallOnConnStop(conn ConnectionInterface)        //调用OnConnStop
}
