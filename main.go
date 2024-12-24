package main

import (
	"fmt"
	"net/http"
	
	"github.com/gin-gonic/gin"

)

type product struct {
	Query string `json:"query"`
	Type string `json:"type"`
}

var products = []product{}

func handleQuery(context *gin.Context) {
	var newProduct product
	if err := context.BindJSON(&newProduct); err != nil {
		fmt.Println(err)
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return
	}
	products = append(products, newProduct)
	context.IndentedJSON(http.StatusCreated, products)
}

func main() {
	// fmt.Println("Hello World")
	router := gin.Default()
	router.POST("/search", handleQuery)

	router.Run("localhost:8080")
}