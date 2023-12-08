package service

import (
	"github.com/gin-gonic/gin"
	"image-gallery/internal/user/service/middleware"
)

func InitRouters(userHandler *Handler, r *gin.Engine) {

	//r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	groupClient := r.Group("/api/v1/user")
	{
		groupClient.POST("/", userHandler.CreateUser)
		groupClient.POST("/confirm", userHandler.ConfirmUser)
	}

	groupAdmin := r.Group("api/v1/admin/user")
	{
		groupAdmin.Use(middleware.JWTMiddleware())
		groupAdmin.GET("/:email", userHandler.GetUser)
		groupAdmin.POST("/", userHandler.CreateUserAdmin)
		groupAdmin.GET("/", userHandler.GetAllUsers)
		groupAdmin.PUT("/:id", userHandler.UpdateUser)
		groupAdmin.DELETE("/:id", userHandler.DeleteUser)
	}

}
