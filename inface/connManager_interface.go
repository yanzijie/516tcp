package inface

// 链接管理模块

type ConnManagerInterface interface {
	AddConn(conn ConnectionInterface)                   // 添加链接
	RemoveConn(conn ConnectionInterface)                // 删除链接
	GetConn(connID uint32) (ConnectionInterface, error) // 根据链接ID获取链接
	GetConnNumberLen() int                              // 获取链接总数
	ClearConn()                                         // 中止并清除全部的链接
}
