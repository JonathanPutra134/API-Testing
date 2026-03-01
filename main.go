package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Global middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// API routes template
	api := e.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/hello", helloHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

func helloHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "hello world",
	})
}