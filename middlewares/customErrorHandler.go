package middlewares

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// CustomHTTPErrorHandler handles errors globally for the application
func CustomHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	code := http.StatusInternalServerError
	message := "Internal Server Error"
	var errors interface{} = nil

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		if he.Message != nil {
			message = fmt.Sprintf("%v", he.Message)
		}
		if he.Internal != nil {
			errors = he.Internal.Error()
		}
	}

	c.Logger().Error(err)

	if c.Request().Header.Get("Accept") == "application/json" || c.Request().Header.Get("Content-Type") == "application/json" {
		response := map[string]interface{}{
			"success": false,
			"status":  code,
			"message": message,
			"errors":  errors,
		}
		if err := c.JSON(code, response); err != nil {
			c.Logger().Error(err)
		}
		return
	}

	errorPage := fmt.Sprintf("errors/%d.html", code)
	if err := c.File(errorPage); err != nil {
		c.Logger().Error(err)
		if fallbackErr := c.String(code, fmt.Sprintf("Error %d: %s", code, message)); fallbackErr != nil {
			c.Logger().Error(fallbackErr)
		}
	}
}
