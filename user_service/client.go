package user_service

import (
	"context"
	"fmt"
	api "github.com/bz-2021/mini_douyin/api/client"
	pb "github.com/bz-2021/mini_douyin/user_service/user_grpc"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 实现了 client 端调用 server 端的方法

func UserLoginAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 客户端
		client := api.C.GetUserLoginClient()

		// 解析请求参数
		var req pb.UserLoginRequest
		req.Username = c.Query("username")
		req.Password = c.Query("password")

		// 调用gRPC服务
		result, err := client.Login(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status_msg": "Failed to call gRPC user_service"})
			return
		}
		// 处理登录响应
		if result.StatusMsg != nil {
			c.JSON(http.StatusOK, gin.H{
				"status_code": 0,
				"status_msg":  result.StatusMsg,
				"user_id":     result.UserId,
				"token":       result.Token,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{"status_msg": result.StatusMsg})
		}
		fmt.Println(result, err)
	}
}
