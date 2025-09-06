// internal/api/router.go
package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/mohan2020coder/mSpace/internal/config"
)

func SetupRouter(cfg *config.Config, log *zap.Logger) *gin.Engine {
	r := gin.New()

	// Attach middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Base health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"version": "v1",
		})
	})

	// --- Modular routes ---
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", loginHandler(cfg, log))
		authGroup.POST("/register", registerHandler(cfg, log))
	}

	itemsGroup := r.Group("/items")
	{
		itemsGroup.POST("/upload", uploadItemHandler(cfg, log))
		itemsGroup.GET("/:id", getItemHandler(cfg, log))
		itemsGroup.GET("/", listItemsHandler(cfg, log))
	}

	searchGroup := r.Group("/search")
	{
		searchGroup.GET("/", searchHandler(cfg, log))
	}

	return r
}
