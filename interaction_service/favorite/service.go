package favorite_service

import (
	"context"
	"fmt"
	pb "github.com/bz-2021/mini_douyin/interaction_service/favorite/favorite_grpc"
	user_pojo "github.com/bz-2021/mini_douyin/user_service/pojo"
	"github.com/bz-2021/mini_douyin/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"strconv"
)

type FavoriteService struct {
	pb.UnimplementedServiceServer
	DB *gorm.DB
}

func NewFavoriteService() *FavoriteService {
	db, err := GetDB()
	if err != nil {
		panic("NewFavoriteService失败")
	}
	return &FavoriteService{
		DB: db,
	}
}

func (fs FavoriteService) FavoriteAction(ctx context.Context, req *pb.FavoriteActionRequest) (resp *pb.FavoriteActionResponse, err error) {
	fmt.Println("微服务-FavoriteAction 调用成功，开始查询")

	//todo: 入参req校验

	token := req.Token
	actionType := req.ActionType

	resp = new(pb.FavoriteActionResponse)

	//parse token & get userId
	myStringId, err := utils.VerifyJWT(token)
	myId, err := strconv.ParseInt(myStringId, 10, 64)
	if err != nil {
		return nil, err
	}
	user, err := fs.getUserById(ctx, myId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		resp.StatusCode = 1
		resp.StatusMsg = utils.PermissionDenied
		return
	}

	//parse action
	switch actionType {
	//action_type: 1 - 点赞 ; 2 - 取消点赞
	case 1:
		fmt.Println("处理点赞操作")
		//点赞操作
		_, err := fs.AddFavoriteRecord(ctx, req, user.Id)
		//失败 - 错误处理
		if err != nil {
			fmt.Println("点赞失败，发生错误")
		}
		//成功 - 返回响应
		return resp, nil
	case 2:
		//取消点赞操作
		_, err := fs.DeleteFavoriteRecord(ctx, req, user.Id)
		//失败 - 错误处理
		if err != nil {
			fmt.Println("取消点赞失败，发生错误")
		}
		//成功 - 返回响应
		return resp, nil
	default:
		//error: action_type != 1 && action_type != 2
		fmt.Println("Incorrect action type")
	}

	// 请求成功
	resp = &pb.FavoriteActionResponse{
		StatusCode: 0,
		StatusMsg:  utils.Succeed,
	}
	return resp, nil
}

func (fs FavoriteService) FavoriteList(ctx context.Context, req *pb.FavoriteListRequest) (
	*pb.FavoriteListResponse, error) {
	fmt.Println("rpc方法FavoriteList调用成功，开始处理")
	item := &pb.FavoriteItem{}
	item.UserId = req.UserId
	//	获取喜欢视频列表
	pos, err := fs.getFavoriteListPo(ctx, item)
	if err != nil {
		fmt.Errorf("favorite list get false - getFavoriteListPo() Error: \n %s", err.Error())
		return nil, status.Errorf(codes.Unavailable, "[Error] Failed to Get Favorite List")
	}

	resp := &pb.FavoriteListResponse{}

	// 根据列表分别获取用户信息
	videoList := make([]*pb.Video, len(pos))
	//each item in pos is a pb.FavoriteItem
	resp.VideoList = videoList

	if len(videoList) == 0 {
		return resp, nil
	}

	// 将用户信息和视频信息组合成Response
	//each item in pos is a pb.FavoriteItem
	//po record : Id  UserId  VideoId
	for i, po := range pos {
		videoRespPo, vErr := fs.GetVideoResp(ctx, po.VideoId)
		if vErr != nil {
			fmt.Errorf(vErr.Error())
			return nil, vErr
		}
		resp.VideoList[i] = videoRespPo
	}

	// 请求成功
	resp.StatusCode = 0
	resp.StatusMsg = utils.Succeed
	return resp, nil
}

func (u *FavoriteService) getUserById(ctx context.Context, id int64) (user *user_pojo.User, err error) {
	user = &user_pojo.User{}
	db := u.DB.WithContext(ctx)
	db = db.Table("user")
	db = db.Where("id = ?", id).Find(user)
	if db.Error != nil {
		return nil, db.Error
	}
	return
}
