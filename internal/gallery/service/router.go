package service

import (
	"github.com/gin-gonic/gin"
	"image-gallery/internal/gallery/service/middleware"
)

func InitRouters(userHandler *Handler, r *gin.Engine) {

	group := r.Group("/api/v1/gallery")
	groupAdmin := r.Group("/api/v1/admin/gallery")

	group.POST("/image", userHandler.CreatePhoto)
	group.POST("/addTag", middleware.JWTVerify(), userHandler.AddTagName)
	group.POST("/follow", middleware.JWTVerify(), userHandler.Follow)
	group.POST("/images", middleware.JWTVerify(), userHandler.SearchPhotosByTag)
	group.GET("/images/sort", middleware.JWTVerify(), userHandler.GetImages)
	group.POST("/like", middleware.JWTVerify(), userHandler.Like)
	group.GET("/images/:id", middleware.JWTVerify(), userHandler.GetImagesByFollowee)

	groupAdmin.GET("/", userHandler.GetAllPhotos)
	groupAdmin.GET("/:id", userHandler.GetGalleryById)
	groupAdmin.DELETE("/:id", userHandler.DeleteImage)
}
