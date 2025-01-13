package helper

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func GetUserByToken(c context.Context, tokens string) (string, string, error) {
	hmacSampleSecret := []byte(os.Getenv("MY_SECRET"))
	tokenString := strings.Replace(tokens, "Bearer ", "", 1)
	// JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})
	if err != nil {
		return "", "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userId, ok := claims["userId"].(string)
		if !ok {
			return "", "", fmt.Errorf("id is not a string")
		}

		userType, ok := claims["userType"].(string)
		if !ok {
			return "", "", fmt.Errorf("userType is not a string")
		}
		return userId, userType, nil
	}
	return "", "", fmt.Errorf("invalid token")

}

func GetUserHeader(c *gin.Context) (*string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, errors.New("Authorization header is missing")
	}

	userID, _, err := GetUserByToken(c, authHeader)
	if err != nil {
		return nil, err
	}

	return &userID, nil
}
