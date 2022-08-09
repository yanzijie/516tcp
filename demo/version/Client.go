package main

import (
	"github.com/yanzijie/516tcp/process"
	"github.com/yanzijie/516tcp/utils"
	"io"
	"net"
	"time"
)

// 模拟客户端
func main() {
	utils.Log.Info("this is client....")

	//1.创建链接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		utils.Log.Error("net dial error: %s", err.Error())
		return
	}

	// 1.创建链接，使用tcpConn
	//tcpAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8999")
	//conn, err := net.DialTCP("tcp", nil, tcpAddr)
	//if err != nil {
	//	utils.Log.Error(" net.DialTCP error: %s", err.Error())
	//	return
	//}
	defer conn.Close()

	for {
		// 2.封包, 发数据
		dp := process.NewDataPackProcess()
		binaryMsg, err := dp.Pack(process.NewMessageProcess(1, []byte("hello i am fuck")))
		if err != nil {
			utils.Log.Error(" package msg error: %s", err.Error())
			break
		}
		_, err = conn.Write(binaryMsg)
		if err != nil {
			utils.Log.Error("Write data error: %s", err.Error())
			break
		}

		// 客户端和服务端的conn不同，不能用服务端封装好的拆包方法
		//1.读取包头数据
		binaryHead := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, binaryHead)
		if err != nil {
			utils.Log.Error("io.ReadFull error: %s", err.Error())
			break
		}
		msg, err := dp.UnpackHead(binaryHead)
		if err != nil {
			utils.Log.Error("UnpackHead error: %s", err.Error())
			break
		}

		//2.读取包体内容
		if msg.GetMsgLen() > 0 {
			data := make([]byte, msg.GetMsgLen())
			_, err = io.ReadFull(conn, data)
			if err != nil {
				utils.Log.Error("io.ReadFull error: %s", err.Error())
				break
			}
			msg.SetMsgData(data)
		}

		utils.Log.Info("receive server data : %s, len : %d", msg.GetMsgData(), msg.GetMsgLen())

		time.Sleep(2 * time.Second)
	}

}
