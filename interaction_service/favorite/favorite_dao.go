package favorite_service

import (
	"context"
	"fmt"
	pojo2 "github.com/bz-2021/mini_douyin/feed_service/pojo"
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
	//1. 新增record
	tx := fs.DB.WithContext(ctx).Create(&favorite)
	if tx.Error != nil {
		fmt.Errorf("create false: %s", tx.Error.Error())
		return nil, tx.Error
	}
	//以下应当使用一组事务
	//2. 增加users表中当前user的favorite_count计数 and 当前视频author的total_favorite计数 && 增加video表的favorite_count计数
	//2-1 user favorite_count
	db_f := fs.DB.WithContext(ctx)
	exprAddFavoriteCount := gorm.Expr("favorite_count + ?", 1)
	db_f = db_f.Table("user").Where("id = ?", userId).Update("favorite_count", exprAddFavoriteCount)
	if db_f.Error != nil {
		fmt.Errorf("update false: %s", db.Error.Error())
		return nil, db.Error
	}
	//2-3 video favorite_count
	videoItem := &pojo2.Video{}
	db_v := fs.DB.WithContext(ctx)
	video_result := db_v.Table("video").Where("id = ?", req.VideoId).First(&videoItem)
	if video_result.Error != nil {
		fmt.Errorf("new videoItem false: %s", video_result.Error.Error())
		return nil, video_result.Error
	}
	authorID := videoItem.UserId
	db_v = fs.DB.WithContext(ctx)
	exprAddVideoFavoriteCount := gorm.Expr("favorite_count + ?", 1)
	db_v = db_v.Table("video").Where("id = ?", req.VideoId).Update("favorite_count", exprAddVideoFavoriteCount)
	if db_v.Error != nil {
		fmt.Errorf("update false: %s", db.Error.Error())
		return nil, db.Error
	}
	//2-2 author total_favorite
	db_au := fs.DB.WithContext(ctx)
	exprAddTotalFavorite := gorm.Expr("total_favorite + ?", 1)
	db_au = db_au.Table("user").Where("id = ?", authorID).Update("total_favorite", exprAddTotalFavorite)
	if db_au.Error != nil {
		fmt.Errorf("update false: %s", db.Error.Error())
		return nil, db.Error
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

	//检查记录是否已存在
	db := fs.DB.WithContext(ctx)
	db = db.Where("user_id = ? AND video_id = ?", favorite.UserId, favorite.VideoId).First(&favorite)
	db.Delete(&favorite)

	//1. 错误处理 - 点赞记录不存在
	if db.RowsAffected == 0 {
		fmt.Errorf("delete false: %s", db.Error.Error())
		return nil, db.Error
	}
	//2. 记录存在并删除
	//以下应当使用一组事务
	//2. 减少users表中当前user的favorite_count计数 and 当前视频author的total_favorite计数 && 减少video表的favorite_count计数
	//2-1 user favorite_count
	db = fs.DB.WithContext(ctx)
	exprAddFavoriteCount := gorm.Expr("favorite_count - ?", 1)
	db = db.Table("user").Where("id = ?", userId).Update("favorite_count", exprAddFavoriteCount)
	if db.Error != nil {
		fmt.Errorf("update false: %s", db.Error.Error())
		return nil, db.Error
	}
	//2-3 video favorite_count
	videoItem := &pojo2.Video{}
	db_v := fs.DB.WithContext(ctx)
	video_result := db_v.Table("video").Where("id = ?", req.VideoId).First(&videoItem)
	if video_result.Error != nil {
		fmt.Errorf("new videoItem false: %s", video_result.Error.Error())
		return nil, video_result.Error
	}
	authorID := videoItem.UserId
	db_v = fs.DB.WithContext(ctx)
	exprAddVideoFavoriteCount := gorm.Expr("favorite_count - ?", 1)
	db_v = db_v.Table("video").Where("id = ?", req.VideoId).Update("favorite_count", exprAddVideoFavoriteCount)
	if db_v.Error != nil {
		fmt.Errorf("update false: %s", db.Error.Error())
		return nil, db.Error
	}
	//2-2 author total_favorite
	db_au := fs.DB.WithContext(ctx)
	exprAddTotalFavorite := gorm.Expr("total_favorite - ?", 1)
	db_au = db_au.Table("user").Where("id = ?", authorID).Update("total_favorite", exprAddTotalFavorite)
	if db_au.Error != nil {
		fmt.Errorf("update false: %s", db.Error.Error())
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
