package service

import (
	"github.com/gin-gonic/gin"
	"image-gallery/internal/gallery/entity"
	"net/http"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) CreatePhoto(c *gin.Context) {
	var photoRequest entity.ImageRequest

	// Bind JSON request body to the ImageRequest struct
	if err := c.ShouldBindJSON(&photoRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the CreatePhoto function from the service
	err := h.Service.CreatePhoto(&photoRequest, c)

	// Check for errors and respond accordingly
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Photo created successfully"})
}

func (h *Handler) GetAllPhotos(c *gin.Context) {
	photos, err := h.Service.GetAllPhotos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	photoResponse := make([]*entity.PhotoResponse, len(photos))
	for i, photo := range photos {
		photoResponse[i] = &entity.PhotoResponse{
			Id:          photo.Id,
			UserId:      photo.UserId,
			Description: photo.Description,
			ImageLink:   photo.ImageLink,
			CreatedAt:   photo.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, photoResponse)
}
