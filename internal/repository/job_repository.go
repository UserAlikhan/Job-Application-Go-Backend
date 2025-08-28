package repository

import (
	"database/sql"
	"job_portal/internal/models"
)

func CreateJob(db *sql.DB, job *models.Job) (*models.Job, error) {
	// check if user with given Id exists
	user, err := GetUserByID(db, job.UserID)
	if err != nil || user == nil {
		return nil, err
	}

	result, err := db.Exec(
		`INSERT INTO jobs (title, description, location, company, salary, user_id) VALUES (?, ?, ?, ?, ?, ?)`,
		job.Title, job.Description, job.Location, job.Company, job.Salary, job.UserID,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	job.ID = int(id)

	return job, err
}

func GetAllJobs(db *sql.DB) ([]*models.Job, error) {
	rows, err := db.Query("SELECT * FROM jobs")
	if err != nil {
		return nil, err
	}

	var jobs []*models.Job

	for rows.Next() {
		job := &models.Job{}
		if err := rows.Scan(
			&job.ID, &job.Title, &job.Description,
			&job.Location, &job.Company, &job.Salary,
			&job.UserID, &job.CreatedAt, &job.UpdatedAt,
		); err != nil {
			return nil, err
		}

		// second query to fetch the user data
		user, err := GetUserByID(db, job.UserID)
		if err == nil {
			job.User = user
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func GetJobsByUserID(db *sql.DB, userID int) ([]*models.Job, error) {
	rows, err := db.Query("SELECT * FROM jobs WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}

	var jobs []*models.Job

	for rows.Next() {
		job := &models.Job{}

		if err := rows.Scan(
			&job.ID, &job.Title, &job.Description,
			&job.Location, &job.Company, &job.Salary,
			&job.UserID, &job.CreatedAt, &job.UpdatedAt,
		); err != nil {
			return nil, err
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func GetJobByID(db *sql.DB, id int) (*models.Job, error) {
	job := &models.Job{}

	err := db.QueryRow(
		"SELECT * FROM jobs WHERE id = ?",
		id,
	).Scan(
		&job.ID, &job.Title, &job.Description, &job.Location,
		&job.Company, &job.Salary, &job.UserID,
		&job.CreatedAt, &job.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func UpdateJob(db *sql.DB, job *models.Job) (*models.Job, error) {
	_, err := db.Exec(`
			UPDATE jobs SET title = ?, description = ?, location = ?, company = ?, salary = ?, user_id = ? WHERE id = ?
		`, job.Title, job.Description, job.Location, job.Company, job.Salary, job.UserID, job.ID,
	)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func DeleteJob(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM jobs WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}
