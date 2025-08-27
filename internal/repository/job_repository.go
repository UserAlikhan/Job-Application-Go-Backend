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
			&job.Title, &job.Description, &job.Location,
			&job.Company, &job.Salary, &job.UserID,
			&job.CreatedAt, &job.UpdatedAt,
		); err != nil {
			return nil, err
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}
