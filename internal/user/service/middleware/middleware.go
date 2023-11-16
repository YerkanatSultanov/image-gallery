package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
)

func JWTVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		var accessToken string
		tokenHeader := c.Request.Header.Get("Authorization")
		tokenFields := strings.Fields(tokenHeader)

		if len(tokenFields) == 2 && tokenFields[0] == "Bearer" {
			accessToken = tokenFields[1]
		} else {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("secretKey"), nil
		})

		if err != nil {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		if !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userID, ok := claims["id"]

		if !ok {
			log.Printf("user id could not be parsed from JWT")
		}

		c.Set("id", userID)
		c.Next()

	}
}
