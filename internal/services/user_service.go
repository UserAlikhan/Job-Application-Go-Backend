package services

import (
	"database/sql"
	"fmt"
	"job_portal/internal/models"
	"job_portal/internal/repository"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetUserByID(db *sql.DB, id int) (*models.User, error) {
	return repository.GetUserByID(db, id)
}

func GetAllUsers(db *sql.DB) ([]*models.User, error) {
	return repository.GetAllUsers(db)
}

func UpdateUserProfile(db *sql.DB, id int, username string, email string) (*models.User, error) {
	user := &models.User{
		ID:       id,
		Username: username,
		Email:    email,
	}

	return repository.UpdateUserProfile(db, user)
}

func UploadProfilePicture(db *sql.DB, id int, profilePicture string) error {
	return repository.UploadProfilePicture(db, id, profilePicture)
}

func DeleteUser(ctx *gin.Context, db *sql.DB, id int) error {
	// We need to delete all associated with users data from
	// another tables
	tx, err := db.BeginTx(ctx.Request.Context(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return fmt.Errorf("Error starting transaction: %v", err)
	}
	// Rollback if not commited
	defer tx.Rollback()

	err = repository.DeleteUserWithTransaction(tx, id)
	if err != nil {
		return fmt.Errorf("Error deleting user: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("Error commiting transaction: %v", err)
	}

	return nil
}

func ChangePassword(db *sql.DB, userID int, currentPassword string, newPassword string) error {
	// get user's credentials by userID
	user, err := GetUserByID(db, userID)
	if err != nil {
		return err
	}
	// compare saved and entered by user passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return fmt.Errorf("Incorrect password %v", err)
	}

	// hash user's new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	return repository.ChangePassword(db, userID, currentPassword, string(hashedPassword))
}
