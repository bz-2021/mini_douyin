package user_service

import (
	"context"
	"fmt"
	api "github.com/bz-2021/mini_douyin/api/client"
	pb "github.com/bz-2021/mini_douyin/feed_service/feed_grpc"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// 实现了 client 端调用 server 端的方法

func UserLoginAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 客户端
		client := api.C.GetServiceClient()

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
				"status_code": result.StatusCode,
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

func UserRegisterAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 客户端
		client := api.C.GetServiceClient()

		// 解析请求参数
		var req pb.UserRegisterRequest
		req.Username = c.Query("username")
		req.Password = c.Query("password")

		// 调用gRPC服务
		result, err := client.Register(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status_msg": "Failed to call gRPC user_service"})
			return
		}
		// 处理登录响应
		if result.StatusMsg != nil {
			c.JSON(http.StatusOK, gin.H{
				"status_code": result.StatusCode,
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

func UserInfoAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 客户端
		client := api.C.GetServiceClient()

		// 解析请求参数
		var req pb.UserInfoRequest
		var err error
		req.UserId, err = strconv.ParseInt(c.Query("user_id"), 10, 64)
		req.Token = c.Query("token")

		// 调用gRPC服务
		result, err := client.UserInfo(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status_msg": "Failed to call gRPC user_service"})
			return
		}
		// 处理登录响应
		if result.StatusMsg != nil {
			c.JSON(http.StatusOK, gin.H{
				"status_code": result.StatusCode,
				"status_msg":  result.StatusMsg,
				"user": map[string]any{
					"id":               result.User.Id,
					"name":             result.User.Name,
					"avatar":           result.User.Avatar,
					"signature":        result.User.Signature,
					"is_follow":        false,
					"follow_count":     result.User.FollowCount,
					"follower_count":   result.User.FollowerCount,
					"background_image": result.User.BackgroundImage,
					"total_favorited":  result.User.TotalFavorited,
					"work_count":       result.User.WorkCount,
					"favorite_count":   result.User.FavoriteCount,
				},
			})
		} else {
			c.JSON(http.StatusOK, gin.H{"status_msg": result.StatusMsg})
		}
		fmt.Println(result, err)
	}
}
