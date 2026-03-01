package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Book struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
	Year  int    `json:"year"`
}

var books = []Book{
}
var nextID = 1
var token string

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/ping", ping)
	e.POST("/echo", echoFunction)
	e.POST("/books", createBooks)
	e.GET("/books", getBooks)
	e.GET("/books/:id", getBookByID)
	e.PUT("/books/:id", updateBook)
	e.DELETE("/books/:id", deleteBook)
	e.POST("/auth/token", generateAuthToken)

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


func getBooks(c echo.Context) error {
	// if(token !== "") {
	// 	return c.JSON(http.StatusUnauthorized, map[string]string{
	// 		"error": "missing auth token",
	// 	})
	// }
	return c.JSON(http.StatusOK, books)
}	

func createBooks(c echo.Context) error {
	var payload Book
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload bos"})
	}
	payload.ID = nextID
	nextID++
	
	books = append(books, payload)
	return c.JSON(http.StatusCreated, payload)
}	

func getBookByID(c echo.Context) error {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "id must be numeric",
		})
	}

	for _, book := range books {
		if book.ID == id {
			return c.JSON(http.StatusOK, book)
		}
	}

	return c.JSON(http.StatusNotFound, map[string]string{
		"error": "book not found",
	})
}	


func updateBook(c echo.Context) error {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "id must be numeric",
		})
	}

	var payload Book

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid json body",
		})
	}

	for i := range books {
		if books[i].ID == id {
			books[i].Title = payload.Title
			books[i].Author = payload.Author
			books[i].Year = payload.Year
			return c.JSON(http.StatusOK, books[i])
		}
	}

	return c.JSON(http.StatusNotFound, map[string]string{
		"error": "book not found",
	})
}	

func deleteBook(c echo.Context) error {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "id must be numeric",
		})
	}

	for i, book := range books {
		if book.ID == id {
			books = append(books[:i], books[i+1:]...)
			return c.JSON(http.StatusOK, map[string]string{
				"message": "book deleted",
			})
		}
	}

	return c.JSON(http.StatusNotFound, map[string]string{
		"error": "book not found",
	})
}	


func generateAuthToken(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"token": "generated-token",
	})
}	

