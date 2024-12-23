package controllers

import (
	"bytes"
	"fmt"
	"io"
	"job-portal/models"
	"job-portal/services"
	"job-portal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	UserService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{UserService: userService}
}

func (uc *UserController) Register(c echo.Context) error {
	var user models.User

	// Read and log the raw request body for debugging
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read request body").SetInternal(err)
	}
	fmt.Println("Request Body:", string(body))

	// Reset the request body so it can be read again by c.Bind
	c.Request().Body = io.NopCloser(bytes.NewBuffer(body))

	// Bind the request body to the user struct
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input").SetInternal(err)
	}
	// Set default role if not provided
	if user.Role == "" {
		user.Role = "user"
	}
	// Validate the user struct
	if err := c.Validate(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Validation failed").SetInternal(err)
	}

	// Register the user
	err = uc.UserService.Register(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, "Failed to register user").SetInternal(err)
	}

	// Return success response with user details
	return utils.SendResponse(c, http.StatusCreated, "User registered successfully", user)
}

// Login handles user login and sets JWT in a cookie
func (uc *UserController) Login(c echo.Context) error {
	// Parse login credentials
	var credentials struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}
	if err := c.Bind(&credentials); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input").SetInternal(err)
	}
	if err := c.Validate(&credentials); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Validation failed").SetInternal(err)
	}

	// Authenticate the user and get token
	token, user, err := uc.UserService.Authenticate(credentials.Email, credentials.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid email or password").SetInternal(err)
	}

	// Set the JWT token in a secure, HttpOnly cookie
	c.SetCookie(&http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		Secure:   true, // Set to true if using HTTPS
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	// Return the response using SendResponse
	loginResponse := utils.CreateLoginResponse(token, user.ID.Hex(), user.Email, user.Role)
	return c.JSON(http.StatusOK, loginResponse)
}
