package socket_router

import (
	"fmt"
	"gitee.com/mrmateoliu/wow_launch.git/conn_handle"
	"gitee.com/mrmateoliu/wow_launch.git/sqlhandle"
	"gitee.com/mrmateoliu/wow_launch.git/utils"
	"gitee.com/mrmateoliu/wow_launch.git/wowiface"
	"gitee.com/mrmateoliu/wow_launch.git/wownet"
	"strconv"
	"strings"
	"time"
)

// 注册自定义路由
func RegisterRouter(s wowiface.IServer) {
	//发送过期时间
	s.AddRouter(100, &SendExpireTime{Server: s})
	//获取用户链接信息
	s.AddRouter(101, &HandClientConnRouter{Server: s})
	//关闭指定连接
	s.AddRouter(102, &CloseCoon{Server: s})
	//Ban指定ip和mac
	s.AddRouter(103, &BanClient{Server: s})
	//解封Ban
	s.AddRouter(104, &RemoveBanClient{})
	//注册账号
	s.AddRouter(1, &RegisterAccount{})
	//修改/找回密码
	s.AddRouter(2, &ChangePassword{})
	//角色解卡
	s.AddRouter(3, &PlayerUnLock{})
}

// 角色解卡
type PlayerUnLock struct {
	wownet.BaseRouter
}

// 角色解卡
func (s *PlayerUnLock) AfterHandle(re wowiface.IRequest) {
	//0. 解析账号,密码,角色名称
	var errStr string
	var userName string   //账号
	var Pwd string        //密码
	var playerName string //角色名称
	msgDat := re.GetData()
	parts := strings.Split(string(msgDat), "#")
	if len(parts) == 3 {
		userName = parts[0]
		Pwd = parts[1]
		playerName = parts[2]
	} else {
		errStr = fmt.Sprintf("[角色解卡]失败原因,账号密码解析错误,解卡账号:%s,解卡角色名称:%s", userName, playerName)
		utils.Logger.Error(errStr)
		_ = re.GetConnection().SendMsg(3, []byte(errStr)) //解卡失败.
		return
	}
	//1. 验证账号密码正确性()
	acc := sqlhandle.NewUserData(userName, Pwd, "")
	accId, err := acc.CheckSaltAndVerifier()
	if err != nil {
		_ = re.GetConnection().SendMsg(3, []byte(err.Error())) //解卡失败.
		return
	}
	if accId == 0 {
		_ = re.GetConnection().SendMsg(3, []byte("没有找到对应的账号")) //解卡失败.
		return
	}
	//2. 检查角色是否存在
	ch := sqlhandle.NewCharacters(userName, Pwd, playerName)
	pName := ch.HasCharacter(accId)
	if pName == "" {
		_ = re.GetConnection().SendMsg(3, []byte("没有找到该角色")) //解卡失败.
		return
	}
	//3. 解卡
	if err = ch.UnLockPlayer(accId); err != nil {
		_ = re.GetConnection().SendMsg(3, []byte("解卡失败,请联系游戏管理员")) //解卡失败.
		return
	}
	_ = re.GetConnection().SendMsg(3, []byte("解卡完成!!")) //解卡成功.
}

type SendExpireTime struct {
	wownet.BaseRouter
	Server wowiface.IServer
}

func (s *SendExpireTime) AfterHandle(re wowiface.IRequest) {

	go func() {
		for {
			serverOwner, err := utils.SServer.GetConnMgr().GetServerOwner()
			if err != nil {
				fmt.Println("服务端UI管理离线", err.Error())
				return
			}
			_ = serverOwner.SendMsg(100, []byte(strconv.FormatInt(utils.RemainTimeSecond, 10))) //发送网关剩余时间
			time.Sleep(1 * time.Second)
		}
	}()
}

type ChangePassword struct {
	wownet.BaseRouter
}

func (s ChangePassword) AfterHandle(request wowiface.IRequest) {
	var userName string   //账号
	var newPwd string     //新密码
	var secureCode string //安全码
	msgDat := request.GetData()
	parts := strings.Split(string(msgDat), "#")
	if len(parts) == 3 {
		userName = parts[0]
		newPwd = parts[1]
		secureCode = parts[2]
	} else {
		utils.Logger.Error(fmt.Sprintf("修改/找回密码失败,对方客户端ip:%s", request.GetConnection().RemoteAddr()))
		return
	}
	//修改/找回密码
	user := sqlhandle.NewUserData(userName, newPwd, secureCode)
	err := user.ChangePassword()
	if err != nil {
		err = request.GetConnection().SendMsg(1, []byte(err.Error())) //注册结果
	}
}

type RegisterAccount struct {
	wownet.BaseRouter
}

func (s RegisterAccount) AfterHandle(request wowiface.IRequest) {
	var userName string   //账号
	var password string   //密码
	var secureCode string //安全码

	msgDat := request.GetData()
	parts := strings.Split(string(msgDat), "#")
	if len(parts) == 3 {
		userName = parts[0]
		password = parts[1]
		secureCode = parts[2]
	} else {
		utils.Logger.Error(fmt.Sprintf("注册账号失败,对方客户端ip:%s", request.GetConnection().RemoteAddr()))
		return
	}
	cIp := request.GetConnection().RemoteAddr().String()
	// 使用 ":" 分割字符串
	partIp := strings.Split(cIp, ":")
	if len(partIp) >= 1 {
		cIp = partIp[0]
	} else {
		fmt.Println("解析客户端ip错误.")
	}
	//注册账号密码
	user := sqlhandle.NewUserData(userName, password, secureCode)
	err := user.CreateAccount(cIp)
	if err != nil {
		err = request.GetConnection().SendMsg(1, []byte(err.Error())) //注册结果
	}
}

// 解封Ban客户端
type RemoveBanClient struct {
	wownet.BaseRouter
}

// 解封Ban客户端
func (s RemoveBanClient) Handle(request wowiface.IRequest) {
	var ipAddr string
	var macAddr string
	sBan := wownet.BanInstance
	msgDat := request.GetData()
	parts := strings.Split(string(msgDat), "#")
	if len(parts) == 2 {
		// parts[0]="ip地址" parts[1]="mac地址"
		ipAddr = parts[0]
		macAddr = parts[1]
	} else {
		utils.Logger.Error(fmt.Sprintf("Ban客户端时信息错误:%s", parts))
		return
	}

	sBan.RemoveBan(macAddr)
	err := sqlhandle.RemoveBan(ipAddr)
	serverOwner, err := utils.SServer.GetConnMgr().GetServerOwner()
	if err != nil {
		return
	}
	if err != nil {
		err = serverOwner.SendMsg(1, []byte(err.Error())) //解封ban结果
	}
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("账号数据库解封BanIp时发生错误,Ip地址:%s,错误信息:%s", ipAddr, err.Error()))
	}
}

// Ban客户端
type BanClient struct {
	wownet.BaseRouter
	Server wowiface.IServer
}

// Ban客户端
func (s *BanClient) Handle(request wowiface.IRequest) {
	sBan := wownet.BanInstance
	var ConnId uint32
	var ipAddr string
	var macAddr string
	var expireTime int64 //应该是秒
	// 1.获取消息
	msgDat := request.GetData()
	parts := strings.Split(string(msgDat), "#")
	if len(parts) == 4 {
		// parts[0]="连接id" parts[1]="ip地址" parts[2]="mac地址" parts[3]="过期时间"
		ipAddr = parts[1]
		macAddr = parts[2]
		// 使用 strconv.ParseInt 函数将字符串转换为 int64
		reTime, err := strconv.ParseInt(parts[3], 10, 64)
		reTime += time.Now().Unix()
		if err != nil {
			utils.Logger.Error(fmt.Sprintf("Ban客户端时候Ip和Mac信息解析错误:%s", err))
			return
		}
		expireTime = reTime
		// 使用 strconv.ParseUint 函数将字符串转换为 uint32
		num, err := strconv.ParseUint(parts[0], 10, 32)
		if err != nil {
			utils.Logger.Error(fmt.Sprintf("Ban客户端时候链接Id解析错误:%s", err))
			return
		}
		ConnId = uint32(num)
	} else {
		utils.Logger.Error(fmt.Sprintf("Ban客户端时信息错误:%s", parts))
		return
	}
	conn, err := s.Server.GetConnMgr().Get(ConnId)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("Ban客户端时:%s", err))
		return
	}
	conn.Stop()
	sBan.AddBan(ipAddr, macAddr, expireTime)
	t := time.Unix(expireTime, 0)
	fmt.Printf("封禁成功:Ip地址:%s,Mac地址:%s,封禁到:%s", ipAddr, macAddr, t.Format(time.DateTime))
	err = sqlhandle.IpBan(ipAddr, expireTime)
	if err != nil {
		err = conn.SendMsg(1, []byte(err.Error())) //ban结果
	}
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("账号数据库Ban时发生错误,Ip地址:%s,Ban时间:%d,错误信息:%s", ipAddr, expireTime, err.Error()))
	}
}

// 获取客户端信息
type HandClientConnRouter struct {
	wownet.BaseRouter
	Server wowiface.IServer
}

// 获取客户端信息
func (s *HandClientConnRouter) AfterHandle(request wowiface.IRequest) {

	// 1.获取消息
	msgDat := request.GetData()
	// 2.使用 "#" 分割字符串
	var mac string
	var Os string
	parts := strings.Split(string(msgDat), "#")
	if len(parts) == 2 {
		// parts[0]="mac地址" parts[1]="系统Id"
		mac = parts[0]
		Os = parts[1]
		request.GetConnection().SetConnMac(mac)
		request.GetConnection().SetConnOs(Os)
		//此处检查是否被封禁
		if conn_handle.IsBanned(request.GetConnection()) {
			_ = request.GetConnection().SendMsg(1001, []byte("您已被封禁")) //ban提示
			request.GetConnection().Stop()
			return
		}
	} else {
		utils.Logger.Error("信息切割错误.", parts)
		request.GetConnection().Stop()
		return
	}

	conn, err := s.Server.GetConnMgr().GetServerOwner()
	if err != nil {
		fmt.Println("没有找到服务端UI管理连接", err.Error())
		return
	}

	if conn.GetConnMac() == utils.SelfMac {
		//return
	}

	cIp := request.GetConnection().GetTCPConnection().RemoteAddr().String()
	// 使用 ":" 分割字符串
	partIp := strings.Split(cIp, ":")
	if len(partIp) >= 1 {
		cIp = partIp[0]
	} else {
		fmt.Println("解析客户端ip错误.")
	}

	result := strconv.Itoa(int(request.GetConnection().GetConnId())) + "#" + cIp + "#" + mac + "#" + Os

	//发送信息给服务端管理UI
	err = conn.SendMsg(102, []byte(result))
	if err != nil {
		fmt.Println("发送给服务端UI消息错误:", err)
		return
	}
}

// 关闭当前链接
type CloseCoon struct {
	wownet.BaseRouter
	Server wowiface.IServer
}

// 关闭当前链接
func (s *CloseCoon) AfterHandle(request wowiface.IRequest) {
	msgId := string(request.GetData())
	uint32Value, err := strconv.ParseUint(msgId, 10, 32)
	if err != nil {
		// 处理转换失败的情况
		fmt.Println("转换失败:", err)
	} else {
		// 转换成功，uint32Value 现在包含转换后的值
		uint32Result := uint32(uint32Value)

		conn, err := s.Server.GetConnMgr().Get(uint32Result)
		if err != nil {
			fmt.Println("没有找到当前服务器UI:", err)
			return
		}
		conn.Stop()
	}
}
