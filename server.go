package main

import (
	"fmt"
	"gitee.com/mrmateoliu/wow_launch.git/conf"
	"gitee.com/mrmateoliu/wow_launch.git/dhttp"
	"gitee.com/mrmateoliu/wow_launch.git/socket_router"
	"gitee.com/mrmateoliu/wow_launch.git/utils"
	"gitee.com/mrmateoliu/wow_launch.git/wowiface"
	"gitee.com/mrmateoliu/wow_launch.git/wownet"
	"strconv"
)

// 创建链接之后的钩子函数
func DoConnectionBegin(conn wowiface.IConnection) {
	if utils.GlobalObject.Develop {
		fmt.Println("===>创建链接钩子已经调用")
	}
	//获取用户mac地址和os版本
	if err := conn.SendMsg(101, []byte("GetClientInfo")); err != nil {
		fmt.Println(err)
	}

	// --------------------------------------------------------
	// - 测试 链接属性设置
	//给这个链接设置一些属性
	if utils.GlobalObject.Develop {
		fmt.Println("----->设置链接属性<-----")
		conn.SetProperty("Name", "魔兽登录器")
		conn.SetProperty("网址", "https://www.getgamesf.com")
	}
}

// 链接断开前的钩子函数
func DoConnectionLost(conn wowiface.IConnection) {
	if utils.GlobalObject.Develop {
		fmt.Println("===>关闭连接钩子已经调用,关闭链接Id=", conn.GetConnId())
	}

	serverOwner, err := utils.SServer.GetConnMgr().GetServerOwner()
	if err != nil {
		return
	}

	//通知管理UI断开连接
	if err := serverOwner.SendMsg(105, []byte(strconv.Itoa(int(conn.GetConnId())))); err != nil {
		fmt.Println(err)
	}
	// --------------------------------------------------------
	// - 测试 链接属性读取
	//if re, err := conn.GetProperty("Name"); err == nil {
	//	fmt.Println(re)
	//}
	//if re, err := conn.GetProperty("网址"); err == nil {
	//	fmt.Println(re)
	//}
}

func main() {
	// --------------0. 初始化配置------------------------------------------
	SysInit()

	// --------------1.验证时间过期------------------------------------------
	dhttp.HandCallDll()

	// --------------2.启用http服务用于下载------------------------------------------
	dhttp.RunHttp()
	//---------------3.创建SocketServer句柄--------------------------------------
	s := wownet.NewServer()
	utils.SServer = s

	// --------------4. 注册连接的Hook方法------------------------------------------
	s.SetAfterStartConn(DoConnectionBegin)
	s.SetBeforeStopConn(DoConnectionLost)

	// --------------5.给当前框架添加自定义router-----------------------------------
	socket_router.RegisterRouter(s)

	// --------------6.启动服务器-----------------------------------
	s.Server() //阻塞
}

func SysInit() {
	// ========================================================
	// =初始化日志配置
	utils.Logger = conf.InitLogger()

	// ========================================================
	// = 读取ban.txt文件到内存中
	conf.ReadBanList()

	// ========================================================
	// =初始化账号数据链接
	if db, err := conf.InitAuthDB(); err != nil {
		if err != nil {
			fmt.Println("账号数据库连接失败...")
		}
	} else {
		utils.AuthDB = db
	}

	// ========================================================
	// =初始化角色数据库连接
	if db, err := conf.InitCharaDB(); err != nil {
		if err != nil {
			fmt.Println("角色数据库连接失败...")
		}
	} else {
		utils.CharaDB = db
	}
}
