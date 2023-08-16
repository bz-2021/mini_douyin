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

func (u *UserLoginService) getUserById(ctx context.Context, id int64) (user *pojo.User, err error) {
	user = &pojo.User{}
	db := u.DB.WithContext(ctx)
	db = db.Table("user")
	db = db.Where("id = ?", id).Find(user)
	if db.Error != nil {
		return nil, db.Error
	}
	return
}

func (u *UserLoginService) insertUser(ctx context.Context, username string, password string) error {
	// 向表中插入数据
	encodedPassword, err := utils.HashAndSalt(password)
	if err != nil {
		return err
	}

	db := u.DB.WithContext(ctx)
	db = db.Table("user")
	db.Create(&pojo.User{
		Name:     username,
		Password: encodedPassword,
	})

	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (u *UserLoginService) maxId(ctx context.Context) (int64, error) {
	db := u.DB.WithContext(ctx)
	var maxID int64
	result := db.Table("user").Select("MAX(id)").Scan(&maxID)
	if result.Error != nil {
		return 0, result.Error
	}
	return maxID, nil
}
