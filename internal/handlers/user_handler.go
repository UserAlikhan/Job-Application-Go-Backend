package handlers

import (
	"database/sql"
	"fmt"
	"job_portal/internal/models"
	"job_portal/internal/services"
	"net/http"
	"os"
	"path/filepath"
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
			return
		}

		ctx.JSON(http.StatusOK, updatedUser)
	}
}

func UpdateUserProfilePicture(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invlaid ID"})
			return
		}
		// these values were set in middleware
		userID := ctx.GetInt("userID")
		isAdmin := ctx.GetBool("isAdmn")

		if !isAdmin && userID != id {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to update this profile"})
			return
		}
		// ctx.FormFile extracts the uploaded file from the form field
		// works with Works with multipart/form-data
		file, err := ctx.FormFile("profile_picture")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error upload the file"})
			return
		}
		// os.MkdirAll creates the specified directory
		if err := os.MkdirAll(os.Getenv("UPLOAD_DIR"), os.ModePerm); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creading upload directory"})
			return
		}
		// constructing the file name (eg 5-profilepic.png)
		filename := fmt.Sprintf("%d-%s", id, filepath.Base(file.Filename))
		filePath := filepath.Join(os.Getenv("UPLOAD_DIR"), filename)

		if err := ctx.SaveUploadedFile(file, filePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error uploading picture"})
			return
		}

		// call the service method to upload the profile picture
		err = services.UploadProfilePicture(db, id, filename)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error uploading picture"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Profile picture was uploaded successfully"})
	}
}

func DeleteUserHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Check if user is Admin
		isAdmin := ctx.GetBool("isAdmin")
		if !isAdmin {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to delete user"})
			return
		}
		// Check if entered id is correct
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
			return
		}
		// Check if user is trying to delete himself
		currentUserID := ctx.GetInt("userID")
		if currentUserID == id {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "User cannot delete himself / herself"})
			return
		}

		if err := services.DeleteUser(ctx, db, id); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "User was deleted successfully"})
	}
}

func ChangePasswordHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req models.ChangePassword
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID := ctx.GetInt("userID")
		err := services.ChangePassword(db, userID, req.CurrentPassword, req.NewPassword)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Password was changed successfully"})
	}
}
