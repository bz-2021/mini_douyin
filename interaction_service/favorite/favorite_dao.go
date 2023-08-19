package favorite_service

import (
	"context"
	"fmt"
	pb "github.com/bz-2021/mini_douyin/interaction_service/favorite/favorite_grpc"
	"github.com/bz-2021/mini_douyin/interaction_service/favorite/pojo"
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

// 增加点赞记录
func (fs *FavoriteService) AddFavoriteRecord(ctx context.Context, req *pb.FavoriteActionRequest, userId int64) (favorite *pojo.Favorite, err error) {
	fmt.Println("\nadd one single favorite record...\n")
	//construct
	favorite, err = fs.NewFavoriteItem(ctx, req, userId)
	if err != nil {
		return nil, err
	}
	fmt.Println("点赞记录构建成功")
	//检查记录是否已存在
	db := fs.DB.WithContext(ctx)
	db = db.Where("user_id = ? AND video_id = ?", favorite.UserId, favorite.VideoId).Order("video_id").Find(&favorite)
	if db.Error != nil {
		fmt.Errorf("select false: %s", db.Error.Error())
		return nil, db.Error
	}
	//点赞记录已存在，基于用户体验，暂不做任何处理
	if db.RowsAffected != 0 {
		//no handle
		return favorite, nil
	}
	//记录不存在
	tx := fs.DB.WithContext(ctx).Create(&favorite)
	if tx.Error != nil {
		fmt.Errorf("create false: %s", tx.Error.Error())
		return nil, tx.Error
	}
	return favorite, nil
}

// 删除点赞记录
func (fs *FavoriteService) DeleteFavoriteRecord(ctx context.Context, req *pb.FavoriteActionRequest, userId int64) (*pojo.Favorite, error) {
	fmt.Println("\ndelete one single favorite record...\n")
	//construct
	favorite, err := fs.NewFavoriteItem(ctx, req, userId)
	if err != nil {
		return nil, err
	}
	db := fs.DB.WithContext(ctx)
	db = db.Where("user_id = ? AND video_id = ?", favorite.UserId, favorite.VideoId).First(&favorite)
	db.Delete(&favorite)
	//错误处理 - 点赞记录不存在
	if db.RowsAffected == 0 {
		fmt.Errorf("delete false: %s", db.Error.Error())
		return nil, db.Error
	}
	return nil, nil
}

// 构建一个点赞记录
func (fs *FavoriteService) NewFavoriteItem(ctx context.Context, req *pb.FavoriteActionRequest, userId int64) (*pojo.Favorite, error) {
	favorite := &pojo.Favorite{}
	favorite.VideoId = req.VideoId
	favorite.UserId = userId
	return favorite, nil
}
