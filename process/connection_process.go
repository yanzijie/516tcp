package process

import (
	"github.com/yanzijie/516tcp/inface"
	"github.com/yanzijie/516tcp/utils"
	"net"
	"sync"
)

/*
链接模块
*/

type ConnectionProcess struct {
	Conn        *net.TCPConn      // 当前链接的socket套接字
	ConnID      uint32            // 链接的ID
	ConnIsClose bool              // 当前链接状态(是否已经关闭)
	HandleApi   inface.HandleFunc // 与当前链接绑定的业务方法
	ExitChan    chan bool         // 管理当前链接是否退出的channel
}

// NewConnection 初始化链接
// conn-客户端socket链接, connID-链接id, callbackApi-链接回调函数
func NewConnection(conn *net.TCPConn, connID uint32, callbackApi inface.HandleFunc) *ConnectionProcess {
	c := &ConnectionProcess{
		Conn:        conn,
		ConnID:      connID,
		ConnIsClose: false,
		HandleApi:   callbackApi,
		ExitChan:    make(chan bool, 1),
	}

	return c
}

// StartConnection 启动链接
func (c *ConnectionProcess) StartConnection() {
	utils.Log.Info("start Conn, connId: %d\n", c.ConnID)
	// 读取当前链接携带的数据
	go c.StartRead()

	// 回写数据给客户端 TODO
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
	close(c.ExitChan)  // 回收资源
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

// Send 发送数据
func (c *ConnectionProcess) Send(data []byte) error {

	return nil
}

// StartRead 读取当前链接携带的数据
func (c *ConnectionProcess) StartRead() {
	utils.Log.Info("read goroutine running connId: %d\n", c.ConnID)
	defer c.StopConnection()
	defer utils.Log.Info("read goroutine is exit, client addr is:%s\n", c.GetClientAddr().String())

	// 现在先把读写都放这里，后面再拆分
	for {
		// 读取客户端数据到buf中
		buf := make([]byte, 512)
		dataLen, err := c.Conn.Read(buf)
		if err != nil && err.Error() != "EOF" {
			utils.Log.Error(" read data from client error: %s, connID:%d\n", err.Error(), c.ConnID)
			continue
		}
		//如果没收到数据，就继续读
		if dataLen == 0 {
			continue
		}
		utils.Log.Info("receive client data: %s, len: %d\n", string(buf), dataLen)

		//调用处理业务的方法
		err = c.HandleApi(c.Conn, buf, dataLen)
		if err != nil {
			utils.Log.Error(" handlerApi error: %s, connID:%d\n", err.Error(), c.ConnID)
			break
		}
	}
}
