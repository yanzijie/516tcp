package inface

// 消息封装模块

type MessageInterface interface {
	GetMsgId() uint32   // 得到消息id
	GetMsgLen() uint32  // 得到消息长度
	GetMsgData() []byte // 得到消息内容

	SetMsgId(id uint32)     // 设置消息id
	SetMsgLen(len uint32)   // 设置消息长度
	SetMsgData(data []byte) // 设置消息内容
}
