package services

import (
	"database/sql"
	"job_portal/internal/models"
	"job_portal/internal/repository"
	"job_portal/package/utils"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(db *sql.DB, user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	return repository.CreateUser(db, user)
}

func LoginUser(db *sql.DB, username string, password string) (string, error) {
	user, err := repository.GetUserByUsername(db, username)
	log.Println("User ", user)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}

	// access the utils package that generates a token
	return utils.GenerateToken(user.ID, user.Username, user.IsAdmin)
}

func ForgetPassword(db *sql.DB, username string) (string, error) {
	user, err := repository.GetUserByUsername(db, username)
	if err != nil {
		return "", err
	}

	generatedPassword := utils.GeneratePassword(6)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(generatedPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}

	user.Password = string(hashedPassword)

	if err := repository.UpdateUserPassword(db, user); err != nil {
		return "", nil
	}

	return generatedPassword, nil
}
