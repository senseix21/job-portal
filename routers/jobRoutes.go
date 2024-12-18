package routers

import (
	"job-portal/controllers"
	"job-portal/middlewares"

	"github.com/labstack/echo/v4"
)

func RegisterJobRoutes(e *echo.Echo, jobController *controllers.JobController) {
	jobGroup := e.Group("/jobs")

	jobGroup.POST("/create", jobController.CreateJobHandler, middlewares.JWTMiddleware("user", "admin"))
	jobGroup.GET("/", jobController.ListJobsHandler, middlewares.JWTMiddleware("user", "admin"))        // Get all jobs
	jobGroup.GET("/:id", jobController.GetJobHandler, middlewares.JWTMiddleware("user", "admin"))       // Get a job by ID
	jobGroup.PATCH("/:id", jobController.UpdateJobHandler, middlewares.JWTMiddleware("user", "admin"))  // Update a job by ID
	jobGroup.DELETE("/:id", jobController.DeleteJobHandler, middlewares.JWTMiddleware("user", "admin")) // Delete a job by ID
}
