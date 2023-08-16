package user_service

import (
	"context"
	"fmt"
	"github.com/bz-2021/mini_douyin/user_service/pojo"
	pb "github.com/bz-2021/mini_douyin/user_service/user_grpc"
	"github.com/bz-2021/mini_douyin/utils"
	"gorm.io/gorm"
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
	token, err := utils.GenerateJWT(username)

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
	if err != nil || user == nil {
		resp.StatusCode = 403
		resp.StatusMsg = &utils.WrongUsernameOrPassword
		return
	}

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

	// 向表中插入数据
	encodedPassword, err := utils.HashAndSalt(password)
	if err != nil {
		panic("密码加密失败")
	}

	db := u.DB.WithContext(ctx)
	db = db.Table("user")
	db.Create(&pojo.User{
		Name:     username,
		Password: encodedPassword,
	})

	if db.Error != nil {

	}
	return
}

func (u *UserLoginService) UserInfo(ctx context.Context, req *pb.UserInfoRequest) (*pb.UserInfoResponse, error) {
	return nil, nil
}
