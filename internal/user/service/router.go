package service

import (
	"github.com/gin-gonic/gin"
	"image-gallery/internal/user/service/middleware"
)

func InitRouters(userHandler *Handler, r *gin.Engine) {

	groupClient := r.Group("/api/v1/user")
	groupAdmin := r.Group("api/v1/admin/user")

	groupClient.POST("/signup", userHandler.CreateUser)
	groupAdmin.GET("/:email", middleware.JWTVerify(), userHandler.GetUser)

	groupAdmin.POST("/", userHandler.CreateUserAdmin)
	groupAdmin.GET("/all", userHandler.GetAllUsers)
	groupAdmin.POST("/update/:id", userHandler.UpdateUser)
	groupAdmin.DELETE("/delete/:id", userHandler.DeleteUser)

}
