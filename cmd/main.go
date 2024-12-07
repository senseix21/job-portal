package main

import (
	"job-portal/config"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config.Connect()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Define the root route
	e.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, "Hello") })
	e.Logger.Fatal(e.Start(":8080"))
}
