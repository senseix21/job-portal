package main

import (
	"job-portal/config"
	"job-portal/controllers"
	"job-portal/middlewares"
	"job-portal/routers"
	"job-portal/services"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	// Initialize Echo
	e := echo.New()

	// Set the custom error handler
	e.HTTPErrorHandler = middlewares.CustomHTTPErrorHandler

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Validator = &CustomValidator{validator: validator.New()}

	// Connect to the database
	config.Connect()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from golang server!")
	})

	// Initialize user service and controller
	userCollection := config.GetCollection("jobportal", "users")
	userService := services.NewUserService(userCollection)
	userController := controllers.NewUserController(userService)

	// Initialize job service and controller
	jobCollection := config.GetCollection("jobportal", "jobs")
	jobService := services.NewJobService(jobCollection)
	jobController := controllers.NewJobController(jobService)

	// Register routes
	routers.RegisterUserRoutes(e, userController)
	routers.RegisterJobRoutes(e, jobController)

	// Start the server
	log.Fatal(e.Start(":8080"))
}
