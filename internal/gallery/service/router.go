package service

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "image-gallery/docs"
	"image-gallery/internal/gallery/service/middleware"
)

func InitRouters(userHandler *Handler, r *gin.Engine) {

	group := r.Group("/api/v1/gallery")
	groupAdmin := r.Group("/api/v1/admin/gallery")

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	group.POST("/image", middleware.JWTMiddleware(), userHandler.CreatePhoto)
	group.POST("/addTag", middleware.JWTMiddleware(), userHandler.AddTagName)
	group.POST("/follow", middleware.JWTMiddleware(), userHandler.Follow)
	group.GET("/images", middleware.JWTMiddleware(), userHandler.SearchPhotosByTag)
	group.GET("/images/sort", middleware.JWTMiddleware(), userHandler.GetImages)
	group.POST("/like", middleware.JWTMiddleware(), userHandler.Like)
	group.GET("/images/:id", middleware.JWTMiddleware(), userHandler.GetImagesByFollowee)
	group.GET("/images/like", middleware.JWTMiddleware(), userHandler.GetLikedImages)
	group.PUT("/image/update", middleware.JWTMiddleware(), userHandler.UpdateImage)

	groupAdmin.GET("/", userHandler.GetAllPhotos)
	groupAdmin.GET("/:id", userHandler.GetGalleryById)
	groupAdmin.DELETE("/:id", userHandler.DeleteImage)
}
