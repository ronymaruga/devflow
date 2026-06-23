package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
)



type book struct{
	ID string `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
	Quantity int `json:"quantity"`
}


var books = []book{
	{ID: "1", Title: "The Hobbit", Author: "J.R.R. Tolkien", Quantity: 12},
	{ID: "2", Title: "The Great Gatsby", Author: "Scot Fitzgerald", Quantity: 7},
	{ID: "3", Title: "War and Peace", Author: "Leo", Quantity: 6},
}


func getBooks (c *gin.Context){
	c.IndentedJSON(http.StatusOK, books)
}

//gin.Context - this stores all information related to a request eg. query parameters, payload, headers
func createBook (c *gin.Context){
	var newBook book

	err := c.BindJSON(&newBook)

	if err != nil{
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)

}

//error can be nil incase the book is found
func getBookById(id string) (*book, error){
	for i, b := range books {
		if b.ID == id{
			return &books[i], nil
		}

	}
	return nil, errors.New("book not found")

}


func bookById(c *gin.Context){
	id := c.Param("id")

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, book)

}

func checkoutBooks(c *gin.Context){
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message":"Missing parameter"})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message":"Book not found"})
		return
	}

	if book.Quantity <= 0{
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "no books available"})
		return
	}

	book.Quantity -= 1

	c.IndentedJSON(http.StatusOK, book)


}


func returnBook (c *gin.Context){
	id, ok := c.GetQuery("id")

	if !ok{
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing parameter"})
		return
	}

	book, err := getBookById(id)

	if err != nil{
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "this book doesn't exist in our collection"})
		return
	}

	book.Quantity += 1

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Book has been received"})
}


func main(){
	router := gin.Default()

	router.GET("/books", getBooks)

	// curl localhost:8080/books --include --header "Content-Type: application/json" -d @body.json --request "POST"
	router.POST("/books", createBook)

	// the colon : is the parameter
	router.GET("/books/:id", bookById)

	// curl localhost:8080/checkout?id=3 --request "PATCH"
	router.PATCH("/checkout", checkoutBooks)

	router.PATCH("/return", returnBook)

	router.Run("localhost:8080")


}