package service

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "image-gallery/docs"
	"image-gallery/internal/gallery/service/middleware"
)

func InitRouters(userHandler *Handler, r *gin.Engine) {

	//r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	group := r.Group("/api/v1/gallery")
	{
		group.Use(middleware.JWTMiddleware())
		group.POST("/image", userHandler.CreatePhoto)
		group.POST("/addTag", userHandler.AddTagName)
		group.POST("/follow", userHandler.Follow)
		group.GET("/images", userHandler.SearchPhotosByTag)
		group.GET("/images/sort", userHandler.GetImages)
		group.POST("/like", userHandler.Like)
		group.GET("/images/:id", userHandler.GetImagesByFollowee)
		group.GET("/images/like", userHandler.GetLikedImages)
		group.PUT("/image/update", userHandler.UpdateImage)
	}

	groupAdmin := r.Group("/api/v1/admin/gallery")
	{
		groupAdmin.Use(middleware.JWTMiddleware())
		groupAdmin.GET("/", userHandler.GetAllPhotos)
		groupAdmin.GET("/:id", userHandler.GetGalleryById)
		groupAdmin.DELETE("/:id", userHandler.DeleteImage)
	}

}
