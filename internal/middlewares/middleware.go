package middlewares

import (
	"job_portal/package/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Missing Authorization header"})
			ctx.Abort()
			return
		}

		claims, err := utils.ValidateToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Invalid Token"})
			ctx.Abort()
			return
		}

		ctx.Set("userID", claims.ID)
		ctx.Set("isAdmin", claims.IsAdmin)

		ctx.Next()
	}
}
