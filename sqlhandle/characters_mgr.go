package sqlhandle

import (
	"errors"
	"fmt"
	"gitee.com/mrmateoliu/wow_launch.git/utils"
)

type Characters struct {
	account    string
	password   string
	playerName string
}

// 构造
func NewCharacters(account, pwd, playerName string) *Characters {
	return &Characters{
		account:    account,
		password:   pwd,
		playerName: playerName,
	}
}

// 查询角色是否存在
func (s *Characters) HasCharacter(accId uint) string {
	db := utils.CharaDB
	if db == nil {
		return ""
	}
	var playerName string

	db.Raw("SELECT name FROM characters WHERE account=? and name=? ", accId, s.playerName).Scan(&playerName)

	return playerName
}

// 角色解卡
func (s *Characters) UnLockPlayer(accId uint) error {
	db := utils.CharaDB
	if db == nil {
		return errors.New("连接数据库失败.请联系游戏管理员")
	}

	tx := db.Begin()

	// 执行原生 SQL 操作
	tx.Exec("UPDATE characters SET `position_x`=-3962.39, `position_y`=-2015.55, `position_z`=96.2567,`map`=1 WHERE name=? AND account=?", s.playerName, accId)

	if err := tx.Commit().Error; err != nil {
		// 回滚
		tx.Rollback()
		utils.Logger.Error(fmt.Sprintf("[解卡失败] 账号:%s,密码:%s,角色名称:%s", s.account, s.password, s.playerName))
		return errors.New("解卡失败失败,请联系游戏管理员")
	}

	return nil
}
