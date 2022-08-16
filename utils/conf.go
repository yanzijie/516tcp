package utils

import (
	"encoding/json"
	"github.com/yanzijie/516tcp/inface"
	"io/ioutil"
)

// 全局参数配置

type GlobalConf struct {
	// server
	TcpServer      inface.ServerInterface // 当前全局的server对象
	Host           string                 // 服务器监听的ip
	TcpPort        int                    // 服务器监听的端口
	Name           string                 // 服务器名称
	Version        string                 // 服务版本号
	MaxConn        int                    // 当前服务的最大链接数
	MaxPackageSize uint32                 // 当前框架收发数据包的最大值
	WorkerPoolSize uint32                 // 协程池的协程数量
	MaxWorkerSize  uint32                 // 每个工作协程对应的消息队列内的任务上限值
}

var GlobalObject *GlobalConf

func init() {
	// 先写默认配置
	GlobalObject = &GlobalConf{
		Host:           "0.0.0.0",
		TcpPort:        8999,
		Name:           "516tcp",
		Version:        "v0.5",
		MaxConn:        1000,
		MaxPackageSize: 4096,
		WorkerPoolSize: 10,
		MaxWorkerSize:  1024,
	}
	//然后从配置文件里面读, 读出来就覆盖
	GlobalObject.Reload()
}

// Reload 加载自定义参数
func (g *GlobalConf) Reload() {
	//data, err := ioutil.ReadFile("demo/version/conf.json") // debug用这个
	data, err := ioutil.ReadFile("conf.json") // 实际跑用这个
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}
