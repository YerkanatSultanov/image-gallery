package service

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"image-gallery/internal/user/entity"
	"image-gallery/internal/user/repo"
	"image-gallery/pkg/util"
	"strconv"
	"time"
)

const (
	secretKey = "secret"
)

type service struct {
	repo.Repository
	timeout time.Duration
}

type Service interface {
	CreateUser(ctx context.Context, req *entity.CreateUserReq) (*entity.CreateUserRes, error)
	LogIn(c context.Context, req *entity.LogInReq) (*entity.LogInRes, error)
}

func NewService(repository repo.Repository) Service {
	return &service{
		repository,
		time.Duration(2) * time.Second,
	}
}

func (s *service) CreateUser(c context.Context, req *entity.CreateUserReq) (*entity.CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	r, err := s.Repository.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}
	res := &entity.CreateUserRes{
		Id:       strconv.Itoa(int(r.Id)),
		Username: r.Username,
		Email:    r.Email,
	}

	return res, nil
}

func (s *service) LogIn(c context.Context, req *entity.LogInReq) (*entity.LogInRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	u, err := s.Repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return &entity.LogInRes{}, err
	}

	err = util.CheckPassword(req.Password, u.Password)
	if err != nil {
		return &entity.LogInRes{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, entity.MyJWTClaims{
		Id:       strconv.Itoa(int(u.Id)),
		Username: u.Username,
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(u.Id)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return &entity.LogInRes{}, err
	}

	refreshTokenClaims := jwt.MapClaims{
		"user_id": u.Id,
		"exp":     time.Now().Add(time.Second * 1800),
	}

	refreshToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), refreshTokenClaims)

	refreshTokenString, err := refreshToken.SignedString([]byte(secretKey))
	if err != nil {
		return &entity.LogInRes{}, err
	}

	return &entity.LogInRes{AccessToken: tokenString, RefreshToken: refreshTokenString, Id: strconv.Itoa(int(u.Id)), Username: u.Username}, nil

}
