package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Book struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
}

var books = []Book{
	{ID: 1, Title: "Clean Code", Author: "Robert C. Martin"},
	{ID: 2, Title: "The Pragmatic Programmer", Author: "Andrew Hunt"},
	{ID: 3, Title: "Go in Action", Author: "William Kennedy"},
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/ping", ping)
	e.POST("/echo", echoFunction)
	// e.POST("/books", notImplemented)
	// e.GET("/books", getBooks)
	// e.GET("/books/:id", getBookByID)
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

func echoFunction(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	if !json.Valid(body) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid json"})
	}
	return c.Blob(http.StatusOK, "application/json", body)
}	


// func getBooks(c echo.Context) error {

// 	return c.JSON(http.StatusOK, map[string]bool{
// 		"success": true,
// 	})
// }	
