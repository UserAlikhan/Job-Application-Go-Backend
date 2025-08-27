package services

import (
	"database/sql"
	"job_portal/internal/models"
	"job_portal/internal/repository"
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
