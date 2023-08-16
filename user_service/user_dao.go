package user_service

import (
	"context"
	"github.com/bz-2021/mini_douyin/user_service/pojo"
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

func (u *UserLoginService) getUserByUsername(ctx context.Context, username string) (user *pojo.User, err error) {
	user = &pojo.User{}
	db := u.DB.WithContext(ctx)
	db = db.Table("user")
	db = db.Where("name = ?", username).Find(user)
	if db.Error != nil {
		return nil, db.Error
	}
	return
}
