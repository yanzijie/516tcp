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
		utils.Log.Info("start tcp %s, ", s.ServerName, "success, begin listen....")

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
					// 把数据读到buf中, cnt是读取到的数据长度
					cnt, err := conn.Read(buf)
					if err != nil {
						utils.Log.Error(" read buf error: %s", err.Error())
						continue
					}

					// 回写
					_, err = conn.Write(buf[:cnt])
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

	// 做一些启动服务器之后的操作

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
