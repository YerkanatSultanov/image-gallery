package service

import (
	"github.com/gin-gonic/gin"
	"image-gallery/internal/user/service/middleware"
)

var r = gin.Default()

func InitRouters(userHandler *Handler) {

	group := r.Group("/api/v1")

	group.POST("/signup", userHandler.CreateUser)
	group.GET("/user/:email", userHandler.GetUser)
	group.GET("user/all", middleware.JWTVerify(), userHandler.GetAllUsers)
}

func Start(addr string) error {
	return r.Run(addr)
}
