package middlewares

import (
	"bytes"
	"io"
	"job_portal/internal/models"
	"job_portal/package/utils"

	"github.com/gin-gonic/gin"
)

func PasswordValidationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Read the request body
		bodyBytes, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "Error reading request body"})
			ctx.Abort()
			return
		}

		// create new reader with the bytes
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Parse the request
		var req models.ChangePassword
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request body"})
			ctx.Abort()
			return
		}

		// Restore the request body for the next middleware
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		isValid, errors := utils.ValidatePasswordStrength(req.NewPassword)
		if !isValid {
			ctx.JSON(400, gin.H{
				"error":   "Password validation failed",
				"details": errors,
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
