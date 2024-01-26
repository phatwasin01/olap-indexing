package api

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	router *gin.Engine
	redis  *redis.Client
}

func NewServer(redis *redis.Client) (*Server, error) {
	router := gin.Default()
	server := &Server{
		router: router,
		redis:  redis,
	}
	router.GET("/ping", server.Ping)
	router.GET("/check-redis", server.CheckRedis)
	router.POST("/add-cube-meta-data", server.AddCubeMetaData)
	router.POST("/get-cube-ids", server.RetreiveCubeByDimensions)
	return server, nil
}
func (server *Server) Start() error {
	return server.router.Run()
}
