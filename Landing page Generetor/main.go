package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itsvictorfy/pkg/chatgpt"
	"github.com/joho/godotenv"
)

var (
	chatGptClient *chatgpt.ChatGPTClient
	prompt        string                  // prompt for ChatGPT
	pageStorageMu = NewLandingPageStore() // Mutex for thread-safe access to pageStorage
	templateMap   map[string]string
)

func LoadTemplates(dir string) error {

	// Read all files in the templates directory
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read templates directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue // Skip directories
		}

		// Build full file path
		filePath := filepath.Join(dir, file.Name())

		// Read file content
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Printf("failed to read template file %s: %v", file.Name(), err)
			continue
		}

		// Use the filename (without extension) as the template ID
		templateID := file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]
		templateMap[templateID] = string(content)
	}

	return nil
}

func init() {
	godotenv.Load()
	chatGptClient, _ = chatgpt.InitClient(os.Getenv("OPENAI_API_KEY"))
	// promptfilepath := "./ptompt.txt"
	// bytePrompt, err := os.ReadFile(promptfilepath)
	// if err != nil {
	// 	slog.Error(fmt.Sprintf("Failed to load prompt file path: %s, errpr: %s", promptfilepath, err))
	// }
	prompt = "You Are a Landing Page Builder, You Recieve a Prompt and return a HTML CSS Code, Your Response sould be start with <html> and end with </html> Within these tags you do all the HTML and CSS needed. Use this HTML as a template-"
	LoadTemplates("./templates")
	// Debugging
	for id, content := range templateMap {
		fmt.Printf("Loaded template: %s\n", id)
		fmt.Printf("Content:\n%s\n", content)
	}
}

func getIDFromCookies(c *gin.Context) string {
	id, err := c.Cookie("victor-website-landingpage-sessionID")
	if err != nil {
		id = uuid.New().String()
		c.SetCookie("victor-website-landingpage-sessionID", id, 3600, "/", "localhost", false, false)
	} else {
		slog.Info(fmt.Sprintf("ID Found in Cookies, ID: %s", id))
	}
	return id
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.*")
	router.Static("/static", "./static")

	if os.Getenv("HOMEPAGE") != "false" {
		router.GET("/", func(c *gin.Context) {
			getIDFromCookies(c)
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
