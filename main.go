package main

import (
	// "fmt"
	"net/http"
	
	"github.com/gin-gonic/gin"

	"github.com/Sarthak-Neupane/Clothing_Backend.git/Crawler"

)


func handleQuery(context *gin.Context) {
	response := crawler.CrawlHnM()
	context.IndentedJSON(http.StatusCreated, response)
}



func main() {
	router := gin.Default()
	router.POST("/search", handleQuery)
	
	router.Run("localhost:8080")
}