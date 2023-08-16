package pkg

import (
	"fmt"
	"github.com/bz-2021/mini_douyin/user_service"
	pb "github.com/bz-2021/mini_douyin/user_service/user_grpc"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"net"
)

func InitRouter(r *gin.Engine) {

	go func() {
		// 创建gRPC服务
		grpcServer := grpc.NewServer()

		// 注册LoginService服务
		//loginSrv := &user_service.UserLoginService{db: db} // 传入GORM数据库连接
		pb.RegisterServiceServer(grpcServer, user_service.NewUserLoginService())
		fmt.Println("grpc server running : 8083 ")

		listen, err := net.Listen("tcp", ":8083")
		if err != nil {
			grpclog.Fatalf("Failed to listen: %v", err)
		}

		if err := grpcServer.Serve(listen); err != nil {

		}
	}()

	//获取请求参数，调用grpc客户端
	r.POST("/douyin/user/login/", user_service.UserLoginAction())
	r.GET("/douyin/user/", user_service.UserInfoAction())
	r.POST("/douyin/user/register/", user_service.UserRegisterAction())

	// basic apis
	//apiRouter.GET("/feed/", controller.Feed)
	//apiRouter.GET("/user/", controller.UserInfo)
	//apiRouter.POST("/user/register/", controller.Register)
	//apiRouter.POST("/user/login/", user_service.Login)
	//apiRouter.POST("/publish/action/", user_service.Publish)

}
