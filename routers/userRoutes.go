package routers

import (
	"job-portal/controllers"

	"github.com/labstack/echo/v4"
)

func RegisterUserRoutes(e *echo.Echo, userController *controllers.UserController) {
	e.POST("/register", userController.Register)
	e.POST("/login", userController.Login)
}
