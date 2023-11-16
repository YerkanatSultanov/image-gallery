package service

import (
	"github.com/gin-gonic/gin"
	"image-gallery/internal/user/entity"
	"image-gallery/internal/user/service/grpc"
	"net/http"
)

type Handler struct {
	*grpc.Service
}

func NewHandler(s *grpc.Service) *Handler {
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

	res, err := h.Service.CreateUser(&u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetUserByEmail(c *gin.Context) {

}

func (h *Handler) GetUser(c *gin.Context) {
	email := c.Param("email")

	u, err := h.Service.GetUser(email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := &entity.User{
		Id:       u.Id,
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetAllUsers(c *gin.Context) {
	users, err := h.Service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userResponses := make([]entity.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = entity.UserResponse{
			Id:       user.Id,
			Username: user.Username,
			Email:    user.Email,
		}
	}

	c.JSON(http.StatusOK, userResponses)
}
