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
		r.POST("/forgetPassword", handlers.ForgetPasswordHandler(db))

		// USER ROUTES (Apply AuthMiddleware)
		authenticated := r.Group("/")
		authenticated.Use(middlewares.AuthMiddleware())
		authenticated.GET("/users/:id", handlers.GetUserByIdHandler(db))
		authenticated.PATCH("/users/:id", handlers.UpdateUserProfileHandler(db))
		authenticated.PATCH("/users/:id/picture", handlers.UpdateUserProfilePicture(db))
		authenticated.DELETE("/users/:id", handlers.DeleteUserHandler(db))
		authenticated.PUT("users/change-password", middlewares.PasswordValidationMiddleware(), handlers.ChangePasswordHandler(db))
		//
		r.GET("/users", handlers.GetAllUsersHandler(db))

		// Job Routes
		authenticated.POST("/jobs", handlers.CreateJobHandler(db))
		authenticated.PUT("/jobs/:id", handlers.UpdateJobHandler(db))
		authenticated.DELETE("/jobs/:id", handlers.DeleteJobHandler(db))
		//
		r.GET("/jobs", handlers.GetAllJobsHandler(db))
		r.GET("/jobs/usersJobs/:userID", handlers.GetJobsByUserIDHandler(db))
		r.GET("/jobs/:id", handlers.GetJobByIDHandler(db))
	}
}
