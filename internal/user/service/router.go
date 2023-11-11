package service

import (
	"github.com/gin-gonic/gin"
)

func InitRouters(userHandler *Handler, r *gin.Engine) {

	group := r.Group("/api/v1")

	group.POST("/signup", userHandler.CreateUser)
	group.GET("/user/:email", userHandler.GetUser)
	group.GET("user/all", userHandler.GetAllUsers)
}
