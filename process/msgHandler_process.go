package process

import (
	"github.com/yanzijie/516tcp/inface"
	"github.com/yanzijie/516tcp/utils"
)

// MsgHandlerProcess 处理消息的结构
type MsgHandlerProcess struct {
	// msgID与router路由对应关系map
	ApiMap map[uint32]inface.RouterInterface
	// 取任务的消息队列, 管道的数组切片[x,x,x,x,x],下标对应每个工作协程的序号, 每个队列对应一个协程
	// 数组内的数据是 inface.RequestInterface 类型的通道
	TaskQueen []chan inface.RequestInterface
	// 协程池内的协程数量
	WorkerPoolSize uint32
}

func NewMsgHandlerProcess() *MsgHandlerProcess {
	return &MsgHandlerProcess{
		ApiMap:         make(map[uint32]inface.RouterInterface),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueen:      make([]chan inface.RequestInterface, utils.GlobalObject.WorkerPoolSize),
	}
}

func (m *MsgHandlerProcess) DoMsgHandler(req inface.RequestInterface) {
	//1. 判断msgId
	handler, ok := m.ApiMap[req.GetMsgId()]
	if !ok {
		utils.Log.Error(" this msgId = %d is not register", req.GetMsgId())
		return
	}
	//2. 根据msgId调用对应的router
	handler.PreHandle(req)
	handler.Handle(req)
	handler.PostHandle(req)
}

func (m *MsgHandlerProcess) AddRouter(msgId uint32, router inface.RouterInterface) {
	// 校验msgID
	_, ok := m.ApiMap[msgId]
	if ok {
		// msgID已经注册
		utils.Log.Error(" this msgId = %d is register", msgId)
		return
	}
	m.ApiMap[msgId] = router
	utils.Log.Info(" msgID = %d register success!", msgId)
}

func (m *MsgHandlerProcess) StartWorkerPool() {
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		// 给worker协程对应的 channel消息管道 开辟空间
		m.TaskQueen[i] = make(chan inface.RequestInterface, utils.GlobalObject.MaxWorkerSize)
		// 启动当前工作协程
		go m.StartOneWorker(i, m.TaskQueen[i])
	}

}

// StartOneWorker   workerID: 协程ID,  taskQueen: 该协程对应的消息队列(请求队列)
func (m *MsgHandlerProcess) StartOneWorker(workerID int, taskQueen chan inface.RequestInterface) {
	utils.Log.Info(" workerID = %d was start", workerID)

	// 阻塞从消息队列中读取请求
	for {
		select {
		case request := <-taskQueen:
			// 请求来了，从channel中取出来，发给对应路由处理
			m.DoMsgHandler(request)
		}
	}
}

func (m *MsgHandlerProcess) SendMsgToTaskQueen(request inface.RequestInterface) {
	// 1. 简单负载均衡(取余轮询分配)，稍微平均的分配给工作协程
	// 其实还可以查一下每个消息队列，给较少的队列分配任务
	// 先根据ConnID来分配
	workerID := request.GetConnection().GetConnID() % m.WorkerPoolSize
	utils.Log.Info(" msgID: %d to workerID: %d process", request.GetMsgId(), workerID)

	// 2. 发送请求给对应的消息队列
	m.TaskQueen[workerID] <- request
}
