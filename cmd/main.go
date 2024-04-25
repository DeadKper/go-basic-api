package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

func (_ Book) new(id string, title string, author string, quantity int) Book {
	return Book{ID: id, Title: title, Author: author, Quantity: quantity}
}

var books = []Book{
	Book{}.new("1", "In Search of a Lost Time", "Marcel Proust", 2),
	Book{}.new("2", "The Great Gatsby", "F. Scott Fitzgerald", 5),
	Book{}.new("3", "War and Peace", "Leo Tolstoy", 6),
}

func getBooks(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, books)
}

func addBook(context *gin.Context) {
	var book Book

	if err := context.BindJSON(&book); err != nil {
		return
	}

	books = append(books, book)
	context.IndentedJSON(http.StatusCreated, book)
}

func bookById(context *gin.Context) {
	var id = context.Param("id")
	var book, err = getBookById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found!"})
		return
	}

	context.IndentedJSON(http.StatusOK, book)
}

func getBookById(id string) (*Book, error) {
	for i, book := range books {
		if book.ID == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("book not found")
}

func checkoutBook(context *gin.Context) {
	var id, ok = context.GetQuery("id")

	if !ok {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter!"})
		return
	}

	var book, err = getBookById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found!"})
		return
	}

	if book.Quantity <= 0 {
		context.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "Book not available!"})
		return
	}

	book.Quantity--
	context.IndentedJSON(http.StatusOK, book)
}

func returnBook(context *gin.Context) {
	var id, ok = context.GetQuery("id")

	if !ok {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter!"})
		return
	}

	var book, err = getBookById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found!"})
		return
	}

	book.Quantity++
	context.IndentedJSON(http.StatusOK, book)
}

func removeBook(context *gin.Context) {
	var id, ok = context.GetQuery("id")

	if !ok {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter!"})
		return
	}

	var to_delete Book
	var found = false

	for i, book := range books {
		if book.ID == id {
			to_delete = book
			found = true
			books = append(books[:i], books[i+1:]...)
		}
	}

	if !found {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found!"})
		return
	}

	context.IndentedJSON(http.StatusOK, to_delete)
}

func main() {
	var router = gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", bookById)
	router.POST("/books", addBook)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.DELETE("/remove", removeBook)
	router.Run("localhost:42069")
}
