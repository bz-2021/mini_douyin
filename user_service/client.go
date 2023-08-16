package user_service

import (
	"context"
	"encoding/json"
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call gRPC user_service 0814"})
			return
		}
		// 处理登录响应
		if result.StatusMsg != nil && *result.StatusMsg == "Succeed" {
			//c.JSON(http.StatusOK, gin.H{
			//	"message":    "登录成功",
			//	"StatusCode": result.StatusCode,
			//	"StatusMsg":  result.StatusMsg,
			//	"Token":      result.Token,
			//})
			jsonBytes, err := json.Marshal(result)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status_msg": "Failed to marshal response to JSON"})
				return
			}
			c.Header("Content-Type", "application/json")
			// 将 JSON 字节流作为响应返回给客户端
			c.Data(http.StatusOK, "application/json", jsonBytes)
		} else {
			c.JSON(http.StatusOK, gin.H{"status_msg": result.StatusMsg})
		}
		fmt.Println(result, err)
	}
}
