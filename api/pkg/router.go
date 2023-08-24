package pkg

import (
	"fmt"
	"github.com/bz-2021/mini_douyin/feed_service"
	pb "github.com/bz-2021/mini_douyin/feed_service/feed_grpc/user"
	"github.com/bz-2021/mini_douyin/feed_service/feed_grpc/video"
	comment_service "github.com/bz-2021/mini_douyin/interaction_service/comment"
	comment "github.com/bz-2021/mini_douyin/interaction_service/comment/comment_grpc"
	"github.com/bz-2021/mini_douyin/interaction_service/favorite"
	favorite "github.com/bz-2021/mini_douyin/interaction_service/favorite/favorite_grpc"
	"github.com/bz-2021/mini_douyin/user_service"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"net"
)

func InitRouter(r *gin.Engine) {

	//LoginService
	go func() {
		// 创建gRPC服务
		grpcServer := grpc.NewServer()

		// 注册LoginService服务
		pb.RegisterServiceServer(grpcServer, user_service.NewUserLoginService())
		fmt.Println("grpc server running : 8083 ")

		listen, err := net.Listen("tcp", ":8083")
		if err != nil {
			grpclog.Fatalf("Failed to listen: %v", err)
		}

		if err := grpcServer.Serve(listen); err != nil {

		}
	}()

	//FeedService
	go func() {
		// 创建gRPC服务
		grpcServer := grpc.NewServer()

		// 注册FeedService服务
		video.RegisterServiceServer(grpcServer, feed_service.NewFeedService())
		fmt.Println("grpc server running : 8084 ")

		listen, err := net.Listen("tcp", ":8084")
		if err != nil {
			grpclog.Fatalf("Failed to listen: %v", err)
		}

		if err := grpcServer.Serve(listen); err != nil {

		}
	}()

	//FavoriteService
	go func() {
		// 创建gRPC服务
		grpcServer := grpc.NewServer()

		// 注册FavoriteService服务
		favorite.RegisterServiceServer(grpcServer, favorite_service.NewFavoriteService())
		fmt.Println("grpc server running : 8085 ")

		listen, err := net.Listen("tcp", ":8085")
		if err != nil {
			grpclog.Fatalf("Failed to listen: %v", err)
		}

		if err := grpcServer.Serve(listen); err != nil {

		}
	}()

	//CommentService
	go func() {
		// 创建gRPC服务
		grpcServer := grpc.NewServer()

		// 注册FavoriteService服务
		comment.RegisterServiceServer(grpcServer, comment_service.NewCommentService())
		fmt.Println("grpc server running : 8086 ")

		listen, err := net.Listen("tcp", ":8086")
		if err != nil {
			grpclog.Fatalf("Failed to listen: %v", err)
		}

		if err := grpcServer.Serve(listen); err != nil {

		}
	}()

	////PublishService
	go func() {
		// 创建gRPC服务
		grpcServer := grpc.NewServer()

		// 注册PublishService服务
		video.RegisterServiceServer(grpcServer, feed_service.NewFeedService())
		fmt.Println("grpc server running : 8087 ")

		listen, err := net.Listen("tcp", ":8087")
		if err != nil {
			grpclog.Fatalf("Failed to listen: %v", err)
		}

		if err := grpcServer.Serve(listen); err != nil {

		}
	}()

	//获取请求参数，调用grpc客户端

	//视频feed流 api
	r.GET("/douyin/feed/", feed_service.FeedAction())
	r.POST("/douyin/publish/action/", feed_service.PublishAction())
	//user apis
	r.POST("/douyin/user/login/", user_service.UserLoginAction())
	r.GET("/douyin/user/", user_service.UserInfoAction())
	r.POST("/douyin/user/register/", user_service.UserRegisterAction())

	//favorite apis
	r.GET("/douyin/favorite/list/", favorite_service.GetFavoriteList())
	r.POST("/douyin/favorite/action/", favorite_service.PostFavoriteAction())

	//comment apis
	r.GET("/douyin/comment/list/", comment_service.GetCommentList())
	r.POST("/douyin/comment/action/", comment_service.PostCommentAction())

}
