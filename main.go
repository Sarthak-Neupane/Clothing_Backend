package main

import (
	// "fmt"

	// "github.com/gocolly/colly/v2"

	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Sarthak-Neupane/Clothing_Backend.git/Crawler"
	// "encoding/json"
)

type queries struct {
	Page       string `json:"page"`
	PageId     string `json:"pageId"`
	PageSize   string `json:"pageSize"`
	CategoryId string `json:"categoryId"`
}

func handleQuery(context *gin.Context) {
	var newQuery queries

	if err := context.BindJSON(&newQuery); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return
	}

	queryParams := map[string]string{
		"page":       newQuery.Page,
		"pageId":     newQuery.PageId,
		"page-size":  newQuery.PageSize,
		"categoryId": newQuery.CategoryId,
	}

	response, err := crawler.CrawlHnM(queryParams)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
	}
	context.IndentedJSON(http.StatusCreated, response)
}

func main() {
	router := gin.Default()
	router.POST("/search", handleQuery)


	router.Run("localhost:8080")
}
