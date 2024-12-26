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
	e.Use(middleware.Logger())    // Log HTTP requests
	e.Use(middleware.Recover())   // Recover from panics
	e.Use(middleware.RequestID()) // Add a unique request ID for each request
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{ "https://job-portal-frontend-pink.vercel.app"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // Enable credentials
	})) // Configure CORS
	// e.Use(middleware.Secure())                                              // Add secure headers (e.g., X-Frame-Options, HSTS)
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20))) // Rate limit: 20 requests per second
	e.Use(middleware.BodyLimit("2M"))                                       // Limit request body size to 2MB

	// Validator
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
