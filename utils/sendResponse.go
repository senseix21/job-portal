package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// LoginResponse defines the structure of the login response
type LoginResponse struct {
	Success bool   `json:"success"`
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Token string `json:"token"`
		User  struct {
			ID    string `json:"id"`
			Email string `json:"email"`
			Role  string `json:"role"`
		} `json:"user"`
	} `json:"data"`
}

// SendResponse sends a generic JSON response
func SendResponse(c echo.Context, statusCode int, message string, data interface{}) error {
	response := map[string]interface{}{
		"success": statusCode == http.StatusOK,
		"status":  statusCode,
		"message": message,
		"data":    data,
	}
	return c.JSON(statusCode, response)
}

// CreateLoginResponse creates a LoginResponse object
func CreateLoginResponse(token string, userID, email, role string) LoginResponse {
	var response LoginResponse
	response.Success = true
	response.Status = http.StatusOK
	response.Message = "Login successful"
	response.Data.Token = token
	response.Data.User.ID = userID
	response.Data.User.Email = email
	response.Data.User.Role = role
	return response
}
