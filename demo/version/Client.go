package main

import (
	"github.com/yanzijie/516tcp/process"
	"github.com/yanzijie/516tcp/utils"
	"net"
	"time"
)

// 模拟客户端
func main() {
	utils.Log.Info("this is client....")

	//1.创建链接
	//conn, err := net.Dial("tcp", "127.0.0.1:8999")
	//if err != nil {
	//	utils.Log.Error("net dial error: %s", err.Error())
	//	return
	//}

	// 1.创建链接，使用tcpConn
	tcpAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8999")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		utils.Log.Error(" net.DialTCP error: %s", err.Error())
		return
	}
	defer conn.Close()

	for {
		// 1.封包, 发数据
		dp := process.NewDataPackProcess()
		binaryMsg, err := dp.Pack(process.NewMessageProcess(2, []byte("hello i am fuck 2")))
		if err != nil {
			utils.Log.Error(" package msg error: %s", err.Error())
			break
		}
		_, err = conn.Write(binaryMsg)
		if err != nil {
			utils.Log.Error("Write data error: %s", err.Error())
			break
		}

		// 2.拆包，读数据
		msg, err := dp.Unpack(conn)
		if err != nil {
			utils.Log.Error(" package Unpack error: %s", err.Error())
			break
		}

		utils.Log.Info("receive server data : %s, len : %d", msg.GetMsgData(), msg.GetMsgLen())

		time.Sleep(2 * time.Second)
	}

}
