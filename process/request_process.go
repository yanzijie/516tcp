package process

import "github.com/yanzijie/516tcp/inface"

type RequestProcess struct {
	Conn    inface.ConnectionInterface // 建立完成的链接句柄Connection
	ReqData inface.MessageInterface    // 客户端的请求数据
}

// GetConnection 得到当前链接
func (r *RequestProcess) GetConnection() inface.ConnectionInterface {
	return r.Conn
}

// GetMsgData 获取请求的消息数据
func (r *RequestProcess) GetMsgData() []byte {
	return r.ReqData.GetMsgData()
}

// GetMsgId 获取消息序列号id
func (r *RequestProcess) GetMsgId() uint32 {
	return r.ReqData.GetMsgId()
}
