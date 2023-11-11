package service

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"image-gallery/internal/auth/config"
	"image-gallery/internal/auth/entity"
	"image-gallery/internal/auth/repo"
	"image-gallery/internal/auth/transport"
	"image-gallery/pkg/util"
	"log"
	"strconv"

	"time"
)

type service struct {
	repository repo.Repository
	timeout    time.Duration
	transport  *transport.UserTransport
	logger     *zap.SugaredLogger
}

type Config config.Auth

type Service interface {
	LogIn(ctx context.Context, req *entity.LogInReq) (*entity.UserTokenResponse, error)
}

func NewService(repository repo.Repository, userTransport *transport.UserTransport, logger *zap.SugaredLogger) Service {
	return &service{
		repository,
		time.Duration(2) * time.Second,
		userTransport,
		logger,
	}
}

func (s *service) LogIn(c context.Context, req *entity.LogInReq) (*entity.UserTokenResponse, error) {

	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	email := req.Email

	u, err := s.transport.GetUser(ctx, email)

	if err != nil {
		s.logger.Info("empty user")
		log.Fatalf("fail in GetUser %s", err)
		return &entity.UserTokenResponse{}, nil
	}

	err = util.CheckPassword(req.Password, u.Password)
	if err != nil {
		fmt.Printf("req password %s, u password %s", req.Password, u.Password)
		fmt.Println()
		s.logger.Info("Incorrect password")
		return &entity.UserTokenResponse{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, entity.MyJWTClaims{
		Id:       strconv.Itoa((u.Id)),
		Username: u.Username,
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:    strconv.Itoa((u.Id)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	tokenString, err := token.SignedString([]byte("secretKey"))
	if err != nil {
		s.logger.Info("incorrect secret key token")
		return &entity.UserTokenResponse{}, err
	}

	refreshTokenClaims := jwt.MapClaims{
		"user_id": u.Id,
		"exp":     time.Now().Add(time.Second * 1800),
	}

	refreshToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), refreshTokenClaims)

	refreshTokenString, err := refreshToken.SignedString([]byte("secretKey"))
	if err != nil {
		s.logger.Info("incorrect secret key refresh token")
		return &entity.UserTokenResponse{}, err
	}

	userToken := entity.UserToken{
		UserId:       u.Id,
		Token:        tokenString,
		RefreshToken: refreshTokenString,
	}

	err = s.repository.CreateUserToken(ctx, userToken)

	if err != nil {
		s.logger.Info("Can not save in database....")
		return nil, nil
	}

	return &entity.UserTokenResponse{UserId: u.Id, Token: tokenString, RefreshToken: refreshTokenString}, nil
}
