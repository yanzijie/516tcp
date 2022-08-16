package process

import (
	"fmt"
	"github.com/yanzijie/516tcp/inface"
	"github.com/yanzijie/516tcp/utils"
	"net"
)

// ServerProcess 这个结构去实现server_interface的接口
type ServerProcess struct {
	// 服务器名称
	ServerName string
	// 服务器绑定的ip版本
	Version string
	// 服务器监听的ip
	IP string
	// 服务器监听的端口
	Port int
	// server链接对应的处理业务, 绑定msgID和对应的handler
	MsgHandler inface.MsgHandlerInterface
}

func (s *ServerProcess) StatServer() {
	utils.Log.Info("[START] Server name: %s, listen at IP: %s, Port %d, version: %s is starting\n",
		s.ServerName, s.IP, s.Port, s.Version)

	go func() {
		// 开启消息队列和工作协程池
		s.MsgHandler.StartWorkerPool()
		//1.获取一个TCP的地址
		ipPort := fmt.Sprintf("%s:%d", s.IP, s.Port)
		addr, err := net.ResolveTCPAddr(s.Version, ipPort)
		if err != nil {
			utils.Log.Error("net.ResolveTCPAddr error: %s", err.Error())
			return
		}
		//2.监听服务器
		listen, err := net.ListenTCP(s.Version, addr)
		if err != nil {
			utils.Log.Error("listen fail: ", err)
			return
		}
		utils.Log.Info("start tcp %s, success, begin listen....", s.ServerName)
		var connId uint32
		connId = 0

		//3.阻塞等待客户端的链接
		for {
			conn, err := listen.AcceptTCP()
			if err != nil {
				utils.Log.Error(" AcceptTCP err: %s", err.Error())
				continue
			}

			// 绑定客户端链接和该链接的业务处理方法
			processConn := NewConnection(conn, connId, s.MsgHandler)
			connId++

			// 开启链接的业务处理
			go processConn.StartConnection()
		}
	}()

}

func (s *ServerProcess) StopServer() {
	// 回收资源
}

func (s *ServerProcess) RunServer() {
	// 在StatServer之前需要先把路由添加上
	s.StatServer()

	// 做一些启动服务器之后的操作 TODO

	// 阻塞等待
	select {}
}

func (s *ServerProcess) AddRouter(msgID uint32, router inface.RouterInterface) {
	s.MsgHandler.AddRouter(msgID, router)
}

// NewServerProcess 初始化server_process模块
func NewServerProcess() inface.ServerInterface {
	return &ServerProcess{
		ServerName: utils.GlobalObject.Name,
		Version:    "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandlerProcess(), // new的时候指定空, AddRouter的时候赋值
	}
}
