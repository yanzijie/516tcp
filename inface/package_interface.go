package inface

import "net"

// 消息封装模块，处理tcp包的粘包和拆包

type DataPackInterface interface {
	// GetHeadLen 获取包头的长度，返回：包头长度
	GetHeadLen() uint32
	// Pack 把一个message封装成一个二进制流的包, 返回: 封装完成的包，处理error
	Pack(msg MessageInterface) ([]byte, error)
	// UnpackHead 把读取到的二进制流的包(这个包是数据的包头)拆成一个message, 返回: 拆出来的msg，处理error
	UnpackHead([]byte) (MessageInterface, error)
	// Unpack 完全的拆包，把所有数据都拆完, 只能服务端使用，客户端的conn不同，不能用
	Unpack(conn *net.TCPConn) (MessageInterface, error)
}
