package middleware

import (
	"errors"
	"go-trades/config"
	"go-trades/entity"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(401, gin.H{"error": "authorization header required"})
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(401, gin.H{"error": "invalid authorization header format"})
			ctx.Abort()
			return
		}

		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return config.GetJWTSecret(), nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(401, gin.H{"error": "invalid or expired token"})
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.JSON(401, gin.H{"error": "invalid token claims"})
			ctx.Abort()
			return
		}

		userId, ok := claims["userId"].(float64)
		if !ok {
			ctx.JSON(401, gin.H{"error": "invalid user ID in token"})
			ctx.Abort()
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			ctx.JSON(401, gin.H{"error": "invalid role in token"})
			ctx.Abort()
			return
		}

		ctx.Set("userId", uint(userId))
		ctx.Set("role", entity.Role(role))
		ctx.Next()
	}
}

func IsAdmin(ctx *gin.Context) (bool, error) {
	role, exists := ctx.Get("role")
	if !exists {
		return false, errors.New("role not found in context")
	}

	if role.(entity.Role) != entity.Admin {
		return false, nil

	}
	return true, nil
}
