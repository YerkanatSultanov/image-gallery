package service

import (
	"github.com/gin-gonic/gin"
)

var r = gin.Default()

func InitRouters(userHandler *Handler) {

	group := r.Group("/api/v1")

	group.POST("/signup", userHandler.CreateUser)

	group.POST("/login", userHandler.LogIn)
	group.GET("/logout", userHandler.LogOut)

}

func Start(addr string) error {
	return r.Run(addr)
}
