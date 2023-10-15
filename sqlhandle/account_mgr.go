package sqlhandle

import (
	"errors"
	"fmt"
	"gitee.com/mrmateoliu/wow_launch.git/model"
	"gitee.com/mrmateoliu/wow_launch.git/utils"
	"strings"
	"time"
)

type UserData struct {
	userName   string
	password   string
	secureCode string
}

func NewUserData(userName, Pwd, secureCode string) *UserData {
	return &UserData{
		userName:   userName,
		password:   Pwd,
		secureCode: secureCode,
	}
}

// 用户名是否存在
func (s *UserData) HasUserName() bool {
	db := utils.AuthDB
	if db == nil {
		return false
	}
	var userExists bool
	if err := db.Table("account").Where("username = ?", s.userName).Select("1").Scan(&userExists).Error; err != nil {
		// 处理错误
		return false
	}
	return userExists
}

// 注册账号
func (s *UserData) CreateAccount(ip string) error {
	db := utils.AuthDB
	if db == nil {
		return errors.New("连接数据库失败.请联系游戏管理员")
	}
	if len(s.userName) > 12 || len(s.password) > 12 || len(s.secureCode) > 12 {
		return errors.New("账号或密码长度不能大于12位")
	}
	if s.HasUserName() {
		return errors.New("账号已被注册")
	}

	// 先计算SHA1
	value := fmt.Sprintf("%s:%s", strings.ToUpper(s.userName), strings.ToUpper(s.password))
	hash := utils.ToHashSHA([]byte(value))
	// 生成Salt
	salt := utils.MakeSalt()
	//转二进制
	nowSalt := utils.FromBigSalt(salt)
	// 生成Verifier
	verifier := utils.MakeVerifier(hash[:], nowSalt)
	//转二进制
	nowVerifier := utils.FromBigSalt(verifier)
	//翻转数组
	utils.ReverseByteArray(nowVerifier)
	tx := db.Begin()

	// 执行原生 SQL 操作
	tx.Exec("INSERT INTO `account` (`username`, `salt`, `verifier`, `joindate`, `locale`, `last_ip`) VALUES (?, ?, ?, ?, ?, ?)",
		s.userName, string(nowSalt), string(nowVerifier), time.Now().Format(time.DateTime), 4, ip)
	if err := tx.Commit().Error; err != nil {
		// 处理提交错误
		tx.Rollback()
		utils.Logger.Error(fmt.Sprintf("[账号注册失败] 账号:%s,密码:%s,安全码:%s,ip地址:%s", s.userName, s.password, s.secureCode, ip))
		return errors.New("注册账号失败,请联系游戏管理员")
	} else {
		fmt.Printf("[账号注册] ip地址:%s,注册账号:%s", ip, s.userName)
	}

	userData := &model.UserInfo{
		UserName:   s.userName,
		Password:   s.password,
		SecureCode: s.secureCode,
	}

	db.Create(&userData)

	return errors.New("注册成功")
}

// 修改/找回密码前的检查
func (s *UserData) CheckPwd() error {
	db := utils.AuthDB
	if db == nil {
		return errors.New("连接数据库失败.请联系游戏管理员")
	}
	var userData = model.UserInfo{}
	var hasCount int64
	res := db.Model(&userData).
		Where("user_name=?", s.userName).First(&userData)
	res.Count(&hasCount)
	if hasCount <= 0 || res.Error != nil || (userData.SecureCode != s.secureCode) {
		return errors.New("安全码错误.找回/修改密码失败")
	}
	return nil
}

// 修改密码
func (s *UserData) ChangePassword() error {
	db := utils.AuthDB
	if db == nil {
		return errors.New("连接数据库失败.请联系游戏管理员")
	}
	if len(s.userName) > 12 || len(s.password) > 12 || len(s.secureCode) > 12 {
		return errors.New("账号或密码长度不能大于12位")
	}
	if err := s.CheckPwd(); err != nil {
		return err
	}

	// 先计算SHA1
	value := fmt.Sprintf("%s:%s", strings.ToUpper(s.userName), strings.ToUpper(s.password))
	hash := utils.ToHashSHA([]byte(value))
	// 生成Salt
	salt := utils.MakeSalt()
	//转二进制
	nowSalt := utils.FromBigSalt(salt)
	// 生成Verifier
	verifier := utils.MakeVerifier(hash[:], nowSalt)
	//转二进制
	nowVerifier := utils.FromBigSalt(verifier)
	//翻转数组
	utils.ReverseByteArray(nowVerifier)

	tx := db.Begin()

	// 执行原生 SQL 操作
	tx.Exec("UPDATE `account` SET `salt`=?, `verifier`=? WHERE username=?", string(nowSalt), string(nowVerifier), s.userName)

	if err := tx.Commit().Error; err != nil {
		// 处理提交错误
		tx.Rollback()
		utils.Logger.Error(fmt.Sprintf("[找回/修改密码失败] 账号:%s,密码:%s", s.userName, s.password))
		return errors.New("找回/修改密码失败,请联系游戏管理员")
	}

	db.Model(&model.UserInfo{}).
		Where("user_name=?", s.userName).
		First(&model.UserInfo{}).
		Select("password").
		Update("password", s.password)

	return errors.New("找回/修改密码成功")
}

// 根据用户名和密码验证账号有效性
func (s *UserData) CheckSaltAndVerifier() (uint, error) {
	db := utils.AuthDB
	if db == nil {
		return 0, errors.New("连接数据库失败.请联系游戏管理员")
	}

	var id uint
	var salt []byte
	var verifier []byte
	//db.Raw("SELECT id, salt, verifier FROM account WHERE username=? ", s.userName).Scan(&check)
	row := db.Raw("SELECT id, salt, verifier FROM account WHERE username=? ", s.userName).Row()
	err := row.Scan(&id, &salt, &verifier)
	if err != nil {
		return 0, errors.New("账号或密码错误")
	}
	if !utils.CheckSaltVerifier(s.userName, s.password, salt, verifier) {
		return 0, errors.New("账号或密码错误")
	}

	return id, nil
}
