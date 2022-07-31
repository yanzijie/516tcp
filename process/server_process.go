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
}

func (s *ServerProcess) StatServer() {
	utils.Log.Info("[START] Server name: %s, listen at IP: %s, Port %d is starting",
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

		//3.阻塞等待客户端的链接
		for {
			conn, err := listen.AcceptTCP()
			if err != nil {
				utils.Log.Error(" AcceptTCP err: %s", err.Error())
				continue
			}

			//4.读写客户端请求
			//先搞个简单的
			go func() {
				for {
					buf := make([]byte, 512)
					// 把数据读到buf中, dataLen是读取到的数据长度
					dataLen, err := conn.Read(buf)
					if err != nil && err.Error() != "EOF" {
						utils.Log.Error(" read buf error: %s", err.Error())
						continue
					}
					//如果没收到数据，就不进行回写了
					if dataLen == 0 {
						continue
					}
					utils.Log.Info("receive client data: %s, len: %d", string(buf), dataLen)
					// 回写
					_, err = conn.Write([]byte("come on baby"))
					if err != nil {
						utils.Log.Error(" write buf error: %s", err.Error())
						continue
					}
				}
			}()
		}
	}()

}

func (s *ServerProcess) StopServer() {
	// 回收资源
}

func (s *ServerProcess) RunServer() {
	s.StatServer()

	// 做一些启动服务器之后的操作 TODO

	// 阻塞等待
	select {}
}

// NewServerProcess 初始化server_process模块
func NewServerProcess(name string) inface.ServerInterface {
	return &ServerProcess{
		ServerName: name,
		IPVersion:  "tcp4",
		IP:         "0.0.0.0",
		Port:       8999,
	}
}
