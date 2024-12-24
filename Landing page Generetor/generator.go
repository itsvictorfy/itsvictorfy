package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itsvictorfy/pkg/chatgpt"
)

// Function to generate a landing page
func generateLandingPage(c *gin.Context) {
	var requestData struct {
		Query string `json:"query"`
	}

	// Parse JSON body
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	input := requestData.Query
	slog.Info(fmt.Sprintf("Received request to generate landing page input- %s", input))

	id, err := c.Cookie("victor-website-landingpage-sessionID")
	if err != nil {
		id = uuid.New().String()
		c.SetCookie("victor-website-landingpage-sessionID", id, 3600, "/", "localhost", false, false)
		slog.Error(fmt.Sprintf("Cant Find victor-website-landingpage-sessionID Cookie - Error: %s, Created new one %s", err, id))
	}

	page, exists := pageStorageMu.Get(id)
	if !exists {
		page = LandingPage{
			ID:      id,
			HTML:    "",
			History: chatgpt.InitChatGptHistory(prompt), // Assume `prompt` is globally available
			Status:  "Initializing",                     // Set initial status
		}
		pageStorageMu.Add(page)
	}

	page.Status = "Waiting for ChatGPT"
	pageStorageMu.Add(page)

	page.HTML, page.History, err = chatGptClient.SendChatGptRequestWithHistory(input, page.History)
	if err != nil {
		page.Status = "Error: Failed to generate content from ChatGPT"
		pageStorageMu.Add(page)
		slog.Error(fmt.Sprintf("Error generating content from ChatGPT: %s", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate content from ChatGPT"})
		return
	}

	page.Status = "Deploying Page"
	pageStorageMu.Add(page)

	c.JSON(http.StatusOK, gin.H{
		"message": "Page generated successfully!",
		"url":     generateUrl(c) + page.ID,
	})

	page.Status = "Completed"
	pageStorageMu.Add(page)
}

func testGeneratedPageWithTestPage(c *gin.Context) {
	id, err := c.Cookie("victor-website-landingpage-sessionID")
	if err != nil {
		id = uuid.New().String()
		c.SetCookie("victor-website-landingpage-sessionID", id, 3600, "/", "localhost", false, false)
		slog.Error(fmt.Sprintf("Cant Find victor-website-landingpage-sessionID Cookie - Error: %s, Created new one %s", err, id))
	}

	page, exists := pageStorageMu.Get(id)
	if !exists {
		page = LandingPage{
			ID:      id,
			HTML:    "",
			History: chatgpt.InitChatGptHistory(prompt), // Assume `prompt` is globally available
		}
	}
	htmlFilePath := "./templates/html-generetor-test.html"

	htmlContent, err := os.ReadFile(htmlFilePath)
	if err != nil {
		slog.Error("Failed to load HTML file", "path", htmlFilePath, "error", err)
		c.String(http.StatusInternalServerError, "Failed to load HTML file")
		return
	}

	page.HTML = string(htmlContent)

	pageStorageMu.Add(page)
	c.JSON(http.StatusOK, gin.H{
		"message": "Page generated successfully!",
		"url":     generateUrl(c) + page.ID,
	})
	slog.Info("Landing page generated successfully")
}

func generateUrl(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	fullURL := fmt.Sprintf("%s://%s/pages/", scheme, c.Request.Host)
	return fullURL
}

func progressEndpoint(c *gin.Context) {
	id := c.Param("id")

	page, exists := pageStorageMu.Get(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": page.Status})
}
