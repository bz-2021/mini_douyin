package main

import (
	"context"
	service "github.com/bz-2021/mini_douyin/user_account_service/protoc_gen"
	"google.golang.org/grpc"
	"log"
	"net"
)

type UserAccountServer struct {
	service.UnimplementedServiceServer
}

func (s UserAccountServer) Register(ctx context.Context, request *service.UserRegisterRequest) (*service.UserRegisterResponse, error) {
	return &service.UserRegisterResponse{
		StatusCode: 0,
		StatusMsg:  nil,
		UserId:     0,
		Token:      "",
	}, nil
}

func (s UserAccountServer) Login(ctx context.Context, request *service.UserLoginRequest) (*service.UserLoginResponse, error) {
	return &service.UserLoginResponse{
		StatusCode: 0,
		StatusMsg:  nil,
		UserId:     0,
		Token:      "",
	}, nil
}

func (s UserAccountServer) UserInfo(ctx context.Context, request *service.UserInfoRequest) (*service.UserInfoResponse, error) {
	return &service.UserInfoResponse{
		StatusCode: 0,
		StatusMsg:  nil,
		User:       nil,
	}, nil
}

func main() {
	listen, err := net.Listen("tcp", ":8090")
	if err != nil {
		log.Fatalf("cannot create listener: %s", err)
	}
	serverRegister := grpc.NewServer()
	s := &UserAccountServer{}

	service.RegisterServiceServer(serverRegister, s)
	err = serverRegister.Serve(listen)
	if err != nil {
		log.Fatalf("failed: %s", err)
	}
}
