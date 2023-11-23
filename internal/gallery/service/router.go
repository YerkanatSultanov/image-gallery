package service

import (
	"github.com/gin-gonic/gin"
)

func InitRouters(userHandler *Handler, r *gin.Engine) {

	group := r.Group("/api/v1/gallery")

	group.POST("/create", userHandler.CreatePhoto)
	group.GET("/getAllPhotos", userHandler.GetAllPhotos)
	group.GET("/getById/:id", userHandler.GetGalleryById)
	group.POST("/addTag", userHandler.AddTagName)
}
