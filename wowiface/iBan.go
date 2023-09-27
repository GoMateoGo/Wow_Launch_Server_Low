package wowiface

// 封禁抽象层
type IBan interface {
	//增加封禁
	AddBan(ip, mac string, expireTime int64)
	//获取具体封禁信息
	GetBanByInfo(ip, Mac string) interface{}
	//移除封禁
	RemoveBan(ipOrMac string) bool
}
