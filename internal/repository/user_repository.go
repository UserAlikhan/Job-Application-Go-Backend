package repository

import (
	"database/sql"
	"fmt"
	"job_portal/internal/models"
)

func CreateUser(db *sql.DB, user *models.User) error {
	_, err := db.Exec(
		`INSERT INTO users (username, password, email) VALUES (?, ?, ?)`,
		user.Username, user.Password, user.Email,
	)

	return err
}

func GetAllUsers(db *sql.DB) ([]*models.User, error) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	// close the DB pool connection
	defer rows.Close()

	// Array for storing all the users
	var users []*models.User
	// iterating over all rows
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(
			&user.ID, &user.Username, &user.Password,
			&user.Email, &user.IsAdmin, &user.ProfilePicture,
			&user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		// If no error appen the array
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetUserByUsername(db *sql.DB, username string) (*models.User, error) {
	user := &models.User{}

	err := db.QueryRow(
		"SELECT * FROM users WHERE username = ?",
		username,
	).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email, &user.IsAdmin,
		&user.ProfilePicture, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByID(db *sql.DB, id int) (*models.User, error) {
	var user models.User
	// use sql.NullString to handle the null value
	var profilePicture sql.NullString

	err := db.QueryRow(
		`SELECT * FROM users WHERE id = ?`,
		id,
	).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email, &user.IsAdmin,
		&profilePicture, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if profilePicture.Valid {
		user.ProfilePicture = &profilePicture.String
	} else {
		user.ProfilePicture = nil
	}

	return &user, nil
}

func UpdateUserProfile(db *sql.DB, user *models.User) (*models.User, error) {
	_, err := db.Exec(
		`UPDATE users SET username = ?, email = ? WHERE id = ?`,
		user.Username, user.Email, user.ID,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func UploadProfilePicture(db *sql.DB, id int, profilePicture string) error {
	_, err := db.Exec("UPDATE users SET profile_picture = ? WHERE id = ?", profilePicture, id)
	if err != nil {
		return err
	}

	return nil
}

func UpdateUserPassword(db *sql.DB, user *models.User) error {
	_, err := db.Exec(
		`UPDATE users SET password = ? WHERE id = ?`,
		user.Password, user.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUserWithTransaction(tx *sql.Tx, userID int) error {
	// Delete associated jobs
	_, err := tx.Exec("DELETE FROM jobs WHERE user_id = ?", userID)
	if err != nil {
		return fmt.Errorf("Error deleting user's jobs: %v", err)
	}

	// Delete the user
	result, err := tx.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func ChangePassword(db *sql.DB, userID int, currentPassword string, hashedNewPassword string) error {
	result, err := db.Exec(`UPDATE users SET password = ? WHERE id = ?`, hashedNewPassword, userID)
	if err != nil {
		return err
	}

	// check if update actually affected a row
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error checking update result: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No user found with id %d", userID)
	}

	return nil
}
