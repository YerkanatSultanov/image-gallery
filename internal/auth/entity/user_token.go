package entity

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type UserToken struct {
	Id           int       `db:"id"`
	Token        string    `db:"token"`
	RefreshToken string    `db:"refresh_token"`
	UserId       int       `db:"user_id"`
	Username     string    `db:"username"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type MyJWTClaims struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	*jwt.RegisteredClaims
}

type UserTokenResponse struct {
	Id           int    `db:"id"`
	UserId       int    `db:"user_id"`
	Username     string `db:"username"`
	Token        string `db:"token"`
	RefreshToken string `db:"refresh_token"`
}

type LogInReq struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"email"`
}
