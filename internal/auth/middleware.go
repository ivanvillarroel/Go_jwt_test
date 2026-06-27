package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Middleware(tokenService TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "bearer token is required"})
			return
		}

		claims, err := tokenService.Validate(parts[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		ctx.Set("subject", claims.Subject)
		ctx.Set("permissions", claims.Permissions)
		ctx.Set("claims", claims)
		ctx.Next()
	}
}

func RequirePermission(permission string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		value, exists := ctx.Get("claims")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token claims are required"})
			return
		}

		claims, ok := value.(*Claims)
		if !ok || !claims.HasPermission(permission) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		ctx.Next()
	}
}
