package utils

import (
	"encoding/json"
	"gitee.com/mrmateoliu/wow_launch.git/wowiface"
	"os"
)

/*
存储一切有关框架的全局参数, 供其他模块使用
一些参数是可以通过XXX.json由用户进行配置
*/

type GlobalObj struct {
	/*
		Server
	*/
	TcpServer        wowiface.IServer //当前全局的Server对象
	Host             string           //当前服务器主机监听的IP
	TcpPort          int              //当前服务器主机监听的端口号
	Name             string           //当前服务器的名称
	Version          string           //当前的版本号
	MaxConn          int              //当前服务器主机允许的最大连接数
	MaxPackageSize   uint32           //当前框架数据包的最大值
	WorkerPoolSize   uint32           //当前业务工作池的Goroutine数量
	MaxWorkerTaskLen uint32           //框架允许用户最大开辟多少个Worker
}

/*
定义一个全局的对外 GlobalObj
*/
var GlobalObject *GlobalObj

// 从 xxx.json去加载用户自定义参数
func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/config.json")
	if err != nil {
		panic(err)
	}
	//将json文件数据解析到GlobalObj struc中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// 提供一个init方法,初始化当前GlobalObj
func init() {
	//如果配置文件没有加载,先生成默认
	GlobalObject = &GlobalObj{
		Host:             "0.0.0.0",
		TcpPort:          8999,
		Name:             "wow_launch_App",
		Version:          "v0.1",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,   //表示整个workerPool的数量
		MaxWorkerTaskLen: 1024, //每个worker对应消息队列(channel)的最大的数量值
	}

	//应该尝试从conf/xxx.json去加载一些用户自定义参数
	GlobalObject.Reload()
}
