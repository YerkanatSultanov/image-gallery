package service

import (
	"github.com/gin-gonic/gin"
	"image-gallery/internal/auth/entity"
	"image-gallery/internal/auth/repo"
	"net/http"
)

type Handler struct {
	Service
	repo.Repository
}

func NewHandler(s Service, repository repo.Repository) *Handler {
	return &Handler{
		Service:    s,
		Repository: repository,
	}
}

func (h *Handler) LogIn(c *gin.Context) {
	var user entity.LogInReq

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.Service.LogIn(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userToken := entity.UserToken{
		Id:           u.Id,
		UserId:       u.UserId,
		Token:        u.Token,
		RefreshToken: u.RefreshToken,
	}

	err = h.Repository.CreateUserToken(c, userToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := &entity.UserTokenResponse{
		Id:           userToken.Id,
		UserId:       userToken.UserId,
		Token:        userToken.Token,
		RefreshToken: userToken.RefreshToken,
	}

	c.JSON(http.StatusOK, res)
}
