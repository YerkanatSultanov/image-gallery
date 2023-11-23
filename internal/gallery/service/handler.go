package service

import (
	"github.com/gin-gonic/gin"
	"image-gallery/internal/gallery/entity"
	"log"
	"net/http"
	"strconv"
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

	if err := c.ShouldBindJSON(&photoRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Service.CreatePhoto(&photoRequest, c)

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
func (h *Handler) GetGalleryById(c *gin.Context) {
	userId := c.Param("id")
	id, err := strconv.Atoi(userId)
	if err != nil {
		log.Fatalf("Error in parsing string user id: %s", err)
	}

	photos, err := h.Service.GetGalleryById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, photos)
}

func (h *Handler) AddTagName(c *gin.Context) {
	var addingTag entity.TageRequest

	if err := c.ShouldBindJSON(&addingTag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Service.AddTag(addingTag.TagName, addingTag.ImageId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error in adding tag handler": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tag added successfully"})
}
