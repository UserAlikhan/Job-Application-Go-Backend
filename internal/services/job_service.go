package services

import (
	"database/sql"
	"errors"
	"job_portal/internal/models"
	"job_portal/internal/repository"
)

func CreateJob(db *sql.DB, job *models.Job) (*models.Job, error) {
	return repository.CreateJob(db, job)
}

func GetAllJobs(db *sql.DB) ([]*models.Job, error) {
	return repository.GetAllJobs(db)
}

func GetJobsByUserID(db *sql.DB, userID int) ([]*models.Job, error) {
	return repository.GetJobsByUserID(db, userID)
}

func GetJobByID(db *sql.DB, id int) (*models.Job, error) {
	return repository.GetJobByID(db, id)
}

func UpdateJob(db *sql.DB, job *models.Job, userID int, isAdmin bool) (*models.Job, error) {
	// check if job exists
	existingJob, err := repository.GetJobByID(db, job.ID)
	if err != nil {
		return nil, err
	}
	// only admin or job owner itself could change info
	if !isAdmin && existingJob.UserID != userID {
		return nil, errors.New("Unauthorized to update this job")
	}
	// calling repository
	return repository.UpdateJob(db, job)
}

func DeleteJob(db *sql.DB, id int, userID int, isAdmin bool) error {
	// check if job exists
	existingJob, err := repository.GetJobByID(db, id)
	if err != nil {
		return err
	}

	if !isAdmin && userID != existingJob.UserID {
		return errors.New("Unauthorized to update this job")
	}

	return repository.DeleteJob(db, id)
}
