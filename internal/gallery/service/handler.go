package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "image-gallery/docs"
	"image-gallery/internal/gallery/entity"
	"image-gallery/internal/gallery/worker"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	Service
	Worker *worker.Worker
}

func NewHandler(s Service, Worker *worker.Worker) *Handler {
	return &Handler{
		Service: s,
		Worker:  Worker,
	}
}

// CreatePhoto godoc
//
//	@Summary		CreatPhoto
//	@Tags			gallery
//	@Description	CreatPhoto
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			input	body		entity.ImageRequest	true	"ImageRequest parameters"
//	@Success		200		{object}	entity.Response
//	@Failure		400		{object}	entity.Response
//	@Router			/api/v1/gallery/image [post]
func (h *Handler) CreatePhoto(c *gin.Context) {
	var photoRequest entity.ImageRequest

	if err := c.ShouldBindJSON(&photoRequest); err != nil {
		c.JSON(http.StatusBadRequest, entity.Response{Message: fmt.Sprintf("error %s", err.Error())})
		return
	}

	err := h.Service.CreatePhoto(&photoRequest, c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Response{Message: fmt.Sprintf("error %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image created successfully"})
}

// GetAllPhotos godoc
//
//	@Summary		GetAllPhotos
//	@Tags			admin
//	@Description	GetAllPhotos
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//
//	@Success		200	{object}	entity.PhotoResponse
//	@Failure		400	{object}	entity.PhotoResponse
//	@Router			/api/v1/admin/gallery/ [get]
func (h *Handler) GetAllPhotos(c *gin.Context) {
	photos, err := h.Service.GetAllPhotos(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Response{Message: fmt.Sprintf("error %s", err.Error())})
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

// GetGalleryById godoc
//
//	@Summary		GetGalleryById
//	@Tags			admin
//	@Description	GetGalleryById
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"user id"
//	@Success		200	{object}	entity.PhotoResponse
//	@Failure		400	{object}	entity.PhotoResponse
//	@Router			/api/v1/admin/gallery/{id} [get]
func (h *Handler) GetGalleryById(c *gin.Context) {
	userId := c.Param("id")
	id, err := strconv.Atoi(userId)
	if err != nil {
		log.Fatalf("Error in parsing string user id: %s", err)
	}

	photos, err := h.Service.GetGalleryById(id, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Response{Message: fmt.Sprintf("error %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, photos)
}

// AddTagName godoc
//
//	@Summary		AddTagName
//	@Tags			gallery
//	@Description	AddTagName
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			input	body		entity.TageRequest	true	"TageRequest parameters"
//	@Success		200		{object}	entity.Response
//	@Failure		400		{object}	entity.Response
//	@Router			/api/v1/gallery/addTag [post]
func (h *Handler) AddTagName(c *gin.Context) {
	var addingTag entity.TageRequest

	if err := c.ShouldBindJSON(&addingTag); err != nil {
		c.JSON(http.StatusBadRequest, entity.Response{Message: fmt.Sprintf("error %s", err.Error())})
		return
	}

	err := h.Service.AddTag(addingTag.TagName, addingTag.ImageId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Response{Message: fmt.Sprintf("error in adding tag handler %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, entity.Response{Message: "Tag added Successfully"})
}

// DeleteImage godoc
//
//	@Summary		Delete Image
//	@Tags			admin
//	@Description	Delete Image
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"image id"
//	@Success		200	{object}	entity.Response
//	@Failure		400	{object}	entity.Response
//	@Router			/api/v1/admin/gallery/{id} [delete]
func (h *Handler) DeleteImage(c *gin.Context) {
	imageId := c.Param("id")
	id, err := strconv.Atoi(imageId)
	if err != nil {
		log.Fatalf("Error in parsing string user id: %s", err)
	}
	if err := h.Service.DeleteImage(id, c); err != nil {
		c.JSON(http.StatusInternalServerError, entity.Response{Message: fmt.Sprintf("error %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, entity.Response{Message: "image delete successfully"})
}

// Follow godoc
//
//	@Summary		Follow
//	@Tags			gallery
//	@Description	Follow
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			input	body		entity.Username	true	"Username parameters"
//	@Success		200		{object}	entity.Response
//	@Failure		400		{object}	entity.Response
//	@Router			/api/v1/gallery/follow [post]
func (h *Handler) Follow(c *gin.Context) {
	var req entity.Username

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.Response{Message: fmt.Sprintf("error %s", err.Error())})
		return
	}

	err := h.Service.Follows(req.Username, c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Response{Message: fmt.Sprintf("error %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, entity.Response{Message: "You successfully follow"})
}

// Like godoc
//
//	@Summary		Like
//	@Tags			gallery
//	@Description	Like
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			input	body		entity.LikesRequest	true	"LikeRequest parameters"
//	@Success		200		{object}	entity.Response
//	@Failure		400		{object}	entity.Response
//	@Router			/api/v1/gallery/like [post]
func (h *Handler) Like(c *gin.Context) {
	var req entity.LikesRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.Response{Message: fmt.Sprintf("error %s", err.Error())})
		return
	}

	err := h.Service.Like(c, req.ImageId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Response{Message: fmt.Sprintf("error %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, entity.Response{Message: "You successfully like"})
}

// SearchPhotosByTag godoc
//
//	@Summary		search
//	@Tags			gallery
//	@Description	Search photos based on the specified tag
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			tag	query		string	true	"Tag to filter by"	default(popular)
//	@Success		200	{array}		entity.Image
//	@Failure		400	{object}	entity.Response
//	@Router			/api/v1/gallery/images [get]
func (h *Handler) SearchPhotosByTag(c *gin.Context) {
	tag := c.Query("tag")
	if tag == "" {
		c.JSON(http.StatusBadRequest, entity.Response{Message: "error: Tag parameter is required"})
		return
	}

	photos, err := h.Service.SearchPhotosByTag(tag, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, photos)
}

// GetImages @Summary Get images
//
//	@Description	Get a list of images based on sorting criteria
//	@Tags			images
//
//	@Security		BearerAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			sortKey	query		string	true	"Key to use for sorting (e.g., name, date)"
//	@Param			sortBy	query		string	true	"Sort order (e.g., asc, desc)"
//	@Success		200		{array}	entity.Image
//
//	@Failure		400		{object}	entity.Response
//	@Router			/api/v1/gallery/images/sort [get]
func (h *Handler) GetImages(c *gin.Context) {
	sortKey := c.Query("sortKey")
	sortBy := c.Query("sortBy")

	images, err := h.Service.GetImages(sortKey, sortBy, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Response{Message: fmt.Sprintf("error %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"images": images})
}

// GetImagesByFollowee @Summary Get images by followee
//
//	@Description	Get a list of images based on the followee's user ID
//	@Tags			images
//
//	@Security		BearerAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			id	path		integer	true	"User ID of the followee"
//	@Success		200	{array}	entity.Image
//	@Failure		400	{object}	entity.Response
//	@Failure		500	{object}	entity.Response
//	@Router			/api/v1/gallery/images/{id} [get]
func (h *Handler) GetImagesByFollowee(c *gin.Context) {
	userId := c.Param("id")
	id, err := strconv.Atoi(userId)
	if err != nil {
		log.Fatalf("Error in parsing string user id: %s", err)
	}

	res, err := h.Service.GetImagesByFollowing(id, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Response{Message: fmt.Sprintf("error %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"images": res})
}

// GetLikedImages @Summary Get liked images
//
//	@Description	Get a list of images that the user has liked
//	@Tags			gallery
//
//	@Security		BearerAuth
//
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	entity.Image
//	@Failure		500	{object}	entity.Response
//	@Router			/api/v1/gallery/images/like [get]
func (h *Handler) GetLikedImages(c *gin.Context) {
	images, err := h.Service.GetLikedImages(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Response{Message: fmt.Sprintf("error %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, gin.H{"images": images})
}

// UpdateImage @Summary Update an image
//
//	@Description	Update the details of an image
//	@Tags			images
//
//	@Security		BearerAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			input	body		entity.UpdateImageRequest	true	"Request body containing updated image details"
//	@Success		200		{object}	entity.Image
//	@Failure		400		{object}	entity.Response
//	@Failure		500		{object}	entity.Response
//	@Router			/api/v1/gallery/image/update [put]
func (h *Handler) UpdateImage(c *gin.Context) {
	var req *entity.UpdateImageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.Response{Message: fmt.Sprintf("error %s", err.Error())})
		return
	}

	image, err := h.Service.UpdateImage(c, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Response{Message: fmt.Sprintf("error %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"image": image})
}
