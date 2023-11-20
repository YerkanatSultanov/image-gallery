package service

import (
	"github.com/gin-gonic/gin"
	"image-gallery/internal/auth/entity"
	"image-gallery/internal/auth/service/grpc"
	"net/http"
)

type Handler struct {
	*grpc.Service
}

func NewHandler(s grpc.Service) *Handler {
	return &Handler{
		Service: &s,
	}
}

func (h *Handler) LogIn(c *gin.Context) {
	var user entity.LogInReq

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.Service.LogIn(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userToken := entity.UserToken{
		Id:           u.Id,
		UserId:       u.UserId,
		Username:     u.Username,
		Token:        u.Token,
		RefreshToken: u.RefreshToken,
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := &entity.UserTokenResponse{
		Id:           userToken.Id,
		UserId:       userToken.UserId,
		Username:     userToken.Username,
		Token:        userToken.Token,
		RefreshToken: userToken.RefreshToken,
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) RenewToken(c *gin.Context) {
	//userId := c.Param("id")

	userTokenResponse, err := h.Service.RenewToken(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userId":       userTokenResponse.UserId,
		"username":     userTokenResponse.Username,
		"token":        userTokenResponse.Token,
		"refreshToken": userTokenResponse.RefreshToken,
	})

}
