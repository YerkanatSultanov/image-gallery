package entity

import "github.com/golang-jwt/jwt/v5"

type MyJWTClaims struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	*jwt.RegisteredClaims
}
