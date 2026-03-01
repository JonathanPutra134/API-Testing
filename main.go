package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Book struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
	Year  int    `json:"year"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var books = []Book{
	// {ID: 1, Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Year: 1925},
	// {ID: 2, Title: "To Kill a Mockingbird", Author: "Harper Lee", Year: 1960},
	// {ID: 3, Title: "1984", Author: "George Orwell", Year: 1949},
}
var nextID = 1
const authToken = "generated-token"

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
	auth := c.Request().Header.Get("Authorization")

	if auth != "Bearer "+authToken {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid or missing token",
		})
	}

	author := c.QueryParam("author")
	fmt.Println("Author filter:", author)
	if author != "" {
		filteredBooks := []Book{}
		for _, book := range books {
			fmt.Println(book.Author)
			fmt.Println(author)
			if book.Author == author {
				fmt.Println("Matched book:", book.Title)
				filteredBooks = append(filteredBooks, book)
			}
		}
		return c.JSON(http.StatusOK, filteredBooks)
	}

	page := c.QueryParam("page")
	limit := c.QueryParam("limit")

	if page != "" && limit != "" {
		pageNum, err1 := strconv.Atoi(page)
		limitNum, err2 := strconv.Atoi(limit)
		if err1 != nil || err2 != nil || pageNum < 1 || limitNum < 1 {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "page and limit must be positive integers",
			})
		}
		fmt.Println("Page:", pageNum, "Limit:", limitNum)
		
		start := (pageNum - 1) * limitNum
		end := start + limitNum

		fmt.Println("Start:", start, "End:", end)
		if start >= len(books) {
			return c.JSON(http.StatusOK, []Book{})
		}
		if end > len(books) {
			end = len(books)
		}
		return c.JSON(http.StatusOK, books[start:end])
	}
	return c.JSON(http.StatusOK, books)
}	

func createBooks(c echo.Context) error {
	var payload Book
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload bos"})
	}

	payload.Title = strings.TrimSpace(payload.Title)
	payload.Author = strings.TrimSpace(payload.Author)

	if payload.Title == "" || payload.Author == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "title and author are required",
		})
	}

	if payload.Year == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "year must be > 0",
		})
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
	var loginReq LoginRequest
	if err := c.Bind(&loginReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid json body",
		})
	}

   if loginReq.Username != "admin" || loginReq.Password != "password" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"token": authToken,
	})
}	

