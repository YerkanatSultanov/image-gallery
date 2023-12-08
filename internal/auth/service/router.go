package service

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "image-gallery/docs"
)

func InitRouters(userHandler *Handler, r *gin.Engine) {
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	group := r.Group("/api/v1/auth")
	{
		group.POST("/login", userHandler.LogIn)
		group.PUT("/renew-token", userHandler.RenewToken)
	}

}
