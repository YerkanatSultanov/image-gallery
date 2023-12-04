package service

import (
	"github.com/gin-gonic/gin"
	"image-gallery/internal/user/service/middleware"
)

func InitRouters(userHandler *Handler, r *gin.Engine) {

	groupClient := r.Group("/api/v1/user")
	groupClient.POST("/", userHandler.CreateUser)
	groupClient.POST("/confirm", userHandler.ConfirmUser)

	groupAdmin := r.Group("api/v1/admin/user")
	groupAdmin.GET("/:email", middleware.JWTVerify(), userHandler.GetUser)
	groupAdmin.POST("/", middleware.JWTVerify(), userHandler.CreateUserAdmin)
	groupAdmin.GET("/", middleware.JWTVerify(), userHandler.GetAllUsers)
	groupAdmin.PUT("/:id", middleware.JWTVerify(), userHandler.UpdateUser)
	groupAdmin.DELETE("/:id", middleware.JWTVerify(), userHandler.DeleteUser)

}
