package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (server *Server) Ping(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "pong",
	})
}

func (server *Server) CheckRedis(ctx *gin.Context) {
	val, err := server.redis.Get(ctx, "foo").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("foo", val)
	ctx.JSON(200, gin.H{
		"message": val,
	})

}
