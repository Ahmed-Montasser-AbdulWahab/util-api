package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type responseStructure map[string]string

func main() {
	// Load environment variables from .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("No .env file found or error loading .env file")
	}

	//CURRENCY_API_HOST := os.Getenv("CURRENCY_API_HOST")
	CURRENCY_API_PORT := os.Getenv("CURRENCY_API_PORT")

	AZAN_API_HOST := os.Getenv("AZAN_API_SERVICE_HOST")
	AZAN_API_PORT := os.Getenv("AZAN_API_SERVICE_PORT")

	UTIL_API_HOST := os.Getenv("UTIL_API_HOST")
	UTIL_API_PORT := os.Getenv("UTIL_API_PORT")

	var AZAN_API string = "http://" + AZAN_API_HOST + AZAN_API_PORT + "/service/get-today/azan-times"
	var CURRENCY_API string = "http://localhost:" + CURRENCY_API_PORT + "/service/get-today/exchange-rate"
	AZAN_API = "http://" + "localhost:6000" + "/service/get-today/azan-times"
	CURRENCY_API = "http://localhost:" + "5000" + "/service/get-today/exchange-rate"
	server := gin.Default()

	server.GET(
		"/services/:serviceName",
		func(c *gin.Context) {
			var API_CALLED = ""
			serviceName := c.Param("serviceName")
			switch serviceName {
			case "1":
				API_CALLED = CURRENCY_API
			case "2":
				API_CALLED = AZAN_API
			default:
				c.JSON(
					http.StatusNotFound, gin.H{"error": "We don't offer this service"},
				)
			}
			response, err := http.Get(API_CALLED)

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

			c.JSON(
				http.StatusOK,
				responseStructure)

		})

	server.GET(
		"/services",
		func(c *gin.Context) {
			c.JSON(
				http.StatusOK,
				gin.H{
					"1": "Get EGP/USD Today",
					"2": "Get Azan Timings today in Cairo, Egypt",
				})

		})

	server.Run(UTIL_API_HOST + ":" + UTIL_API_PORT)

}
