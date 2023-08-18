package feed_service

import (
	"context"
	service "github.com/bz-2021/mini_douyin/feed_service/feed_grpc"
	"github.com/bz-2021/mini_douyin/feed_service/feed_grpc/video"
	"gorm.io/gorm"
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
	return
}