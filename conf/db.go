package conf

import (
	"gitee.com/mrmateoliu/wow_launch.git/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

func InitDB() (*gorm.DB, error) {

	logMode := logger.Info
	if !utils.GlobalObject.Develop {
		logMode = logger.Error
	}
	db, err := gorm.Open(mysql.Open(utils.GlobalObject.DbDsn1), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "sys_",
			SingularTable: true,
			//NameReplacer:        nil,
			//NoLowerCase:         false,
			//IdentifierMaxLength: 0,
		},
		Logger: logger.Default.LogMode(logMode),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(utils.GlobalObject.DbmaxIdleConn) // 最多空闲链接数
	sqlDB.SetMaxOpenConns(utils.GlobalObject.DbmaxOpenConn) // 最多打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour * 24)                //生命周期为1天

	err = db.AutoMigrate(
	//&model.UserBaseData{},
	)

	return db, nil
}
