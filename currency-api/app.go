package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const API string = "https://api.exchangerate-api.com/v4/latest/USD"

type responseStructure map[string]any

func main() {
	// Load environment variables from .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("No .env file found or error loading .env file")
	}

	HOST := os.Getenv("CURRENCY_API_HOST")
	PORT := os.Getenv("CURRENCY_API_PORT")

	server := gin.Default()

	server.GET(
		"/service/get-today/exchange-rate",
		func(c *gin.Context) {

			response, err := http.Get(API)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			defer response.Body.Close()
			responseData, err := io.ReadAll(response.Body)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			var responseStructure responseStructure
			if err := json.Unmarshal(responseData, &responseStructure); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid JSON from API"})
				return
			}

			rates, ok := responseStructure["rates"]

			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Try Later"})
				return
			}

			value, ok := rates.(map[string]any)

			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Try Later"})
				return
			}
			c.JSON(
				http.StatusOK,
				gin.H{
					"1 USD": fmt.Sprint(value["EGP"], " EGP"),
				})

		})

	server.Run(HOST + ":" + PORT)

}
