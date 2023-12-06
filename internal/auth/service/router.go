package service

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "image-gallery/docs"
)

func InitRouters(userHandler *Handler, r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	group := r.Group("/api/v1/auth")
	group.POST("/login", userHandler.LogIn)
	group.PUT("/renew-token", userHandler.RenewToken)
}
