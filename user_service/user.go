package user_service

import (
	"context"
	"fmt"
	pb "github.com/bz-2021/mini_douyin/user_service/user_grpc"
	"github.com/bz-2021/mini_douyin/utils"
	"gorm.io/gorm"
	"strconv"
)

type UserLoginService struct {
	pb.UnimplementedServiceServer
	DB *gorm.DB
}

func NewUserLoginService() *UserLoginService {
	db, err := GetDB()
	if err != nil {
		panic("NewUserLoginService失败")
	}
	return &UserLoginService{
		DB: db,
	}
}

func (u *UserLoginService) Login(ctx context.Context, req *pb.UserLoginRequest) (resp *pb.UserLoginResponse, err error) {
	fmt.Println("微服务调用成功，开始查询")
	fmt.Printf("username : %v ", req.Username)
	fmt.Printf("password : %v ", req.Password)

	username := req.Username
	password := req.Password

	resp = new(pb.UserLoginResponse)

	// 参数为空，请求失败
	if len(username) == 0 || len(password) == 0 {
		resp.StatusCode = 400
		resp.StatusMsg = &utils.BadRequest
		return
	}

	//用户登录验证逻辑
	user, err := u.getUserByUsername(ctx, username)
	fmt.Println(user)
	if err != nil || len(user.Name) == 0 {
		resp.StatusCode = 403
		resp.StatusMsg = &utils.WrongUsernameOrPassword
		return
	}
	token, err := utils.GenerateJWT(strconv.FormatInt(user.Id, 10))

	// 密码错误，登陆失败
	if !utils.ComparePasswords(user.Password, password) {
		fmt.Println(user.Password, password)
		resp.StatusCode = 403
		resp.StatusMsg = &utils.WrongUsernameOrPassword
		return
	}

	// 请求成功
	resp = &pb.UserLoginResponse{
		StatusCode: 0,
		StatusMsg:  &utils.Succeed,
		UserId:     user.Id,
		Token:      token,
	}
	return
}

func (u *UserLoginService) Register(ctx context.Context, req *pb.UserRegisterRequest) (resp *pb.UserRegisterResponse, err error) {
	username := req.Username
	password := req.Password

	//token, err := utils.GenerateJWT(username)
	resp = new(pb.UserRegisterResponse)

	// 参数为空，请求失败
	if len(username) == 0 || len(password) == 0 {
		resp.StatusCode = 400
		resp.StatusMsg = &utils.BadRequest
		return
	}

	user, err := u.getUserByUsername(ctx, username)
	if len(user.Name) != 0 || err != nil {
		resp.StatusCode = 1
		resp.StatusMsg = &utils.UsernameAlreadyExist
		fmt.Println(err)
		return
	}

	err = u.insertUser(ctx, username, password)
	if err != nil {
		resp.StatusCode = 1
		resp.StatusMsg = &utils.InternalServerErr
		return
	}
	maxID, err := u.maxId(ctx)
	if err != nil {
		resp.StatusCode = 1
		resp.StatusMsg = &utils.InternalServerErr
		return
	}
	token, err := utils.GenerateJWT(strconv.FormatInt(maxID, 10))
	if err != nil {
		resp.StatusCode = 1
		resp.StatusMsg = &utils.InternalServerErr
		return
	}
	resp.Token = token
	resp.UserId = maxID
	resp.StatusCode = 0
	resp.StatusMsg = &utils.Succeed
	return
}

func (u *UserLoginService) UserInfo(ctx context.Context, req *pb.UserInfoRequest) (resp *pb.UserInfoResponse, err error) {
	token := req.Token
	userId := req.UserId

	myStringId, err := utils.VerifyJWT(token)
	resp = new(pb.UserInfoResponse)

	if err != nil {
		resp.StatusCode = 1
		resp.StatusMsg = &utils.PermissionDenied
		return
	}
	myId, err := strconv.ParseInt(myStringId, 10, 64)
	if err != nil {
		return nil, err
	}

	user, err := u.getUserById(ctx, myId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		resp.StatusCode = 1
		resp.StatusMsg = &utils.PermissionDenied
		return
	}
	thisUser, err := u.getUserById(ctx, userId)
	if err != nil {
		return nil, err
	}
	resp.User = &pb.User{
		Id:        thisUser.Id,
		Name:      thisUser.Name,
		Avatar:    &thisUser.Avatar,
		Signature: &thisUser.Signature,
	}
	resp.StatusCode = 0
	resp.StatusMsg = &utils.Succeed
	return

}
