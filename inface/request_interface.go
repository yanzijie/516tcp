package inface

// 封装客户端请求链接和发过来的请求数据, 封装到一个request中

type RequestInterface interface {
	GetConnection() ConnectionInterface // 得到当前链接
	GetReqData() []byte                 // 获取请求的消息数据
}
