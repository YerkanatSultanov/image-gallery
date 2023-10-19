package user

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"image-gallery/util"
	"strconv"
	"time"
)

const (
	secretKey = "secret"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(repository Repository) Service {
	return &service{
		repository,
		time.Duration(2) * time.Second,
	}
}

type MyJWTClaims struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s *service) CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	r, err := s.Repository.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}
	res := &CreateUserRes{
		Id:       strconv.Itoa(int(r.Id)),
		Username: r.Username,
		Email:    r.Email,
	}

	return res, nil
}

func (s *service) LogIn(c context.Context, req *LogInReq) (*LogInRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	u, err := s.Repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return &LogInRes{}, err
	}

	err = util.CheckPassword(req.Password, u.Password)
	if err != nil {
		return &LogInRes{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJWTClaims{
		Id:       strconv.Itoa(int(u.Id)),
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(u.Id)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	ss, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return &LogInRes{}, err
	}

	return &LogInRes{accessToken: ss, Id: strconv.Itoa(int(u.Id)), Username: u.Username}, nil

}
