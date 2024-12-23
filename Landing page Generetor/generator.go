package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Function to generate a landing page
func generateLandingPage(c *gin.Context) {
	input := c.PostForm("query")

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
			History: InitChatGptHistory(prompt), // Assume `prompt` is globally available
		}
	}
	page.HTML, page.History, err = chatGptClient.SendChatGptRequestWithHistory(input, page.History)
	if err != nil {
		slog.Error(fmt.Sprintf("Error generating content from ChatGPT: %s", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate content from ChatGPT"})
		return
	}

	pageStorageMu.Add(page)

	c.JSON(http.StatusOK, gin.H{
		"message": "Page generated successfully!",
		"url":     generateUrl(c) + page.ID,
	})

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
			History: InitChatGptHistory(prompt), // Assume `prompt` is globally available
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
