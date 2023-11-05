package service

import (
	"context"
	"fmt"
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
	GetUser(ctx context.Context, email string) (*entity.User, error)
	GetAllUsers(c context.Context) ([]*entity.UserResponse, error)
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

func (s *service) GetUser(c context.Context, email string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	u, err := s.Repository.GetUserByEmail(ctx, email)

	if err != nil {
		return &entity.User{}, err
	}

	res := &entity.User{
		Id:       u.Id,
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
	}

	return res, nil
}

func (s *service) GetAllUsers(c context.Context) ([]*entity.UserResponse, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	users, err := s.Repository.GetAllUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("error in service GetAllUsers method %s", err)
	}
	userResponses := make([]*entity.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = &entity.UserResponse{
			Id:       user.Id,
			Username: user.Username,
			Email:    user.Email,
		}
	}

	return userResponses, nil
}
