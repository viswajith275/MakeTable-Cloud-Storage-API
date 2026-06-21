package middleware

import (
	"fmt"
	"net/http"
	"timetable_api/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleWare(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		accessTokenString, err := ctx.Cookie("access_token")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		token, err := jwt.Parse(accessTokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(cfg.SecretKEY), nil
		})

		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
			})
			return
		}

		ctx.Set("user_id", claims["user_id"])
		ctx.Next()
	}
}
