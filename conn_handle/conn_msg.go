package conn_handle

import (
	"fmt"
	"gitee.com/mrmateoliu/wow_launch.git/wowiface"
	"gitee.com/mrmateoliu/wow_launch.git/wownet"
	"strings"
)

// 检查是否被ban
func IsBanned(conn wowiface.IConnection) bool {
	sBan := wownet.BanInstance
	// 使用 ":" 分割字符串
	var cIp string
	partIp := strings.Split(conn.RemoteAddr().String(), ":")
	if len(partIp) >= 1 {
		cIp = partIp[0]
	} else {
		fmt.Println("解析客户端ip错误.")
		return true
	}
	res := sBan.GetBanByInfo(cIp, conn.GetConnMac())
	if res == nil {
		return false
	}
	return true
}
