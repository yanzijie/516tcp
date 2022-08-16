package main

import (
	"github.com/yanzijie/516tcp/process"
	"github.com/yanzijie/516tcp/utils"
	"math/rand"
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
		// 取随机数
		rand.Seed(time.Now().UnixNano())
		//num := rand.Intn(3)
		//if num <= 0 || num > 3 {
		//	continue
		//}
		binaryMsg, err := dp.Pack(process.NewMessageProcess(3, []byte("hello i am fuck 3")))
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
