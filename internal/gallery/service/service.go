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
	Follows(followeeUsername string, c *gin.Context) error
	Like(c *gin.Context, imageId int) error
	SearchPhotosByTag(tagString string, c *gin.Context) ([]*entity.Image, error)
	GetImages(sortKey, sortBy string, c *gin.Context) ([]*entity.Image, error)
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
	tokenString, id, _, err := token.Claims(c)

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

	_, id, _, err := token.Claims(c)
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
	_, id, _, err := token.Claims(c)
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

	_, id, _, err := token.Claims(c)

	if err != nil {
		s.logger.Errorf("error in claims of token: %s", err)
	}

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

func (s *service) Follows(followeeUsername string, c *gin.Context) error {
	_, id, _, err := token.Claims(c)
	if err != nil {
		s.logger.Errorf("error in claims of token: %s", err)
	}

	u, err := s.userGrpc.GetUserByUsername(c, followeeUsername)
	//TODO: Check user for authorization

	if err != nil {
		s.logger.Errorf("U can not to follow this user: %s", err)
	}

	err = s.repository.Follow(id, int(u.Id))

	if err != nil {
		s.logger.Fatalf("error to follow: %s", err)
	}

	return nil
}

func (s *service) Like(c *gin.Context, imageId int) error {
	_, id, _, err := token.Claims(c)
	if err != nil {
		s.logger.Errorf("error in claims of token: %s", err)
	}

	has, err := s.repository.UserHasImage(imageId)
	if !has {
		s.logger.Errorf("User does not have images")
	}
	if err != nil {
		s.logger.Fatalf("error in checkin images: %s", err)
	}
	ok, err := s.repository.UserLikedPhoto(id, imageId)

	req := &entity.Likes{
		UserId:  id,
		ImageId: imageId,
	}

	if !ok {
		err := s.repository.LikePhoto(req)
		if err != nil {
			s.logger.Fatalf("Can not like the image: %s", err)
		}
	}

	return nil
}

func (s *service) SearchPhotosByTag(tagString string, c *gin.Context) ([]*entity.Image, error) {
	tokenString, _, _, err := token.Claims(c)
	if err != nil {
		s.logger.Fatalf("Eroor in token: %s", err)
		return nil, err
	}
	ok, err := s.authGrpc.IsUserAuthorized(c, tokenString)
	if !ok.Authorized {
		s.logger.Fatalf("You are not authorized")
		return nil, err
	}

	photos, err := s.repository.FindImagesByTag(tagString)
	if err != nil {
		return nil, fmt.Errorf("error in service Search By Tag method %s", err)
	}

	photoResponse := make([]*entity.Image, len(photos))
	for i, photo := range photos {
		photoResponse[i] = &entity.Image{
			Id:          photo.Id,
			UserId:      photo.UserId,
			Description: photo.Description,
			ImageLink:   photo.ImageLink,
			CreatedAt:   photo.CreatedAt,
			UpdatedAt:   photo.UpdatedAt,
		}
	}

	return photoResponse, nil
}
func (s *service) GetImages(sortKey, sortBy string, c *gin.Context) ([]*entity.Image, error) {
	//tokenString, _, _, err := token.Claims(c)
	//if err != nil {
	//	s.logger.Fatalf("Eroor in token: %s", err)
	//	return nil, err
	//}
	//
	//ok, err := s.authGrpc.IsUserAuthorized(c, tokenString)
	//if !ok.Authorized {
	//	s.logger.Fatalf("You are not authorized")
	//	return nil, err
	//}

	images, err := s.repository.GetImages(sortKey, sortBy)
	if err != nil {
		return nil, fmt.Errorf("error in service GetImages: %w", err)
	}
	return images, nil
}
