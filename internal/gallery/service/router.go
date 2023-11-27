package service

import (
	"github.com/gin-gonic/gin"
)

func InitRouters(userHandler *Handler, r *gin.Engine) {

	group := r.Group("/api/v1/image")
	groupAdmin := r.Group("/api/v1/admin/image")

	group.POST("/create", userHandler.CreatePhoto)
	group.POST("/addTag", userHandler.AddTagName)

	groupAdmin.GET("/all", userHandler.GetAllPhotos)
	groupAdmin.GET("/:id", userHandler.GetGalleryById)
	groupAdmin.DELETE("/delete/:id", userHandler.DeleteImage)
}
