package routes

import (
	"database/sql"
	"job_portal/internal/handlers"
	"job_portal/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine, db *sql.DB) {
	{
		// AUTH ROUTES
		r.POST("/login", handlers.LoginHandler(db))
		r.POST("/register", handlers.RegisterHandler(db))

		// USER ROUTES // employer
		authenticated := r.Group("/")
		// Apply AuthMiddleware
		authenticated.Use(middlewares.AuthMiddleware())
		authenticated.GET("/users/:id", handlers.GetUserByIdHandler(db))
		authenticated.PATCH("/users/:id", handlers.UpdateUserProfileHandler(db))
		authenticated.PATCH("/users/:id/picture", handlers.UpdateUserProfilePicture(db))
		//
		r.GET("/users", handlers.GetAllUsersHandler(db))

		// Job Routes
		authenticated.POST("/jobs", handlers.CreateJobHandler(db))
		r.GET("/jobs", handlers.GetAllJobsHandler(db))
	}
}
