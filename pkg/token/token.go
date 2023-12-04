package token

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"strings"
)

func ExtractTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header is missing")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("Invalid Authorization header format")
	}

	return parts[1], nil
}

func ParseJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secretKey"), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("Invalid JWT token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Failed to extract claims from JWT token")
	}

	return claims, nil
}

func Claims(c *gin.Context) (string, int, string, error) {
	tokenString, err := ExtractTokenFromHeader(c)
	if err != nil {
		return "", 0, "", fmt.Errorf("failed to extract token: %s", err)
	}

	claims, err := ParseJWT(tokenString)
	if err != nil {
		return "", 0, "", fmt.Errorf("failed to parse JWT: %s", err)
	}

	userID, ok := claims["id"].(string)
	if !ok {
		return "", 0, "", fmt.Errorf("failed to extract user ID from JWT claims")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", 0, "", fmt.Errorf("failed to extract user username from JWT claims")
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		return "", 0, "", fmt.Errorf("can not convert string to int: %s", err)
	}

	return tokenString, id, username, nil
}
