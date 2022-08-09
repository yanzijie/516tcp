package process

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/yanzijie/516tcp/inface"
	"github.com/yanzijie/516tcp/utils"
	"io"
	"net"
)

// 消息封装模块，处理tcp包的粘包和拆包
// 包结构:
//   包头：4字节包长 + 4字节包的序列号id
//	 包体：msg的数据

type DataPackProcess struct {
}

func NewDataPackProcess() *DataPackProcess {
	return &DataPackProcess{}
}

func (d *DataPackProcess) GetHeadLen() uint32 {
	// 4字节包的长度 + 4字节包的序列号id
	return 8
}

// Pack 把一个message封装成一个二进制流的包, 返回: 封装完成的包，处理error
func (d *DataPackProcess) Pack(msg inface.MessageInterface) ([]byte, error) {
	// 创建buf的缓冲区
	dataBuff := bytes.NewBuffer([]byte{})

	// 写dataBuff, 顺序不能, 用小端对齐写，也有小端对齐拆,
	// 这里怎么知道要写多少字节？？
	// 因为
	//	msg.GetMsgLen() 和 msg.GetMsgId() 得到的是 uint32类型数据
	//	而uint32类型占4个字节, 所以写入的 dataLen 和 序列号id 都是4个字节，加起来就是8个字节
	// 1.写dataLen
	err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen())
	if err != nil {
		return nil, err
	}
	// 2.写msg序列号id
	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId())
	if err != nil {
		return nil, err
	}
	// 3.写msg数据
	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgData())
	if err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// UnpackHead 拆包, 这里实际上知识进行了包头的拆包.....
// 先读取固定长度的包头数据，得到消息的长度和序列号id，这个时候位置已经偏移到真正的数据的起始位置。
// 然后根据包头的长度，再去对应的conn链接里读取数据。
// binaryData : 数据的包头 4字节包长 + 4字节包的序列号id
func (d *DataPackProcess) UnpackHead(binaryData []byte) (inface.MessageInterface, error) {
	// 把二进制数据读到缓冲区
	dataBuff := bytes.NewReader(binaryData)
	msg := &MessageProcess{}

	// 读msg的长度
	err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen)
	if err != nil {
		return nil, err
	}

	// 判断包是否超过规定的大小
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("package too big")
	}

	// 读msg的id
	err = binary.Read(dataBuff, binary.LittleEndian, &msg.Id)
	if err != nil {
		return nil, err
	}

	// 下次直接读取数据就可以了，这里已经拿到了包的长度和id
	// ? 这里怎么拿到的, read函数怎么知道我要拿多少字节？？？
	// 因为：
	//	msg.DataLen 和 msg.Id 都是 uint32类型, 而 uint32类型占4个字节
	//	而 binaryData 是8个字节的包头数据，所以刚好可以读取到 msg.DataLen 和 msg.Id里面
	return msg, nil
}

func (d *DataPackProcess) Unpack(conn *net.TCPConn) (inface.MessageInterface, error) {
	// 拆包
	dp := NewDataPackProcess()
	// 1.第一次读，把包的head的二进制流读出来
	headData := make([]byte, dp.GetHeadLen())
	// 从conn中一次把headData读满,
	// 这个时候就读取到了8个字节（偏移到了真实数据的开始处），正好是全部的包头（前提是发包方必须用我们自定义的协议封包）
	_, err := io.ReadFull(conn, headData)
	if err != nil {
		utils.Log.Error(" io.ReadFull package head error: %s", err.Error())
		return nil, err
	}
	// 拆包,这个时候headData有包头8个字节的数据
	msg, err := dp.UnpackHead(headData)
	if err != nil {
		utils.Log.Error(" UnpackHead error: %s", err.Error())
		return nil, err
	}
	// 2.第二次读取，根据包长度，设置data长度，然后读取数据
	if msg.GetMsgLen() > 0 {
		data := make([]byte, msg.GetMsgLen())
		_, err = io.ReadFull(conn, data)
		if err != nil {
			utils.Log.Error(" io.ReadFull package data error: %s", err.Error())
			return nil, err
		}
		msg.SetMsgData(data)
	}

	return msg, nil
}
