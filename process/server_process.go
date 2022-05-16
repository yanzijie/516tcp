package process

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

}

func (s *ServerProcess) StopServer() {

}

func (s *ServerProcess) RunServer() {

}
