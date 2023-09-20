package wownet

import "gitee.com/mrmateoliu/wow_launch.git/wowiface"

// 实现基础路由时, 先嵌入这个BaseRouter基类, 然后根据需要,对这个基类的某个方法进行重写
type BaseRouter struct{}

/*说明:
1.先用BaseRouter实现抽象层IRouter所有方法
2.当其他路由需要使用BaseRouter中的方法时,直接继承BaseRouter结构体即可.
然后使用BaseRouter中任意方法即可.
这样就不需要再重新实现IRouter中全部的方法了.因为BaseRouter已经实现了抽象层的IRouter
*/

// 处理conn业务之前的Hook方法
func (br *BaseRouter) BeforeHandle(request wowiface.IRequest) {}

// 处理conn业务的主Hook方法
func (br *BaseRouter) Handle(request wowiface.IRequest) {}

// 处理conn业务之后的Hook方法
func (br *BaseRouter) AfterHandle(request wowiface.IRequest) {}
