package main

import (
	"github.com/bz-2021/mini_douyin/api/pkg"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	//访问地址，处理我们的请求 Request Response

	pkg.InitRouter(r)

	err := r.Run(":8080")
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
