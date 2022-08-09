package process

type MessageProcess struct {
	Id      uint32 // 消息id
	DataLen uint32 // 消息长度
	Data    []byte // 消息内容
}

// NewMessageProcess 创建一个消息包
func NewMessageProcess(msgId uint32, data []byte) *MessageProcess {
	return &MessageProcess{
		Id:      msgId,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

func (m *MessageProcess) GetMsgId() uint32 {
	return m.Id
}

func (m *MessageProcess) GetMsgLen() uint32 {
	return m.DataLen
}

func (m *MessageProcess) GetMsgData() []byte {
	return m.Data
}

func (m *MessageProcess) SetMsgId(id uint32) {
	m.Id = id
}

func (m *MessageProcess) SetMsgLen(len uint32) {
	m.DataLen = len
}

func (m *MessageProcess) SetMsgData(data []byte) {
	m.Data = data
}
