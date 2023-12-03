package service

import (
	"github.com/gin-gonic/gin"
	"image-gallery/internal/gallery/service/middleware"
)

func InitRouters(userHandler *Handler, r *gin.Engine) {

	group := r.Group("/api/v1/gallery")
	groupAdmin := r.Group("/api/v1/admin/gallery")

	group.POST("/image", userHandler.CreatePhoto)
	group.POST("/addTag", userHandler.AddTagName)
	group.POST("/follow", userHandler.Follow)
	group.POST("/images", userHandler.SearchPhotosByTag)
	group.POST("/like", middleware.JWTVerify(), userHandler.Like)

	groupAdmin.GET("/all", userHandler.GetAllPhotos)
	groupAdmin.GET("/:id", userHandler.GetGalleryById)
	groupAdmin.DELETE("/delete/:id", userHandler.DeleteImage)
}
