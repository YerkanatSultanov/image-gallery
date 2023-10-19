package user

import "context"

type User struct {
	Id       int64  `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"email"`
}

type CreateUserReq struct {
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"email"`
}

type CreateUserRes struct {
	Id       string `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
}

type LogInReq struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"email"`
}

type LogInRes struct {
	accessToken string
	Id          string `json:"id" db:"id"`
	Username    string `json:"username" db:"username"`
}

type Repository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type Service interface {
	CreateUser(ctx context.Context, req *CreateUserReq) (*CreateUserRes, error)
	LogIn(c context.Context, req *LogInReq) (*LogInRes, error)
}
