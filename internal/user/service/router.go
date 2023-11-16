package service

import (
	"github.com/gin-gonic/gin"
	"image-gallery/internal/user/service/middleware"
)

func InitRouters(userHandler *Handler, r *gin.Engine) {

	group := r.Group("/api/v1")

	group.POST("/signup", userHandler.CreateUser)
	group.GET("/user/:email", middleware.JWTVerify(), userHandler.GetUser)
	group.GET("/user/all", userHandler.GetAllUsers)
}
