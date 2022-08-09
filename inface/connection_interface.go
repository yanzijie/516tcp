package inface

import "net"

// 链接模块

type ConnectionInterface interface {
	StartConnection()                     // 启动链接
	StopConnection()                      // 停止链接
	GetTCPSocketConn() *net.TCPConn       // 获取当前链接绑定的socket conn
	GetConnID() uint32                    // 获取当前链接的链接ID
	GetClientAddr() net.Addr              // 获取客户端的tcp状态，ip和port
	Send(msgId uint32, data []byte) error // 发送数据, 先封包，再发送
	StartRead()                           // 读协程
	StartWrite()                          // 回写客户端协程
}

// HandleFunc 处理链接的业务方法
// conn-客户端链接，data-要处理的数据内容 ，dataLen-数据长度
type HandleFunc func(conn *net.TCPConn, data []byte, dataLen int) error
