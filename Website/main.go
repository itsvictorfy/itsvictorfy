package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatsResponse struct {
	AvgLatencyMs  int `json:"avgLatencyMs"`
	ProjectsCount int `json:"projectsCount"`
	TimeSaved     int `json:"timeSaved"`
}

func main() {
	router := gin.Default()

	// Create an "/api" group
	api := router.Group("/api")
	{
		api.GET("/stats", func(c *gin.Context) {
			// In a real app, you'd query a DB or metrics aggregator here
			response := loadStats()
			c.JSON(http.StatusOK, response)
		})
	}

	router.Run(":8080") // Listen on port 8080
}
