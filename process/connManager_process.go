package process

import (
	"errors"
	"github.com/yanzijie/516tcp/inface"
	"github.com/yanzijie/516tcp/utils"
	"strconv"
	"sync"
)

// 链接管理模块

type ConnManagerProcess struct {
	/*
		0 - 链接1
		1 - 链接2
		2 - 链接3
		3 - 链接4
		....
	*/
	connections map[uint32]inface.ConnectionInterface // 创建好的Connection集合
	connLock    sync.RWMutex
}

// NewConnManager 创建当前链接管理
func NewConnManager() *ConnManagerProcess {
	return &ConnManagerProcess{
		connections: make(map[uint32]inface.ConnectionInterface),
	}
}

func (c *ConnManagerProcess) AddConn(conn inface.ConnectionInterface) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	c.connections[conn.GetConnID()] = conn
	utils.Log.Info(" connection add to connManager success, connID = %d, connNumLen = %d",
		conn.GetConnID(), c.GetConnNumberLen())
}

func (c *ConnManagerProcess) RemoveConn(conn inface.ConnectionInterface) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	delete(c.connections, conn.GetConnID())
	utils.Log.Info(" connection delete from connManager success, connID = %d, connNumLen = %d",
		conn.GetConnID(), c.GetConnNumberLen())
}

func (c *ConnManagerProcess) GetConn(connID uint32) (inface.ConnectionInterface, error) {
	c.connLock.RLock()
	defer c.connLock.RUnlock()

	conn, ok := c.connections[connID]
	if ok {
		return conn, nil
	} else {
		id := strconv.Itoa(int(connID))
		return nil, errors.New(" connID:" + id + " not register")
	}
}

func (c *ConnManagerProcess) GetConnNumberLen() int {
	return len(c.connections)
}

func (c *ConnManagerProcess) ClearConn() {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	for connID, conn := range c.connections {
		conn.StopConnection()         // 中止
		delete(c.connections, connID) // 删除
	}

	utils.Log.Info(" clear all connection, connNumLen = %d", c.GetConnNumberLen())
}
