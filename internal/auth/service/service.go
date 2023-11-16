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
	userGrpc   *transport.UserGrpcTransport
	logger     *zap.SugaredLogger
}

type Config config.Auth

type Service interface {
	LogIn(req *entity.LogInReq) (*entity.UserTokenResponse, error)
	RenewToken(userID string) (*entity.UserTokenResponse, error)
}

func NewService(repository repo.Repository, userTransport *transport.UserTransport, userGrpc *transport.UserGrpcTransport, logger *zap.SugaredLogger) Service {
	return &service{
		repository,
		time.Duration(2) * time.Second,
		userTransport,
		userGrpc,
		logger,
	}
}

func (s *service) LogIn(req *entity.LogInReq) (*entity.UserTokenResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	email := req.Email

	u, err := s.userGrpc.GetUserByEmail(ctx, email)

	if err != nil {
		log.Fatalf("fail in GetUserByEmail %s", err)
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
		Id: strconv.Itoa(int((u.Id))),
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int((u.Id))),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
		},
	})

	tokenString, err := token.SignedString([]byte("secretKey"))
	if err != nil {
		s.logger.Info("incorrect secret key token")
		return &entity.UserTokenResponse{}, err
	}

	refreshTokenClaims := jwt.MapClaims{
		"user_id": u.Id,
		"exp":     time.Now().Add(time.Second * 180),
	}

	refreshToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), refreshTokenClaims)

	refreshTokenString, err := refreshToken.SignedString([]byte("secretKey"))
	if err != nil {
		s.logger.Info("incorrect secret key refresh token")
		return &entity.UserTokenResponse{}, err
	}

	userToken := entity.UserToken{
		UserId:       int(u.Id),
		Username:     u.Username,
		Token:        tokenString,
		RefreshToken: refreshTokenString,
	}

	if err := s.repository.CreateUserToken(ctx, userToken); err != nil {
		s.logger.Errorf("failed to create user token: %s", err)
	}

	return &entity.UserTokenResponse{UserId: int(u.Id), Username: u.Username, Token: tokenString, RefreshToken: refreshTokenString}, nil
}

func (s *service) RenewToken(userID string) (*entity.UserTokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	id, err := strconv.Atoi(userID)
	if err != nil {
		s.logger.Fatalf("can not convert string to int: %s", err)
	}

	existingUserToken, err := s.repository.GetUserTokenByUserID(id)
	if err != nil {
		s.logger.Errorf("failed to get user token: %s", err)
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &entity.MyJWTClaims{
		Id:       strconv.Itoa(id),
		Username: existingUserToken.Username,
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(id),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
		},
	})

	tokenString, err := token.SignedString([]byte("secretKey"))
	if err != nil {
		s.logger.Error("failed to sign access token:", err)
		return nil, err
	}

	refreshTokenClaims := jwt.MapClaims{
		"user_id": id,
		"exp":     time.Now().Add(time.Second * 180).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	refreshTokenString, err := refreshToken.SignedString([]byte("secretKey"))
	if err != nil {
		s.logger.Error("failed to sign refresh token:", err)
		return nil, err
	}

	updatedUserToken := entity.UserToken{
		UserId:       id,
		Username:     existingUserToken.Username,
		Token:        tokenString,
		RefreshToken: refreshTokenString,
	}

	if err := s.repository.UpdateUserToken(ctx, updatedUserToken); err != nil {
		s.logger.Errorf("failed to update user token: %s", err)
		return nil, err
	}

	return &entity.UserTokenResponse{
		UserId:       id,
		Username:     updatedUserToken.Username,
		Token:        tokenString,
		RefreshToken: refreshTokenString,
	}, nil
}
