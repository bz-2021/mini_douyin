package feed_service

import (
	"context"
	"fmt"
	service "github.com/bz-2021/mini_douyin/feed_service/feed_grpc"
	"github.com/bz-2021/mini_douyin/feed_service/feed_grpc/video"
	"github.com/bz-2021/mini_douyin/utils"
	"gorm.io/gorm"
	"strconv"
)

type FeedService struct {
	video.UnimplementedServiceServer
	DB *gorm.DB
}

func NewFeedService() *FeedService {
	db, err := GetDB()
	if err != nil {
		panic("NewUserLoginService失败")
	}
	return &FeedService{
		DB: db,
	}
}

func (f *FeedService) PublishAction(ctx context.Context, req *service.PublishActionRequest) (resp *service.PublishActionResponse, err error) {
	return
}

func (f *FeedService) PublishList(ctx context.Context, req *service.PublishListRequest) (resp *service.PublishListResponse, err error) {
	return
}

func (f *FeedService) FeedAction(ctx context.Context, req *service.FeedRequest) (resp *service.FeedResponse, err error) {
	lastTime := req.LatestTime
	token := req.Token

	myStringId, err := utils.VerifyJWT(token)
	myId, err := strconv.ParseInt(myStringId, 10, 64)
	if err != nil {
		return nil, err
	}

	user, err := f.getUserById(ctx, myId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		resp.StatusCode = 1
		resp.StatusMsg = &utils.PermissionDenied
		return
	}
	list, err := f.getVideoListByDate(ctx, lastTime)
	if err != nil {
		return
	}
	videoList := make([]*service.Video, len(list))
	avatar := "https://cdn.acwing.com/media/user/profile/photo/220156_lg_9ddb2ec392.jpg"
	for i, v := range list {
		videoList[i] = &service.Video{
			Id: v.Id,
			Author: &service.User{
				Id:     1,
				Avatar: &avatar,
			},
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: int64(v.FavoriteCount),
			CommentCount:  int64(v.CommentCount),
			IsFavorite:    false,
			Title:         v.Title,
		}
	}
	fmt.Println("videoList在这里", videoList)
	return &service.FeedResponse{
		StatusCode: 0,
		StatusMsg:  &utils.Succeed,
		VideoList:  videoList,
		NextTime:   0,
	}, nil
}
