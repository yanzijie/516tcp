package process

import (
	"github.com/yanzijie/516tcp/utils"
	"io"
	"net"
	"testing"
)

// 拆包封包功能的单元测试

func TestDataPack(t *testing.T) {
	/*******************************模拟服务器start*******************************/
	// 1.创建socketTcp
	listener, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		utils.Log.Error(" server listen error: %s", err.Error())
		return
	}

	// 2.从客户端读取数据
	go func() {
		for {
			// 阻塞监听
			conn, err := listener.Accept()
			if err != nil {
				utils.Log.Error(" listener Accept error: %s", err.Error())
				continue
			}

			// 处理客户端请求
			go func() {
				/***********拆包开始************/
				dp := NewDataPackProcess()
				for {
					// 1.第一次读，把包的head的二进制流读出来
					headData := make([]byte, dp.GetHeadLen())
					// 从conn中一次把headData读满,
					//这个时候就读取到了8个字节，正好是全部的包头（前提是发包方必须用我们自定义的协议封包）
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						utils.Log.Error(" io.ReadFull package head error: %s", err.Error())
						return
					}
					// 这个时候headData有包头8个字节的数据
					msgHead, err := dp.UnpackHead(headData)
					if err != nil {
						utils.Log.Error(" UnpackHead error: %s", err.Error())
						return
					}
					if msgHead.GetMsgLen() > 0 {
						// 2.第二次读，根据包长读数据
						msg := msgHead.(*MessageProcess)
						// 根据包长设置data长度
						msg.Data = make([]byte, msg.GetMsgLen())
						_, err = io.ReadFull(conn, msg.Data)
						if err != nil {
							utils.Log.Error(" io.ReadFull package data error: %s", err.Error())
							return
						}

						utils.Log.Info(" receive msgId = %d, dataLen = %d, data = %s",
							msg.Id, msg.DataLen, string(msg.Data))
					}

				}

				/***********拆包结束************/
			}()
		}
	}()
	/*******************************模拟服务器end*******************************/

	/*******************************模拟客户端start*******************************/
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		utils.Log.Error("Dial err: %s", err.Error())
		return
	}

	// 封包
	dp := NewDataPackProcess()
	//模拟粘包过程, 封装两个msg包一次性发送，如果服务端都能拆包，说明粘包问题已经解决
	msgOne := &MessageProcess{
		Id:      1,
		DataLen: 5,
		Data:    []byte{'f', 'u', 'n', 'k', '!'},
	}
	dataOne, err := dp.Pack(msgOne)
	if err != nil {
		utils.Log.Error("dp.Pack msgOne err: %s", err.Error())
		return
	}

	msgTwo := &MessageProcess{
		Id:      2,
		DataLen: 7,
		Data:    []byte{'c', 'o', 'm', 'e', ' ', 'o', 'n'},
	}
	dataTwo, err := dp.Pack(msgTwo)
	if err != nil {
		utils.Log.Error("dp.Pack msgTwo err: %s", err.Error())
		return
	}

	// 把两个包合拼接成一个 []byte包
	dataOne = append(dataOne, dataTwo...)
	// 发包
	_, err = conn.Write(dataOne)
	if err != nil {
		utils.Log.Error("conn.Write err: %s", err.Error())
		return
	}

	// 阻塞等待服务端返回
	select {}
	/*******************************模拟客户端end*******************************/
}
