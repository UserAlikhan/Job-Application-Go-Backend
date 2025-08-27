package handlers

import (
	"database/sql"
	"job_portal/internal/models"
	"job_portal/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateJobHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var job models.Job

		if err := ctx.ShouldBindJSON(&job); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID := ctx.GetInt("userID")
		job.UserID = userID

		createdJob, err := services.CreateJob(db, &job)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, createdJob)
	}
}

func GetAllJobsHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jobs, err := services.GetAllJobs(db)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, jobs)
	}
}
