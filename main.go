package main

import (
	"todo_app/controllers"
	"todo_app/models"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	models.ConnectDatabase()

	// Fetch all books
	r.GET("/books", controllers.FindBooks)
	// Add a new book
	r.POST("/books", controllers.CreateBook)
	// Get a single book by ID
	r.GET("/books/:id", controllers.FindBook)
	// Update a book
	r.PATCH("/books/:id", controllers.UpdateBook)
	// Delete a book
	r.DELETE("/books/:id", controllers.DeleteBook)

	r.Run()
}
