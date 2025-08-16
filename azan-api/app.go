package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var API string = "https://api.aladhan.com/v1/timingsByCity/" + time.Now().Format(time.DateOnly) + "?city=Cairo&country=EGY&method=5"

type responseStructure map[string]any

func main() {
	// Load environment variables from .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("No .env file found or error loading .env file")
	}

	HOST := os.Getenv("AZAN_API_HOST")
	PORT := os.Getenv("AZAN_API_PORT")

	server := gin.Default()

	server.GET(
		"/service/get-today/azan-times",
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

			data, ok := responseStructure["data"]

			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Try Later"})
				return
			}

			value, ok := data.(map[string]any)

			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Try Later"})
				return
			}

			timings := value["timings"]

			timingsValue, ok := timings.(map[string]any)

			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Try Later"})
				return
			}

			c.JSON(
				http.StatusOK,
				timingsValue)
		})

	server.Run(HOST + ":" + PORT)

}
