package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/tushar27x/music-lib-api/config"
	"github.com/tushar27x/music-lib-api/models"
)

var jwtSecret = []byte(config.GetEnv("JWT_SECRET"))

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Header missing"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invaid token"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userId := uint(claims["user_id"].(float64))
		role := string(claims["role"].(string))

		var user models.User

		if err := config.DB.First(&user, userId).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.Set("userId", userId)
		c.Set("role", role)
		c.Next()
	}

}
