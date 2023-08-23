package comment_service

import (
	"context"
	"fmt"
	api "github.com/bz-2021/mini_douyin/api/client"
	pb "github.com/bz-2021/mini_douyin/interaction_service/comment/comment_grpc"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func PostCommentAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 客户端
		client := api.C.GetCommentClient()

		// 解析请求参数
		var req pb.CommentActionRequest
		var err error

		req.Token = c.Query("token")

		videoIDStr := c.Query("video_id")
		req.VideoId, _ = strconv.ParseInt(videoIDStr, 10, 64)

		actionTypeStr := c.Query("action_type")
		actionTypeInt64, _ := strconv.ParseInt(actionTypeStr, 10, 64)
		req.ActionType = int32(actionTypeInt64)

		commentIDStr := c.Query("comment_id")
		req.CommentId, _ = strconv.ParseInt(commentIDStr, 10, 64)

		req.CommentText = c.Query("comment_text")

		// 调用gRPC服务
		result, err := client.CommentAction(context.Background(), &req)
		if err != nil {
			fmt.Println("调用出错")
			c.JSON(http.StatusInternalServerError, gin.H{"status_msg": "Failed to call gRPC favorite_service"})
			fmt.Println(err)
			return
		}
		// 处理登录响应
		if &result.StatusMsg != nil {
			c.JSON(http.StatusOK, gin.H{
				"status_code": result.StatusCode,
				"status_msg":  result.StatusMsg,
				"comment":     result.Comment,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{"status_msg": result.StatusMsg})
		}
		fmt.Println(result, err)
	}

}

func GetCommentList() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 客户端
		client := api.C.GetCommentClient()

		// 解析请求参数
		var req pb.CommentListRequest
		var err error

		req.Token = c.Query("token")

		videoIDStr := c.Query("video_id")
		req.VideoId, _ = strconv.ParseInt(videoIDStr, 10, 64)

		// 调用gRPC服务
		result, err := client.CommentList(context.Background(), &req)
		if err != nil {
			fmt.Println("调用出错")
			c.JSON(http.StatusInternalServerError, gin.H{"status_msg": "Failed to call gRPC favorite_service"})
			fmt.Println(err)
			return
		}
		// 处理登录响应
		if &result.StatusMsg != nil {
			c.JSON(http.StatusOK, gin.H{
				"status_code":  result.StatusCode,
				"status_msg":   result.StatusMsg,
				"comment_list": result.CommentList,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{"status_msg": result.StatusMsg})
		}
		fmt.Println(result, err)
	}

}
