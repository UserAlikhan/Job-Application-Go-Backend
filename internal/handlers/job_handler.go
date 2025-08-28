package handlers

import (
	"database/sql"
	"job_portal/internal/models"
	"job_portal/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateJobHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var job models.Job

		if err := ctx.ShouldBindJSON(&job); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

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

func GetJobsByUserIDHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, err := strconv.Atoi(ctx.Param("userID"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "userID was not provided"})
			return
		}

		jobs, err := services.GetJobsByUserID(db, userID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, jobs)
	}
}

func GetJobByIDHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		job, err := services.GetJobByID(db, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, job)
	}
}

func UpdateJobHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get id parameter from url
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// variable to store job
		var job models.Job
		job.ID = id
		// Store body parmeters in job variable
		if err := ctx.ShouldBindJSON(&job); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		// getting these variables from middleware
		userID := ctx.GetInt("userID")
		isAdmin := ctx.GetBool("isAdmin")
		// calling service
		if _, err := services.UpdateJob(db, &job, userID, isAdmin); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// printing response if everything is successfull
		ctx.JSON(http.StatusOK, job)
	}
}

func DeleteJobHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		userID := ctx.GetInt("userID")
		isAdmin := ctx.GetBool("isAdmin")

		if err := services.DeleteJob(db, id, userID, isAdmin); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Job was deleted successfully"})
	}
}
