package comment_service

import (
	"context"
	"fmt"
	pb "github.com/bz-2021/mini_douyin/interaction_service/comment/comment_grpc"
	user_pojo "github.com/bz-2021/mini_douyin/user_service/pojo"
	"github.com/bz-2021/mini_douyin/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"strconv"
)

type CommentService struct {
	pb.UnimplementedServiceServer
	DB *gorm.DB
}

func NewCommentService() *CommentService {
	db, err := GetDB()
	if err != nil {
		panic("NewCommentService失败")
	}
	return &CommentService{
		DB: db,
	}
}

func (fc *CommentService) CommentAction(ctx context.Context, req *pb.CommentActionRequest) (resp *pb.CommentActionResponse, err error) {
	//todo: 入参req校验
	token := req.Token
	actionType := req.ActionType

	resp = new(pb.CommentActionResponse)

	//parse token & get userId
	myStringId, err := utils.VerifyJWT(token)
	myId, err := strconv.ParseInt(myStringId, 10, 64)
	if err != nil {
		return nil, err
	}
	user, err := fc.getUserById(ctx, myId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		resp.StatusCode = 1
		resp.StatusMsg = &utils.PermissionDenied
		return
	}

	//parse action
	switch actionType {
	//action_type: 1 - 点赞 ; 2 - 取消点赞
	case 1:
		fmt.Println("新增评论操作")
		//评论操作
		cmtResp, err := fc.AddCommentRecord(ctx, req, user.Id)
		//失败 - 错误处理
		if err != nil {
			fmt.Println("评论失败，发生错误")
		}
		//成功 - 返回响应
		resp = cmtResp
		//return resp, nil
	case 2:
		//删除评论操作
		cmtResp, err := fc.DeleteCommentRecord(ctx, req, user.Id)
		//失败 - 错误处理
		if err != nil {
			fmt.Println("删除评论失败，发生错误")
		}
		//成功 - 返回响应
		resp = cmtResp
		//return resp, nil
	default:
		//error: action_type != 1 && action_type != 2
		fmt.Errorf("Incorrect action type")
		return nil, nil
	}

	// 请求成功
	resp.StatusCode = 0
	resp.StatusMsg = &utils.Succeed
	return resp, nil
}

func (fc *CommentService) CommentList(ctx context.Context, req *pb.CommentListRequest) (
	*pb.CommentListResponse, error) {
	fmt.Println("rpc方法CommentList调用成功，开始处理")
	item := &pb.Comment{}
	item.VideoId = req.VideoId
	//	获取喜欢视频列表
	pos, err := fc.getCommentListPo(ctx, item)
	if err != nil {
		fmt.Errorf("favorite list get false - getFavoriteListPo() Error: \n %s", err.Error())
		return nil, status.Errorf(codes.Unavailable, "[Error] Failed to Get Favorite List")
	}

	resp := &pb.CommentListResponse{}

	// 根据列表分别获取用户信息
	commentList := make([]*pb.Comment, len(pos))
	//each item in pos is a pb.FavoriteItem
	resp.CommentList = commentList

	if len(commentList) == 0 {
		return resp, nil
	}

	// 将用户信息和视频信息组合成Response
	//each item in pos is a pb.Comment
	//po record : Id  UserId  VideoId
	for i, po := range pos {
		commentPo, vErr := fc.GetComment(ctx, po)
		if vErr != nil {
			fmt.Errorf(vErr.Error())
			return nil, vErr
		}
		resp.CommentList[i] = commentPo
	}

	// 请求成功
	resp.StatusCode = 0
	resp.StatusMsg = utils.Succeed
	return resp, nil
}

func (u *CommentService) getUserById(ctx context.Context, id int64) (user *user_pojo.User, err error) {
	user = &user_pojo.User{}
	db := u.DB.WithContext(ctx)
	db = db.Table("user")
	db = db.Where("id = ?", id).Find(user)
	if db.Error != nil {
		return nil, db.Error
	}
	return
}
