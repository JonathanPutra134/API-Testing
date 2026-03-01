package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/ping", ping)
	// e.POST("/echo", notImplemented)
	// e.POST("/books", notImplemented)
	// e.GET("/books", notImplemented)
	// e.GET("/books/:id", notImplemented)
	// e.PUT("/books/:id", notImplemented)
	// e.DELETE("/books/:id", notImplemented)
	// e.POST("/auth/token", notImplemented)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
}		

func ping(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]bool{
		"success": true,
	})
}	
