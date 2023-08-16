package dao

import (
	"github.com/bz-2021/mini_douyin/utils"
	"gorm.io/gorm"
)

// 有些混乱（？），为了复用 utils 中的代码

var GormDB *gorm.DB

var MySQLDatabase *utils.Mysql

func GetDB() (*gorm.DB, error) {
	if MySQLDatabase == nil {
		MySQLDatabase = utils.DefaultMySQLDB()
	}
	return MySQLDatabase.GetDB(GormDB)
}
