package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

var (
	chatGptClient *ChatGPTClient
	prompt        string                  // prompt for ChatGPT
	pageStorageMu = NewLandingPageStore() // Mutex for thread-safe access to pageStorage
)

func init() {
	godotenv.Load()
	chatGptClient = InitChatGptClient(os.Getenv("OPENAI_API_KEY"))
	// promptfilepath := "./ptompt.txt"
	// bytePrompt, err := os.ReadFile(promptfilepath)
	// if err != nil {
	// 	slog.Error(fmt.Sprintf("Failed to load prompt file path: %s, errpr: %s", promptfilepath, err))
	// }
	prompt = "You Are a Landing Page Builder, You Recieve a Prompt and return a HTML CSS Code, Your Response sould be start with <html> and end with </html> Within these tags you do all the HTML and CSS needed"
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.*")
	router.Static("/static", "./static")

	if os.Getenv("HOMEPAGE") != "false" {
		router.GET("/", func(c *gin.Context) {
			id := uuid.New().String()
			c.SetCookie("victor-website-landingpage-sessionID", id, 3600, "/", "localhost", false, false)
			slog.Info("victor-website-landingpage-sessionID: %s", slog.String("ID", id))
			c.HTML(http.StatusOK, "html-generator.html", nil)
		})
	}
	router.GET("/pages/:id", func(c *gin.Context) {
		id := c.Param("id")

		page, exists := pageStorageMu.Get(id)
		if !exists {
			c.String(http.StatusNotFound, "Page not found")
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(page.HTML))
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.GET("/progress/:id", progressEndpoint)

	router.POST("/test", testGeneratedPageWithTestPage)
	router.POST("/generate", generateLandingPage)

	router.Run() // Listen and serve on 0.0.0.0:8080
}
