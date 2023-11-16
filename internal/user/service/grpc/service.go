package grpc

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"image-gallery/internal/user/entity"
	"image-gallery/internal/user/repo"
	pb "image-gallery/pkg/protobuf/userservice/gw"
	"image-gallery/pkg/util"
	"strconv"
	"time"
)

type Service struct {
	pb.UnimplementedUserServiceServer
	logger  *zap.SugaredLogger
	repo    repo.Repository
	timeout time.Duration
}

//type ServiceInt interface {
//	CreateUser(ctx context.Context, req *entity.CreateUserReq) (*entity.CreateUserRes, error)
//	GetUser(ctx context.Context, email string) (*entity.User, error)
//	GetAllUsers(ctx context.Context) ([]*entity.UserResponse, error)
//	GetUserByEmail(c context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error)
//}

func NewService(repository *repo.Repository, logger *zap.SugaredLogger) *Service {
	return &Service{
		repo:   *repository,
		logger: logger,
	}
}

func (s *Service) CreateUser(req *entity.CreateUserReq) (*entity.CreateUserRes, error) {

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	r, err := s.repo.CreateUser(u)
	if err != nil {
		s.logger.Errorf("Can not Create user: %s", err)
		return nil, err
	}
	res := &entity.CreateUserRes{
		Id:       strconv.Itoa(int(r.Id)),
		Username: r.Username,
		Email:    r.Email,
	}

	return res, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error) {
	u, err := s.repo.GetUserByEmail(req.Email)

	if err != nil {
		s.logger.Errorf("failed to GetUserByLogin err: %v", err)
		return nil, fmt.Errorf("GetUserByLogin err: %w", err)
	}

	return &pb.GetUserByEmailResponse{
		Result: &pb.User{
			Id:       int32(u.Id),
			Email:    u.Email,
			Username: u.Username,
			Password: u.Password,
		},
	}, nil
}

func (s *Service) GetUser(email string) (*entity.User, error) {
	u, err := s.repo.GetUserByEmail(email)

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

func (s *Service) GetAllUsers() ([]*entity.UserResponse, error) {
	users, err := s.repo.GetAllUsers()
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
