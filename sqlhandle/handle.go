package sqlhandle

import (
	"errors"
	"gitee.com/mrmateoliu/wow_launch.git/utils"
	"time"
)

func IpBan(ip string, banTime int64) error {
	if !utils.GlobalObject.BanSql {
		return nil
	}
	db := utils.AuthDB
	if db == nil {
		return errors.New("账号数据库未连接,游戏内账号ip封禁失败")
	}
	err := db.Exec("INSERT INTO `ip_banned` (`ip`, `bandate`, `unbandate`, `bannedby`, `banreason`) VALUES (?, ?, ?, '登录器', '登录器封禁')",
		ip, time.Now().Unix(), banTime).Error
	if err != nil {
		return err
	}
	return nil
}

func RemoveBan(ip string) error {
	if !utils.GlobalObject.BanSql {
		return nil
	}
	db := utils.AuthDB
	if db == nil {
		return errors.New("账号数据库未连接,游戏内账号ip封禁失败")
	}
	err := db.Exec("DELETE FROM `ip_banned` WHERE ip = ?", ip).Error
	if err != nil {
		return err
	}
	return nil
}
