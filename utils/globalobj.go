package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gitee.com/mrmateoliu/wow_launch.git/wowiface"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io"
	"net"
	"os"
	"strings"
)

var (
	Logger           *zap.SugaredLogger
	AuthDB           *gorm.DB
	CharaDB          *gorm.DB
	SServer          wowiface.IServer
	RemainTimeSecond int64 = 1
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
	HttpPort         int              //http服务器端口
	Name             string           //当前服务器的名称
	Version          string           //当前的版本号
	MaxConn          int              //当前服务器主机允许的最大连接数
	MaxPackageSize   uint32           //当前框架数据包的最大值
	WorkerPoolSize   uint32           //当前业务工作池的Goroutine数量
	MaxWorkerTaskLen uint32           //框架允许用户最大开辟多少个Worker
	LogMaxSize       uint32           //日志文件最大尺寸(M),超限后自动分隔
	LogMaxBackups    uint32           //保留旧文件的最大个数
	LogMaxAge        uint32           //保留旧文件的最大天数
	Develop          bool             //是否为开发者模式
	AuthDsn          string           //数据库连接 1 角色库
	CharaDsn         string           //数据库链接2 角色库
	AuthMaxIdleConn  int              //最多空闲链接数
	AuthMaxOpenConn  int              //最多打开连接数
	BanSql           bool             //是否也禁用数据库(账号)ip地址,需要链接数据库
	WebSite          string           //网站主页
	PayWebSite       string           //充值页面
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
		Host:             "127.0.0.1",
		TcpPort:          8999,
		HttpPort:         8888, //http服务器
		Name:             "wow_launch_App",
		Version:          "v0.1",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,   //表示整个workerPool的数量
		MaxWorkerTaskLen: 1024, //每个worker对应消息队列(channel)的最大的数量值
		LogMaxSize:       1,    //日志文件最大尺寸(M),超限后自动分隔
		LogMaxBackups:    10,   //保留旧文件的最大个数
		LogMaxAge:        90,   //保留旧文件的最大天数
		Develop:          true, //是否为开发者模式
		AuthDsn:          "root:root@tcp(127.0.0.1:3307)/wow_launch?charset=utf8mb4&parseTime=True&loc=Local",
		CharaDsn:         "root:root@tcp(127.0.0.1:3307)/wow_launch?charset=utf8mb4&parseTime=True&loc=Local",
		AuthMaxIdleConn:  100,   //最多空闲链接数
		AuthMaxOpenConn:  100,   //最多打开连接数
		BanSql:           false, //是否也禁用数据库(账号)ip地址,需要链接数据库
		WebSite:          "",    //网站主页
		PayWebSite:       "",    //充值页面
	}

	//应该尝试从conf/xxx.json去加载一些用户自定义参数
	GlobalObject.Reload()
}

// 获取本机Mac地址
func GetMACAddress() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("无法获取网络接口: %v", err)
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			addrs, err := iface.Addrs()
			if err != nil {
				return "", fmt.Errorf("无法获取网卡地址 %s: %v", iface.Name, err)
			}

			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					hwAddr := iface.HardwareAddr.String()
					return strings.ToUpper(hwAddr), nil
				}
			}
		}
	}

	return "", fmt.Errorf("找不到有效的MAC地址")
}

func CheckMD5(md5Value string) bool {
	filePath := "downloaded.zip"

	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return false
	}

	calculatedMD5 := hex.EncodeToString(hash.Sum(nil))
	return strings.EqualFold(calculatedMD5, md5Value)
}
