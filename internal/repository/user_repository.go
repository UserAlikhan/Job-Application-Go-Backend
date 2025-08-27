package repository

import (
	"database/sql"
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
