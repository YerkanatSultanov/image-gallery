package entity

import "github.com/golang-jwt/jwt/v5"

type User struct {
	Id       int64  `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"email"`
	Role     string `json:"role"`
}

type Claims struct {
	ID   int    `json:"id"`
	Role string `json:"role"`
	*jwt.RegisteredClaims
}

type UserResponse struct {
	Id       int64  `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
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
	Id           string `json:"id" db:"id"`
	Username     string `json:"username" db:"username"`
	AccessToken  string
	RefreshToken string
}
