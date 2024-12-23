package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/rand"
)

func progressEndpoint(c *gin.Context) {
	// id, err := c.Cookie("victor-website-landingpage-sessionID")
	// if err != nil {
	// 	slog.Error(fmt.Sprintf("Cannot get cookie error: %s", err))
	// 	progress := fmt.Sprintf("Processing... %s%% complete", rand.Intn(100))
	// 	c.String(http.StatusOK, progress)
	// }

	// Mock progress updates for simplicity
	// In a real implementation, fetch actual progress from a store
	progress := fmt.Sprintf("Processing... %d%% complete", rand.Intn(100))

	c.String(http.StatusOK, progress)
}
