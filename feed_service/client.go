package feed_service

import (
	"context"
	"fmt"
	api "github.com/bz-2021/mini_douyin/api/client"
	pb "github.com/bz-2021/mini_douyin/feed_service/feed_grpc"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func FeedAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 客户端
		client := api.C.GetFeedClient()

		// 解析请求参数
		var req pb.FeedRequest
		var err error
		req.LatestTime, err = strconv.ParseInt(c.Query("latest_time"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status_msg": "fail to convert latest_time to int64",
			})
		}
		req.Token = c.Query("token")

		// 调用gRPC服务
		result, err := client.FeedAction(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status_msg": "Failed to call gRPC user_service"})
			fmt.Println(err)
			return
		}
		// 处理登录响应
		if result.StatusMsg != nil {
			c.JSON(http.StatusOK, gin.H{
				"status_code": result.StatusCode,
				"status_msg":  result.StatusMsg,
				"video_list":  result.VideoList,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{"status_msg": result.StatusMsg})
		}
		fmt.Println(result, err)
	}

}
