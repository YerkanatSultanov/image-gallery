package service

import (
	"github.com/gin-gonic/gin"
)

var r = gin.Default()

func InitRouters(userHandler *Handler) {

	group := r.Group("/api/v1/auth")

	group.POST("/login", userHandler.LogIn)
}

func Start(addr string) error {
	return r.Run(addr)
}
