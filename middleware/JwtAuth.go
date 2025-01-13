package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func CheckJwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check Header Authorization
		header := c.Request.Header.Get("Authorization")
		hmacSampleSecret := []byte(os.Getenv("MY_SECRET"))
		tokenString := strings.Replace(header, "Bearer ", "", 1)

		// JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return hmacSampleSecret, nil
		})

		// Check Token
		if err != nil || !token.Valid {
			var message string
			if err != nil {
				if ve, ok := err.(*jwt.ValidationError); ok {
					if ve.Errors&jwt.ValidationErrorExpired != 0 {
						message = "Token is expired"
					} else {
						message = "Token is invalid"
					}
				} else {
					message = err.Error()
				}
			} else {
				message = "Invalid token"
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "error", "message": message})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := claims["userId"]
			userType := claims["userType"]

			if userID == nil || userType == nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid token"})
				return
			}

			userTypeStr, ok := userType.(string)
			if !ok || (userTypeStr != "user" && userTypeStr != "manager") {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid token"})
				return
			}

			c.Set("userId", claims["userId"])
			c.Set("userType", claims["userType"])
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid token"})
			return
		}

		c.Next()
	}
}
