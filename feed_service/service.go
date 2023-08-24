package feed_service

import (
	"context"
	"fmt"
	service "github.com/bz-2021/mini_douyin/feed_service/feed_grpc"
	"github.com/bz-2021/mini_douyin/feed_service/feed_grpc/video"
	"github.com/bz-2021/mini_douyin/utils"
	"gorm.io/gorm"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type FeedService struct {
	video.UnimplementedServiceServer

	DB *gorm.DB
}

func NewFeedService() *FeedService {
	db, err := GetDB()
	if err != nil {
		panic("NewFeedService失败")
	}
	return &FeedService{
		DB: db,
	}
}

func (f *FeedService) PublishAction(ctx context.Context, req *service.PublishActionRequest) (resp *service.PublishActionResponse, err error) {
	// 获取参数
	token := req.Token
	data := req.Data
	title := req.Title

	//鉴权
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

	//处理data
	tempFile, err := os.CreateTemp("", "video")
	if err != nil {
		fmt.Println("Error creating temp file:", err)
		return
	}
	defer func(tempFile *os.File) {
		err := tempFile.Close()
		if err != nil {

		}
	}(tempFile)
	_, err = tempFile.Write(data)
	if err != nil {
		fmt.Println("Error writing to temp file:", err)
		return
	}

	outputFilePath := "./public/"
	cmd := exec.Command("ffmpeg", "-i", tempFile.Name(), outputFilePath)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error executing ffmpeg:", err)
		return
	}
	fmt.Println("Video conversion successful")

	playURL := outputFilePath + tempFile.Name()
	coverURL := outputFilePath + "cover.jpg"

	cmd = exec.Command("ffmpeg", "-i", playURL, "-ss", "00:00:00.001", "-vframes", "1", coverURL)

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Image extracted successfully")

	f.insertVideo(ctx, myStringId, title, playURL, coverURL, time.Now().Format("2006-01-02 15:04:05"))

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
	list, err := f.getVideoListByDate(ctx, lastTime, myId)
	if err != nil {
		return
	}
	videoList := make([]*service.Video, len(list))
	for i, v := range list {
		videoList[i] = &service.Video{
			Id: v.Id,
			Author: &service.User{
				Id:              v.Author.ID,
				Avatar:          &list[i].Author.Avatar,
				Name:            v.Author.Name,
				FollowCount:     v.Author.FollowCount,
				FollowerCount:   v.Author.FollowerCount,
				IsFollow:        v.Author.IsFollow,
				BackgroundImage: &list[i].Author.BackgroundImage,
				Signature:       &list[i].Author.Signature,
				TotalFavorited:  v.Author.TotalFavorite,
				WorkCount:       v.Author.WorkCount,
				FavoriteCount:   v.Author.FavoriteCount,
			},
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    v.IsFavorite,
			Title:         v.Title,
		}
	}
	return &service.FeedResponse{
		StatusCode: 0,
		StatusMsg:  &utils.Succeed,
		VideoList:  videoList,
		NextTime:   0,
	}, nil
}
