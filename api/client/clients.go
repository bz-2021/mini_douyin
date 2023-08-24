package client

import (
	"github.com/bz-2021/mini_douyin/feed_service/feed_grpc/user"
	"github.com/bz-2021/mini_douyin/feed_service/feed_grpc/video"
	comment "github.com/bz-2021/mini_douyin/interaction_service/comment/comment_grpc"
	favorite "github.com/bz-2021/mini_douyin/interaction_service/favorite/favorite_grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
)

// 实现了 client 的单例模式

var C Clients

type Clients struct {
	ServiceClient  user.ServiceClient
	FeedClient     video.ServiceClient
	FavoriteClient favorite.ServiceClient
	CommentClient  comment.ServiceClient
	PublishClient  video.ServiceClient
	lock           sync.Mutex
}

func (c *Clients) GetServiceClient() user.ServiceClient {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.ServiceClient == nil {
		addr := ":8083"
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic("创建连接失败")
		}
		client := user.NewServiceClient(conn)
		c.ServiceClient = client
	}
	return c.ServiceClient
}

func (c *Clients) GetFeedClient() video.ServiceClient {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.FeedClient == nil {
		addr := ":8084"
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic("创建连接失败")
		}
		client := video.NewServiceClient(conn)
		c.FeedClient = client
	}
	return c.FeedClient
}

func (c *Clients) GetFavoriteClient() favorite.ServiceClient {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.FavoriteClient == nil {
		addr := ":8085"
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic("创建连接失败")
		}
		client := favorite.NewServiceClient(conn)
		c.FavoriteClient = client
	}
	return c.FavoriteClient
}

func (c *Clients) GetCommentClient() comment.ServiceClient {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.CommentClient == nil {
		addr := ":8086"
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic("创建连接失败")
		}
		client := comment.NewServiceClient(conn)
		c.CommentClient = client
	}
	return c.CommentClient
}

func (c *Clients) GetPublishClient() video.ServiceClient {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.PublishClient == nil {
		addr := ":8087"
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic("创建连接失败")
		}
		client := video.NewServiceClient(conn)
		c.PublishClient = client
	}
	return c.PublishClient
}
