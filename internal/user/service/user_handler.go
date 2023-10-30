package service

import (
	"github.com/gin-gonic/gin"
	"image-gallery/internal/user/entity"
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

func (h *Handler) CreateUser(c *gin.Context) {
	var u entity.CreateUserReq
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Service.CreateUser(c.Request.Context(), &u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
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

	c.SetCookie("access_token", u.AccessToken, 3600, "/", "localhost", false, true)
	c.SetCookie("refresh_token", u.RefreshToken, 3600, "/", "localhost", false, true)

	res := &entity.LogInRes{
		AccessToken:  u.AccessToken,
		RefreshToken: u.RefreshToken,
		Id:           u.Id,
		Username:     u.Username,
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) LogOut(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
