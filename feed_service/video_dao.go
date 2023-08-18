package feed_service

import (
	"context"
	"fmt"
	"github.com/bz-2021/mini_douyin/feed_service/pojo"
	user_pojo "github.com/bz-2021/mini_douyin/user_service/pojo"
	"github.com/bz-2021/mini_douyin/utils"
	"gorm.io/gorm"
)

var GormDB *gorm.DB

var MySQLDatabase *utils.Mysql

func GetDB() (*gorm.DB, error) {
	if MySQLDatabase == nil {
		MySQLDatabase = utils.DefaultMySQLDB()
	}
	return MySQLDatabase.GetDB(GormDB)
}

func (f *FeedService) getVideoListByDate(ctx context.Context, date int64) (list []pojo.Video, err error) {
	t := make([]pojo.Video, 10)
	db := f.DB.WithContext(ctx)
	db = db.Table("video").Find(&t)
	if db.Error != nil {
		return nil, db.Error
	}
	fmt.Println("list在这咯", t)
	return t, nil
}

func (u *FeedService) getUserById(ctx context.Context, id int64) (user *user_pojo.User, err error) {
	user = &user_pojo.User{}
	db := u.DB.WithContext(ctx)
	db = db.Table("user")
	db = db.Where("id = ?", id).Find(user)
	if db.Error != nil {
		return nil, db.Error
	}
	return
}
