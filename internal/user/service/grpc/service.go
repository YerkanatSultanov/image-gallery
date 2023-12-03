package grpc

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"image-gallery/internal/user/entity"
	"image-gallery/internal/user/repo"
	pb "image-gallery/pkg/protobuf/userservice/gw"
	"image-gallery/pkg/token"
	"image-gallery/pkg/util"
	"strconv"
	"time"
)

type Service struct {
	pb.UnimplementedUserServiceServer
	logger  *zap.SugaredLogger
	repo    repo.Repository
	timeout time.Duration
	//userVerificationProducer *kafka.Producer
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
		//userVerificationProducer: userVerificationProducer,
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

	//randNum1 := rand.Intn(10)
	//randNum2 := rand.Intn(10)
	//randNum3 := rand.Intn(10)
	//randNum4 := rand.Intn(10)

	//msg := dto.UserCode{Code: fmt.Sprintf("%d%d%d%d", randNum1, randNum2, randNum3, randNum4)}
	//if err != nil {
	//	return nil, fmt.Errorf("failed to convert UserCode to int: %w", err)
	//}

	//code := &entity.UserCode{
	//	UserCode: msg.Code,
	//}

	r, err := s.repo.CreateUser(u)
	if err != nil {
		s.logger.Errorf("Cannot create user: %s", err)
		return nil, err
	}

	//code.UserId = int(r.Id)
	//
	//err = s.repo.UserCodeInsert(code)
	//if err != nil {
	//	s.logger.Errorf("Cannot insert user code: %s", err)
	//	return nil, err
	//}

	//b, err := json.Marshal(&msg)
	//if err != nil {
	//	s.logger.Errorf("Failed to marshal UserCode: %s", err)
	//	return nil, err
	//}
	//s.userVerificationProducer.ProduceMessage(b)

	res := &entity.CreateUserRes{
		Id:       strconv.FormatInt(r.Id, 10),
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
			Role:     u.Role,
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

func (s *Service) UpdateUser(targetUserID int, newUsername string, c *gin.Context) error {

	_, adminId, _, err := token.Claims(c)

	admin, err := s.repo.GetUserById(adminId)

	if admin.Role != "admin" {
		return fmt.Errorf("You don't have a permission for updating")
	}

	err = s.repo.UpdateUser(targetUserID, newUsername)
	if err != nil {
		return fmt.Errorf("error in updating user: %s", err)
	}

	return nil
}

func (s *Service) DeleteUser(id int, c *gin.Context) error {
	_, adminId, _, err := token.Claims(c)
	if err != nil {
		s.logger.Errorf("error extract an admin id: %s", err)
		return err
	}
	admin, err := s.repo.GetUserById(adminId)

	if admin.Role != "admin" {
		return fmt.Errorf("You don't have a permission for updating")
	}

	err = s.repo.DeleteUser(id)
	if err != nil {
		return fmt.Errorf("error in deleting user: %s", err)
	}

	return nil

}

func (s *Service) GetUserById(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.GetUserByIdResponse, error) {
	u, err := s.repo.GetUserById(int(req.Id))

	if err != nil {
		s.logger.Errorf("failed to GetUserByLogin err: %v", err)
		return nil, fmt.Errorf("GetUserByLogin err: %w", err)
	}

	return &pb.GetUserByIdResponse{
		Result: &pb.User{
			Id:       int32(u.Id),
			Email:    u.Email,
			Username: u.Username,
			Password: u.Password,
			Role:     u.Role,
		},
	}, nil
}

func (s *Service) GetUserByUsername(ctx context.Context, req *pb.GetUserByUsernameRequest) (*pb.GetUserByUsernameResponse, error) {
	u, err := s.repo.GetUserByUsername(req.Username)

	if err != nil {
		s.logger.Errorf("failed to GetUserByLogin err: %v", err)
		return nil, fmt.Errorf("GetUserByLogin err: %w", err)
	}

	return &pb.GetUserByUsernameResponse{
		Result: &pb.User{
			Id:       int32(u.Id),
			Email:    u.Email,
			Username: u.Username,
			Password: u.Password,
			Role:     u.Role,
		},
	}, nil
}

func (s *Service) CreateUserAdmin(req *entity.CreateUserReq, c *gin.Context) (*entity.CreateUserRes, error) {
	_, id, _, err := token.Claims(c)
	if err != nil {
		s.logger.Errorf("error in claims: %s", err)
		return nil, err
	}

	u, err := s.repo.GetUserById(id)
	if err != nil {
		s.logger.Errorf("error in getting admin id: %s", err)
		return nil, err
	}

	if u.Role != "admin" {
		s.logger.Info("You don't have a permissions")
		return nil, nil
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	r, err := s.repo.CreateUser(user)
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

func (s *Service) ConfirmUser(userCode string) error {
	return s.repo.ConfirmUser(userCode)
}
