package client

import (
	"github.com/MantoCoding/grpcDouyinDemo/user_service/user_grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
)

// 实现了 client 的单例模式

var C Clients

type Clients struct {
	UserLoginClient user_grpc.ServiceClient

	lock sync.Mutex
}

func (c *Clients) GetUserLoginClient() user_grpc.ServiceClient {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.UserLoginClient == nil {
		addr := ":8083"
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic("创建连接失败")
		}
		client := user_grpc.NewServiceClient(conn)
		c.UserLoginClient = client
	}
	return c.UserLoginClient
}
