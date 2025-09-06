// internal/api/handlers.go
package api

import (
	"net/http"

	"github.com/mohan2020coder/mSpace/internal/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// --- AUTH HANDLERS ---

func loginHandler(cfg *config.Config, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		// TODO: check user in Postgres
		log.Info("User login attempt", zap.String("username", creds.Username))

		c.JSON(http.StatusOK, gin.H{"token": "mock-jwt-token"})
	}
}

func registerHandler(cfg *config.Config, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Email    string `json:"email"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		// TODO: insert user in Postgres
		log.Info("User registration", zap.String("username", input.Username))

		c.JSON(http.StatusCreated, gin.H{"status": "registered"})
	}
}

// --- ITEM HANDLERS ---

func uploadItemHandler(cfg *config.Config, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no file provided"})
			return
		}

		// TODO: save file to MinIO using cfg.Storage
		log.Info("File uploaded", zap.String("filename", file.Filename))

		c.JSON(http.StatusOK, gin.H{"status": "uploaded", "filename": file.Filename})
	}
}

func getItemHandler(cfg *config.Config, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		// TODO: query metadata from Postgres
		log.Info("Get item", zap.String("id", id))

		c.JSON(http.StatusOK, gin.H{
			"id":       id,
			"title":    "Sample Title",
			"author":   "John Doe",
			"abstract": "This is a mock abstract.",
		})
	}
}

func listItemsHandler(cfg *config.Config, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: fetch items list from Postgres
		log.Info("List items")

		items := []map[string]string{
			{"id": "1", "title": "First Item"},
			{"id": "2", "title": "Second Item"},
		}

		c.JSON(http.StatusOK, items)
	}
}

// --- SEARCH HANDLER ---

func searchHandler(cfg *config.Config, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("q")
		// TODO: search in Bleve/Elasticsearch
		log.Info("Search request", zap.String("query", query))

		c.JSON(http.StatusOK, gin.H{
			"query":   query,
			"results": []string{"mock-result-1", "mock-result-2"},
		})
	}
}
