package process

import "github.com/yanzijie/516tcp/inface"

// 具体的业务路由，继承（嵌入）基础路由BaseRouter, 然后重写方法
// 所以这里的实现就空着就行, 由具体业务去重写

type BaseRouter struct {
}

// PreHandle 处理业务之前的方法（钩子）
func (b *BaseRouter) PreHandle(req inface.RequestInterface) {

}

// Handle 处理业务的主要方法
func (b *BaseRouter) Handle(req inface.RequestInterface) {

}

// PostHandle 处理业务之后的方法（钩子）
func (b *BaseRouter) PostHandle(req inface.RequestInterface) {

}
