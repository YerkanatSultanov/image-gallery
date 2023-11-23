package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"image-gallery/internal/gallery/entity"
	"image-gallery/internal/gallery/repo"
	"image-gallery/internal/gallery/service/token"
	"image-gallery/internal/gallery/transport"
	"log"
	"strconv"
	"time"
)

type service struct {
	repository repo.Repository
	timeout    time.Duration
	logger     *zap.SugaredLogger
	authGrpc   *transport.AuthGrpcTransport
}

type Service interface {
	CreatePhoto(ph *entity.ImageRequest, c *gin.Context) error
	GetAllPhotos() ([]*entity.PhotoResponse, error)
	GetGalleryById(id int) ([]*entity.PhotoResponse, error)
	AddTag(tagName string, imageId int) error
}

func NewService(repository repo.Repository, logger *zap.SugaredLogger, authGrpc *transport.AuthGrpcTransport) Service {
	return &service{
		repository,
		time.Duration(2) * time.Second,
		logger,
		authGrpc,
	}
}

func (s *service) CreatePhoto(ph *entity.ImageRequest, c *gin.Context) error {
	tokenString, err := token.ExtractTokenFromHeader(c)
	if err != nil {
		s.logger.Error("failed to extract token:", err)
		return err
	}

	claims, err := token.ParseJWT(tokenString)
	if err != nil {
		s.logger.Error("failed to parse JWT:", err)
		return err
	}

	userID, ok := claims["id"].(string)
	if !ok {
		return errors.New("failed to extract user ID from JWT claims")
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		s.logger.Fatalf("can not convert string to int: %s", err)
	}

	photo := &entity.Photo{
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

func (s *service) GetAllPhotos() ([]*entity.PhotoResponse, error) {
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

func (s *service) GetGalleryById(id int) ([]*entity.PhotoResponse, error) {
	photos, err := s.repository.GetGalleryById(id)
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
