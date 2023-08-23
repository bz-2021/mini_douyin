package comment_service

import (
	"context"
	"fmt"
	pojo2 "github.com/bz-2021/mini_douyin/feed_service/pojo"
	pb "github.com/bz-2021/mini_douyin/interaction_service/comment/comment_grpc"
	"github.com/bz-2021/mini_douyin/interaction_service/comment/pojo"
	"github.com/bz-2021/mini_douyin/utils"
	"gorm.io/gorm"
	"time"
)

var GormDB *gorm.DB

var MySQLDatabase *utils.Mysql

func GetDB() (*gorm.DB, error) {
	if MySQLDatabase == nil {
		MySQLDatabase = utils.DefaultMySQLDB()
	}
	return MySQLDatabase.GetDB(GormDB)
}

// 增加评论记录
func (fc *CommentService) AddCommentRecord(ctx context.Context, req *pb.CommentActionRequest, userId int64) (comment *pb.CommentActionResponse, err error) {
	fmt.Println("\nadd one single comment record...\n")
	//construct
	commentItem, err := fc.NewCommentItem(ctx, req, userId)
	if err != nil {
		return nil, err
	}
	fmt.Println("评论记录构建成功")
	//检查记录是否已存在
	db := fc.DB.WithContext(ctx)
	result := db.Create(&commentItem)
	if result.Error != nil {
		fmt.Errorf("comment record create failed")
	}

	//
	//2. 增加当前视频的comment_count计数
	// video comment_count
	videoItem := &pojo2.Video{}
	db_v := fc.DB.WithContext(ctx)
	video_result := db_v.Table("video").Where("id = ?", req.VideoId).First(&videoItem)
	if video_result.Error != nil {
		fmt.Errorf("new videoItem false: %s", video_result.Error.Error())
		return nil, video_result.Error
	}
	db_v = fc.DB.WithContext(ctx)
	exprAddVideoCommentCount := gorm.Expr("comment_count + ?", 1)
	db_v = db_v.Table("video").Where("id = ?", req.VideoId).Update("comment_count", exprAddVideoCommentCount)
	if db_v.Error != nil {
		fmt.Errorf("update false: %s", db.Error.Error())
		return nil, db.Error
	}
	comment, err = fc.GetCommentResp(ctx, commentItem)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

// 删除评论记录
func (fc *CommentService) DeleteCommentRecord(ctx context.Context, req *pb.CommentActionRequest, userId int64) (comment *pb.CommentActionResponse, err error) {
	fmt.Println("\ndelete one single favorite record...\n")
	//construct
	commentItem, err := fc.NewCommentItem(ctx, req, userId)
	if err != nil {
		return nil, err
	}

	//检查评论记录是否存在
	db := fc.DB.WithContext(ctx)
	db = db.Where("id = ?", commentItem.Id).First(&commentItem)
	db.Delete(&commentItem)

	//1. 错误处理 - 点赞记录不存在
	if db.RowsAffected == 0 {
		fmt.Errorf("delete false: %s", db.Error.Error())
		return nil, db.Error
	}
	//2. 记录存在并删除
	//2. 减少当前视频的comment_count计数
	// video comment_count
	videoItem := &pojo2.Video{}
	db_v := fc.DB.WithContext(ctx)
	video_result := db_v.Table("video").Where("id = ?", req.VideoId).First(&videoItem)
	if video_result.Error != nil {
		fmt.Errorf("new videoItem false: %s", video_result.Error.Error())
		return nil, video_result.Error
	}
	db_v = fc.DB.WithContext(ctx)
	exprAddVideoCommentCount := gorm.Expr("comment_count - ?", 1)
	db_v = db_v.Table("video").Where("id = ?", req.VideoId).Update("comment_count", exprAddVideoCommentCount)
	if db_v.Error != nil {
		fmt.Errorf("update false: %s", db.Error.Error())
		return nil, db.Error
	}
	comment, err = fc.GetCommentResp(ctx, commentItem)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

// 获取评论列表
func (fc *CommentService) getCommentListPo(ctx context.Context, po *pb.Comment) (
	[]*pojo.Comment, error) {

	//向数据库查询所有数据
	db := fc.DB.WithContext(ctx)
	pos := make([]*pojo.Comment, 0)
	if po.VideoId > 0 {
		db = db.Table("comment")
		db = db.Where("video_id = ?", po.VideoId)
	} else {
		fmt.Println("get list false: Parameter false")
		return pos, nil
	}
	//在comment表中查找当前视频所有的评论记录并进行绑定
	db.Order("create_date DESC").Find(&pos) //record : Id  UserId  VideoId
	if db.Error != nil {
		fmt.Errorf("select comment records false: %s", db.Error.Error())
		return nil, db.Error
	}
	return pos, nil
}

// 构建一个评论记录
func (fc *CommentService) NewCommentItem(ctx context.Context, req *pb.CommentActionRequest, userId int64) (*pojo.Comment, error) {
	comment := &pojo.Comment{}
	comment.VideoId = req.VideoId
	comment.UserId = userId
	comment.Id = req.CommentId
	comment.Content = req.CommentText
	comment.CreateDate = time.Now().Format("2006-01-02 15:04:05")
	//time.Now().Format("2023-08-22 15:04:05")
	return comment, nil
}

// 构建一个评论响应 - 将用户信息与评论组合
func (fc *CommentService) GetCommentResp(ctx context.Context, req *pojo.Comment) (resp *pb.CommentActionResponse, err error) {
	resp = &pb.CommentActionResponse{}
	//comment
	comment := &pb.Comment{}
	comment.Id = req.Id
	user := pb.User{}
	db := fc.DB.WithContext(ctx)
	db = db.Table("user").Where("id = ?", req.UserId).First(&user)
	comment.User = &user
	comment.Content = req.Content
	comment.CreateDate = req.CreateDate

	resp.Comment = comment
	return resp, nil
}

// 构建一个评论 - 将用户信息与评论组合，与GetCommentResp的区别在于返回类型
func (fc *CommentService) GetComment(ctx context.Context, req *pojo.Comment) (resp *pb.Comment, err error) {
	resp = &pb.Comment{}
	//comment
	resp.Id = req.Id
	user := pb.User{}
	db := fc.DB.WithContext(ctx)
	db = db.Table("user").Where("id = ?", req.UserId).First(&user)
	resp.User = &user
	resp.Content = req.Content
	resp.CreateDate = req.CreateDate
	return resp, nil
}
