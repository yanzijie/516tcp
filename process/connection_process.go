package process

import (
	"errors"
	"github.com/yanzijie/516tcp/inface"
	"github.com/yanzijie/516tcp/utils"
	"net"
	"sync"
)

/*
链接模块
*/

type ConnectionProcess struct {
	FatherServer inface.ServerInterface     // 当前链接所属的server
	Conn         *net.TCPConn               // 当前链接的socket套接字
	ConnID       uint32                     // 链接的ID
	ConnIsClose  bool                       // 当前链接状态(是否已经关闭)
	ExitChan     chan bool                  // 管理当前链接是否退出的channel, read协程告诉write协程是否退出
	msgChan      chan []byte                //无缓冲channel，用于读写协程之间消息通信
	MsgHandler   inface.MsgHandlerInterface // 该链接的对应处理方法(handler)
}

// NewConnection 初始化链接
// conn-客户端socket链接, connID-链接id, msgHandler-消息处理模块
func NewConnection(server inface.ServerInterface, conn *net.TCPConn, connID uint32, msgHandler inface.MsgHandlerInterface) *ConnectionProcess {
	c := &ConnectionProcess{
		FatherServer: server,
		Conn:         conn,
		ConnID:       connID,
		ConnIsClose:  false,
		ExitChan:     make(chan bool, 1),
		msgChan:      make(chan []byte),
		MsgHandler:   msgHandler,
	}

	c.FatherServer.GetConnManager().AddConn(c)

	return c
}

// StartConnection 启动链接
func (c *ConnectionProcess) StartConnection() {
	utils.Log.Info("start Conn, connId: %d\n", c.ConnID)
	// 读取当前链接携带的数据
	go c.StartRead()
	// 回写数据给客户端
	go c.StartWrite()
}

// StopConnection 停止链接
func (c *ConnectionProcess) StopConnection() {
	utils.Log.Info("stop Conn, connId: %d\n", c.ConnID)
	if c.ConnIsClose {
		return
	}
	var lock sync.Mutex

	lock.Lock()
	c.ConnIsClose = true
	lock.Unlock()

	_ = c.Conn.Close() // 关闭链接
	c.ExitChan <- true // 告诉write协程退出

	// 移除链接管理模块
	c.FatherServer.GetConnManager().RemoveConn(c)

	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)
}

// GetTCPSocketConn 获取当前链接绑定的socket conn
func (c *ConnectionProcess) GetTCPSocketConn() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前链接的链接ID
func (c *ConnectionProcess) GetConnID() uint32 {
	return c.ConnID
}

// GetClientAddr 获取客户端的tcp状态，ip和port
func (c *ConnectionProcess) GetClientAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// Send 发送数据, 先封包，再发送
func (c *ConnectionProcess) Send(msgId uint32, data []byte) error {
	if c.ConnIsClose {
		return errors.New("connection is close")
	}

	//封包
	dp := NewDataPackProcess()
	binaryMsg, err := dp.Pack(NewMessageProcess(msgId, data))
	if err != nil {
		utils.Log.Error(" Pack msg error: %s", err.Error())
		return err
	}

	c.msgChan <- binaryMsg // 发给写goroutine进行处理

	return nil
}

// StartRead 读取当前链接携带的数据
func (c *ConnectionProcess) StartRead() {
	utils.Log.Info("[read goroutine running connId: %d]", c.ConnID)
	defer c.StopConnection()
	defer utils.Log.Info("[read goroutine is exit, client addr is:%s, connId: %d]", c.GetClientAddr().String(), c.ConnID)

	// 现在先把读写都放这里，后面再拆分
	for {
		// 拆包读取数据
		dp := NewDataPackProcess()
		msg, err := dp.Unpack(c.GetTCPSocketConn())
		if err != nil {
			utils.Log.Error("Unpack data error: %s", err.Error())
			break
		}

		// 得到当前链接对象的Request数据
		req := &RequestProcess{
			Conn:    c,
			ReqData: msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 那请求发送给消息队列，然后让协程池的协程去处理
			c.MsgHandler.SendMsgToTaskQueen(req)
		} else {
			// 如果协程池没有协程，就直接自己开一个协程处理
			go c.MsgHandler.DoMsgHandler(req)
		}

	}
}

// StartWrite 回写数据给客户端
func (c *ConnectionProcess) StartWrite() {
	utils.Log.Info("[write goroutine running connId: %d]", c.ConnID)
	defer utils.Log.Info("[write goroutine is exit, client addr is:%s, connId: %d]", c.GetClientAddr().String(), c.ConnID)

	for {
		select {
		case data := <-c.msgChan:
			// 有数据来了，回写客户端
			_, err := c.Conn.Write(data)
			if err != nil {
				utils.Log.Error("Conn.Write error: %s", err.Error())
				return
			}
		case <-c.ExitChan:
			// read已经退出了，write也退出
			return
		}

	}
}
