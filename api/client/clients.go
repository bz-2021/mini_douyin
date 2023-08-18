package client

import (
	"github.com/bz-2021/mini_douyin/feed_service/feed_grpc/user"
	"github.com/bz-2021/mini_douyin/feed_service/feed_grpc/video"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
)

// 实现了 client 的单例模式

var C Clients

type Clients struct {
	ServiceClient user.ServiceClient
	FeedClient    video.ServiceClient

	lock sync.Mutex
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
