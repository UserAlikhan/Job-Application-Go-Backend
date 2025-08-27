package handlers

import (
	"database/sql"
	"job_portal/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllUsersHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		users, err := services.GetAllUsers(db)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, users)
	}
}

func GetUserByIdHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		user, err := services.GetUserByID(db, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, user)
	}
}

func UpdateUserProfileHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var userUpdate struct {
			Username string `json:"username"`
			Email    string `json:"email"`
		}

		if err := ctx.ShouldBindJSON(&userUpdate); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// userId and isAdmin parameters was set in middleware
		userID := ctx.GetInt("userID")
		isAdmin := ctx.GetBool("isAdmin")
		// If user is not the admin he cannot
		// update someone's profile
		if !isAdmin && userID != id {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to update this user profile"})
			return
		}

		updatedUser, err := services.UpdateUserProfile(db, id, userUpdate.Username, userUpdate.Email)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		ctx.JSON(http.StatusOK, updatedUser)
	}
}
