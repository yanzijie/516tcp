package inface

// 路由抽象接口, 路由里面的数据都是RequestInterface

type RouterInterface interface {
	// PreHandle 处理业务之前的方法（钩子）
	PreHandle(req RequestInterface)
	// Handle 处理业务的主要方法
	Handle(req RequestInterface)
	// PostHandle 处理业务之后的方法（钩子）
	PostHandle(req RequestInterface)
}
