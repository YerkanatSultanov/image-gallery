package router

import (
	"github.com/gin-gonic/gin"
	"image-gallery/internal/user"
)

var r *gin.Engine

func InitRouters(userHandler *user.Handler) {
	r = gin.Default()

	r.POST("/signup", userHandler.CreateUser)
	r.POST("/login", userHandler.LogIn)
	r.GET("/logout", userHandler.LogOut)
}

func Start(addr string) error {
	return r.Run(addr)
}
