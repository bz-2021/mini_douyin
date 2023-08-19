package favorite_service

import (
	"context"
	"fmt"
	api "github.com/bz-2021/mini_douyin/api/client"
	pb "github.com/bz-2021/mini_douyin/interaction_service/favorite/favorite_grpc"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func PostFavoriteAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 客户端
		client := api.C.GetFavoriteClient()

		// 解析请求参数
		var req pb.FavoriteActionRequest
		var err error

		req.Token = c.Query("token")

		videoIDStr := c.Query("video_id")
		req.VideoId, _ = strconv.ParseInt(videoIDStr, 10, 64)

		actionTypeStr := c.Query("action_type")
		actionTypeInt64, _ := strconv.ParseInt(actionTypeStr, 10, 64)
		req.ActionType = int32(actionTypeInt64)

		fmt.Printf("\n req.Token: %v \n req.VideoId: %v \n req.ActionType: %v \n", req.Token, req.VideoId, req.ActionType)

		// 调用gRPC服务
		result, err := client.FavoriteAction(context.Background(), &req)
		if err != nil {
			fmt.Println("调用出错")
			c.JSON(http.StatusInternalServerError, gin.H{"status_msg": "Failed to call gRPC user_service"})
			fmt.Println(err)
			return
		}
		// 处理登录响应
		if &result.StatusMsg != nil {
			c.JSON(http.StatusOK, gin.H{
				"status_code": result.StatusCode,
				"status_msg":  result.StatusMsg,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{"status_msg": result.StatusMsg})
		}
		fmt.Println(result, err)
	}

}
