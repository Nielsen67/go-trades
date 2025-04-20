package middleware

import (
	"go-trades/entity"
	"log"

	"github.com/gin-gonic/gin"
)

func RBACMiddleware(requiredRoles ...entity.Role) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("role")
		if !exists {
			ctx.JSON(400, gin.H{"error": "role not found in context"})
			ctx.Abort()
			return
		}

		userRole, ok := role.(entity.Role)
		if !ok {
			ctx.JSON(403, gin.H{"error": "invalid role type"})
			ctx.Abort()
			return
		}

		log.Printf("UserRole: %v, RequiredRoles: %v\n", userRole, requiredRoles)

		for _, requiredRole := range requiredRoles {
			if userRole == requiredRole {
				ctx.Next()
				return
			}
		}

		ctx.JSON(403, gin.H{"error": "access denied"})
		ctx.Abort()
	}
}
