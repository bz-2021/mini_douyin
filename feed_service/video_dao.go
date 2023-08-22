package feed_service

import (
	"context"
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

func (f *FeedService) getVideoListByDate(ctx context.Context, date int64, userID int64) (list []pojo.Video, err error) {
	videos := make([]pojo.Video, 10)
	db := f.DB.WithContext(ctx)
	db = db.Table("video").Preload("Author").Find(&videos)
	if db.Error != nil {
		return nil, db.Error
	}

	var favorites []pojo.Favorite
	db = f.DB.WithContext(ctx)
	err = db.Table("favorite").Where("user_id = ?", userID).Find(&favorites).Error
	if err != nil {
		return nil, err
	}

	// 设置is_favorite字段
	favoriteIDs := make(map[int64]bool)
	for _, favorite := range favorites {
		favoriteIDs[favorite.VideoId] = true
	}
	for i := range videos {
		if _, ok := favoriteIDs[videos[i].Id]; ok {
			videos[i].IsFavorite = true
		}
	}

	return videos, nil
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
