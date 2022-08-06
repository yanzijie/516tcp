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
	IPVersion string
	// 服务器监听的ip
	IP string
	// 服务器监听的端口
	Port int
	// server链接对应的处理业务
	Router inface.RouterInterface
}

// CallBackToClient 当前客户端链接绑定的处理函数, 交由router去处理, 不在这里处理
//func CallBackToClient(conn *net.TCPConn, data []byte, dataLen int) (err error) {
//
//	// 回写数据
//	utils.Log.Info("write msg to client...,msg is: %s", string(data))
//	_, err = conn.Write(data[:dataLen])
//	if err != nil {
//		utils.Log.Error(" write buf error: %ss", err.Error())
//		return errors.New("write msg to client error")
//	}
//	return
//}

func (s *ServerProcess) StatServer() {
	utils.Log.Info("[START] Server name: %s, listen at IP: %s, Port %d is starting\n",
		s.ServerName, s.IP, s.Port)

	go func() {
		//1.获取一个TCP的地址
		ipPort := fmt.Sprintf("%s:%d", s.IP, s.Port)
		addr, err := net.ResolveTCPAddr(s.IPVersion, ipPort)
		if err != nil {
			utils.Log.Error("net.ResolveTCPAddr error: %s", err.Error())
			return
		}
		//2.监听服务器
		listen, err := net.ListenTCP(s.IPVersion, addr)
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
			processConn := NewConnection(conn, connId, s.Router)
			connId++

			// 开启链接的业务处理
			go processConn.StartRead()
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

func (s *ServerProcess) AddRouter(router inface.RouterInterface) {
	s.Router = router
	utils.Log.Info(" add router success")
}

// NewServerProcess 初始化server_process模块
func NewServerProcess(name string) inface.ServerInterface {
	return &ServerProcess{
		ServerName: name,
		IPVersion:  "tcp4",
		IP:         "0.0.0.0",
		Port:       8999,
		Router:     nil, // new的时候指定空, AddRouter的时候赋值
	}
}
