// internal/api/router.go
package api

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/mohan2020coder/mSpace/internal/config"
	"github.com/mohan2020coder/mSpace/internal/models"
	"github.com/mohan2020coder/mSpace/internal/search"
	"github.com/mohan2020coder/mSpace/internal/storage"
	"gorm.io/gorm"
)

type App struct {
	Cfg    *config.Config
	DB     *gorm.DB
	Minio  *storage.MinioClient
	Logger *zap.Logger
}

func SetupRouter(app *App, searchIndex *search.SearchIndex) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	// Custom logger can be added (zap -> ginadapter) but keep default for now
	r.Use(gin.Logger())

	// ----------------- CORS -----------------
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // allow all for now, can restrict domains
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/report", func(c *gin.Context) {
		start := time.Now()

		// Read body safely
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
			// Restore body for Gin to avoid consuming it
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Gather headers
		headers := map[string]string{}
		for k, v := range c.Request.Header {
			headers[k] = v[0]
		}

		// Gather cookies
		cookies := []map[string]string{}
		for _, cookie := range c.Request.Cookies() {
			cookies = append(cookies, map[string]string{
				"name":  cookie.Name,
				"value": cookie.Value,
			})
		}

		// Gather query parameters
		queryParams := map[string]string{}
		for k, v := range c.Request.URL.Query() {
			queryParams[k] = v[0]
		}

		// Build diagnostic report
		report := gin.H{
			"time":             time.Now().Format(time.RFC3339),
			"client_ip":        c.ClientIP(),
			"method":           c.Request.Method,
			"path":             c.Request.URL.Path,
			"protocol":         c.Request.Proto,
			"host":             c.Request.Host,
			"headers":          headers,
			"query_params":     queryParams,
			"cookies":          cookies,
			"body":             string(bodyBytes),
			"response_time_ms": time.Since(start).Milliseconds(),
		}

		c.JSON(http.StatusOK, report)
	})

	// health
	r.GET("/health", func(c *gin.Context) {
		start := time.Now()

		// Basic client info
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		userAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		proto := c.Request.Proto
		host := c.Request.Host
		contentType := c.GetHeader("Content-Type")
		accept := c.GetHeader("Accept")
		language := c.GetHeader("Accept-Language")
		encoding := c.GetHeader("Accept-Encoding")

		// Log a detailed client report
		log.Printf(`
		[HEALTH CHECK]
		Time: %s
		Client IP: %s
		Host: %s
		HTTP Method: %s
		Path: %s
		Protocol: %s
		User-Agent: %s
		Referer: %s
		Content-Type: %s
		Accept: %s
		Accept-Language: %s
		Accept-Encoding: %s
		`, time.Now().Format(time.RFC3339), clientIP, host, method, path, proto, userAgent, referer, contentType, accept, language, encoding)

		// Respond to client
		c.JSON(http.StatusOK, gin.H{"status": "ok"})

		// Log response time
		duration := time.Since(start)
		log.Printf("[INFO] Response sent in %v\n", duration)
	})

	// Communities
	r.GET("/api/communities", listCommunitiesHandler(app))
	r.POST("/api/communities", createCommunityHandler(app))

	// Collections
	r.GET("/api/collections", listCollectionsHandler(app))
	r.POST("/api/collections", createCollectionHandler(app))

	r.GET("/api/search", func(c *gin.Context) {
		q := c.Query("q")
		if q == "" {
			c.JSON(400, gin.H{"error": "query required"})
			return
		}

		collectionID := uint(0)
		if cid := c.Query("collection_id"); cid != "" {
			parsed, _ := strconv.Atoi(cid)
			collectionID = uint(parsed)
		}

		author := c.Query("author")

		ids, err := searchIndex.Search(q, collectionID, author)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		var items []models.Item
		if len(ids) > 0 {
			app.DB.Where("id IN ?", ids).Find(&items)
		}

		c.JSON(200, items)
	})

	// Items
	items := r.Group("/api/items")
	items.POST("", createItemHandler(app))
	items.GET("", listItemsHandler(app))
	items.GET("/:id", getItemHandler(app))
	items.POST("/:id/file", uploadFileHandler(app, searchIndex))
	items.POST("/:id/publish", publishItemHandler(app))
	items.POST("/:id/reject", rejectItemHandler(app))

	return r
}
