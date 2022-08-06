package main

import (
	"github.com/yanzijie/516tcp/utils"
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

	for {
		//2.写数据给服务端
		_, err = conn.Write([]byte("hello i am coming"))
		if err != nil {
			utils.Log.Error("Write data error: %s", err.Error())
			return
		}

		//3.接收服务器数据
		buf := make([]byte, 512)
		// dataLen:读取到的数据长度
		dataLen, err := conn.Read(buf)
		if err != nil {
			utils.Log.Error("read data error: %s", err.Error())
			return
		}
		utils.Log.Info("receive server data : %s, len : %d", string(buf), dataLen)

		time.Sleep(2 * time.Second)
	}
}
