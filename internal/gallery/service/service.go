package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"image-gallery/internal/gallery/entity"
	"image-gallery/internal/gallery/repo"
	"image-gallery/internal/gallery/transport"
	"image-gallery/pkg/token"
	"log"
	"time"
)

type service struct {
	repository repo.Repository
	timeout    time.Duration
	logger     *zap.SugaredLogger
	authGrpc   *transport.AuthGrpcTransport
	userGrpc   *transport.UserGrpc
}

type Service interface {
	CreatePhoto(ph *entity.ImageRequest, c *gin.Context) error
	GetAllPhotos(c *gin.Context) ([]*entity.PhotoResponse, error)
	GetGalleryById(id int, c *gin.Context) ([]*entity.PhotoResponse, error)
	AddTag(tagName string, imageId int) error
	DeleteImage(imageId int, c *gin.Context) error
}

func NewService(repository repo.Repository, logger *zap.SugaredLogger, authGrpc *transport.AuthGrpcTransport, userGrpc *transport.UserGrpc) Service {
	return &service{
		repository,
		time.Duration(2) * time.Second,
		logger,
		authGrpc,
		userGrpc,
	}
}

func (s *service) CreatePhoto(ph *entity.ImageRequest, c *gin.Context) error {
	tokenString, id, err := token.TokenStringClaims(c)

	photo := &entity.Image{
		UserId:      id,
		Description: ph.Description,
		ImageLink:   ph.ImageLink,
	}

	b, err := s.authGrpc.IsUserAuthorized(c, tokenString)

	if err != nil {
		return fmt.Errorf("ERROR IN SERVICE CHECK OF AUTHORIZATION: %S", err)
	}

	if !b.Authorized {
		return fmt.Errorf("USER NOT AUTHORIZED")
	}

	if err := s.repository.CreatePhoto(photo); err != nil {
		return fmt.Errorf("Error in creating photo: %s", err)
	}

	return nil
}

func (s *service) GetAllPhotos(c *gin.Context) ([]*entity.PhotoResponse, error) {

	_, id, err := token.TokenStringClaims(c)
	u, err := s.userGrpc.GetUserById(c, id)
	if u.Role != "admin" {
		s.logger.Fatalf("You don't have a permissions for getting gallery: %s", err)
		return nil, err
	}

	photos, err := s.repository.GetAllPhotos()

	if err != nil {
		return nil, fmt.Errorf("error in service GetAllPhotos method %s", err)
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

	return photoResponse, nil
}

func (s *service) GetGalleryById(targetId int, c *gin.Context) ([]*entity.PhotoResponse, error) {
	_, id, err := token.TokenStringClaims(c)
	u, err := s.userGrpc.GetUserById(c, id)
	if u.Role != "admin" {
		s.logger.Fatalf("You don't have a permissions for getting gallery: %s", err)
		return nil, err
	}

	photos, err := s.repository.GetGalleryById(targetId)
	if err != nil {
		return nil, fmt.Errorf("error in service GetAllPhotos method %s", err)
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

	return photoResponse, nil
}
func (s *service) AddTag(tagName string, imageId int) error {
	tx, err := s.repository.BeginTransaction()
	if err != nil {
		return fmt.Errorf("error starting transaction: %s", err)
	}
	defer func() {
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				log.Fatalf("Error in RollBack")
				return
			}
		}
	}()

	t := &entity.Tags{
		TagName: tagName,
	}
	tagID, err := s.repository.AddTagName(t)
	if err != nil {
		return fmt.Errorf("error adding tag name: %s", err)
	}

	err = s.repository.AddTagImage(tagID, imageId)
	if err != nil {
		return fmt.Errorf("error adding image with tag: %s", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %s", err)
	}

	return nil
}

func (s *service) DeleteImage(imageId int, c *gin.Context) error {

	_, id, err := token.TokenStringClaims(c)
	u, err := s.userGrpc.GetUserById(c, id)
	if u.Role != "admin" {
		s.logger.Fatalf("You don't have a permissions for getting gallery: %s", err)
		return err
	}

	err = s.repository.DeleteImage(imageId)

	if err != nil {
		s.logger.Fatalf("Can not delete user: %s", err)
	}

	return nil
}
