// internal/api/handlers.go
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/ledongthuc/pdf"

	"github.com/mohan2020coder/mSpace/internal/models"
	"github.com/mohan2020coder/mSpace/internal/search"
)

// type createItemReq struct {
// 	Title    string `json:"title" binding:"required"`
// 	Author   string `json:"author"`
// 	Abstract string `json:"abstract"`
// }

// func createItemHandler(app *App) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var req createItemReq
// 		if err := c.ShouldBindJSON(&req); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		item := models.Item{
// 			Title:    req.Title,
// 			Author:   req.Author,
// 			Abstract: req.Abstract,
// 			Status:   "SUBMITTED",
// 		}
// 		if err := app.DB.Create(&item).Error; err != nil {
// 			app.Logger.Error("db create failed", zap.Error(err))
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create item"})
// 			return
// 		}
// 		c.JSON(http.StatusCreated, item)
// 	}
// }
// ---------------- Items ----------------

type createItemReq struct {
	Title        string `json:"title" binding:"required"`
	Author       string `json:"author"`
	Abstract     string `json:"abstract"`
	CollectionID uint   `json:"collection_id" binding:"required"`
	Visibility   string `json:"visibility"` // PUBLIC/PRIVATE
	LegalJSON    string `json:"legal_json"` // optional
}

func listItemsHandler(app *App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var items []models.Item
		if err := app.DB.Order("created_at desc").Find(&items).Error; err != nil {
			app.Logger.Error("db list failed", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list items"})
			return
		}
		c.JSON(http.StatusOK, items)
	}
}

func getItemHandler(app *App) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		var item models.Item
		if err := app.DB.First(&item, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
			return
		}
		c.JSON(http.StatusOK, item)
	}
}

// func uploadFileHandler(app *App) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		idStr := c.Param("id")
// 		id, err := strconv.Atoi(idStr)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
// 			return
// 		}
// 		var item models.Item
// 		if err := app.DB.First(&item, id).Error; err != nil {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
// 			return
// 		}

// 		// file in form field "file"
// 		fh, err := c.FormFile("file")
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
// 			return
// 		}

// 		src, err := fh.Open()
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
// 			return
// 		}
// 		defer src.Close()

// 		ext := filepath.Ext(fh.Filename)
// 		objectName := fmt.Sprintf("item-%d-%d%s", item.ID, time.Now().Unix(), ext)

// 		// Upload to minio
// 		ctx := context.Background()
// 		if _, err := app.Minio.UploadStream(ctx, objectName, src, fh.Size, fh.Header.Get("Content-Type")); err != nil {
// 			app.Logger.Error("minio upload failed", zap.Error(err))
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file"})
// 			return
// 		}

// 		// Generate presigned URL (valid 24 hours)
// 		url, err := app.Minio.PresignedURL(ctx, objectName, 24*time.Hour)
// 		if err != nil {
// 			app.Logger.Error("presign failed", zap.Error(err))
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get file url"})
// 			return
// 		}

// 		// update DB
// 		item.FileURL = url
// 		item.Status = "SUBMITTED" // keep submitted until review; or set PUBLISHED if auto
// 		if err := app.DB.Save(&item).Error; err != nil {
// 			app.Logger.Error("db update failed", zap.Error(err))
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update item"})
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{"message": "file uploaded", "file_url": url})
// 	}
// }

// package api

// import (
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"path/filepath"
// 	"strconv"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"go.uber.org/zap"

// 	"github.com/mohan2020coder/mSpace/internal/models"
// )

// ---------------- Communities ----------------

func createCommunityHandler(app *App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var community models.Community
		if err := c.ShouldBindJSON(&community); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := app.DB.Create(&community).Error; err != nil {
			app.Logger.Error("db create community failed", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create community"})
			return
		}
		c.JSON(http.StatusCreated, community)
	}
}

func listCommunitiesHandler(app *App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var communities []models.Community
		if err := app.DB.Preload("Collections").Find(&communities).Error; err != nil {
			app.Logger.Error("db list communities failed", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list communities"})
			return
		}
		c.JSON(http.StatusOK, communities)
	}
}

// ---------------- Collections ----------------

func createCollectionHandler(app *App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var collection models.Collection
		if err := c.ShouldBindJSON(&collection); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := app.DB.Create(&collection).Error; err != nil {
			app.Logger.Error("db create collection failed", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create collection"})
			return
		}
		c.JSON(http.StatusCreated, collection)
	}
}

func listCollectionsHandler(app *App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var collections []models.Collection
		if err := app.DB.Preload("Items").Find(&collections).Error; err != nil {
			app.Logger.Error("db list collections failed", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list collections"})
			return
		}
		c.JSON(http.StatusOK, collections)
	}
}

func createItemHandler(app *App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createItemReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Initialize item
		item := models.Item{
			Title:        req.Title,
			Author:       req.Author,
			Abstract:     req.Abstract,
			CollectionID: req.CollectionID,
			Status:       "DRAFT",
			Version:      0,
			Visibility:   req.Visibility,
			FullText:     "",   // fine as empty string
			LegalJSON:    "{}", // must be valid JSON
		}

		if err := app.DB.Create(&item).Error; err != nil {
			app.Logger.Error("db create item failed", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create item"})
			return
		}

		c.JSON(http.StatusCreated, item)
	}
}

// Upload file with versioning

// uploadFileHandler handles file upload + extraction + indexing
// func uploadFileHandler(app *App, searchIndex *search.SearchIndex) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// --- Parse item ID ---
// 		idStr := c.Param("id")
// 		id, err := strconv.Atoi(idStr)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
// 			return
// 		}

// 		// --- Load item from DB ---
// 		var item models.Item
// 		if err := app.DB.First(&item, id).Error; err != nil {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
// 			return
// 		}

// 		// --- Get uploaded file ---
// 		fh, err := c.FormFile("file")
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
// 			return
// 		}

// 		src, err := fh.Open()
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
// 			return
// 		}
// 		defer src.Close()

// 		// --- Increment version + object name ---
// 		item.Version++
// 		ext := strings.ToLower(filepath.Ext(fh.Filename))
// 		objectName := fmt.Sprintf("item-%d-v%d-%d%s", item.ID, item.Version, time.Now().Unix(), ext)

// 		// --- Upload to MinIO ---
// 		ctx := context.Background()
// 		if _, err := app.Minio.UploadStream(ctx, objectName, src, fh.Size, fh.Header.Get("Content-Type")); err != nil {
// 			app.Logger.Error("minio upload failed", zap.Error(err))
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file"})
// 			return
// 		}

// 		url, err := app.Minio.PresignedURL(ctx, objectName, 24*time.Hour)
// 		if err != nil {
// 			app.Logger.Error("presign failed", zap.Error(err))
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get file url"})
// 			return
// 		}

// 		item.FileURL = url
// 		item.Status = "SUBMITTED"

// 		// --- Handle PDF text extraction ---
// 		if ext == ".pdf" {
// 			tempDir := "./temp"
// 			os.MkdirAll(tempDir, 0755)
// 			tempPath := filepath.Join(tempDir, fmt.Sprintf("%d_%s", time.Now().UnixNano(), fh.Filename))

// 			if err := c.SaveUploadedFile(fh, tempPath); err == nil {
// 				app.Logger.Info("Saved PDF to temp", zap.String("path", tempPath), zap.Int64("size", fh.Size))

// 				extracted := ExtractPDFTextWithTikaFallback(tempPath, app.Logger)
// 				app.Logger.Info("Extracted PDF text", zap.Int("length", len(extracted)))

// 				item.FullText = extracted
// 				os.Remove(tempPath)
// 			} else {
// 				app.Logger.Error("Failed to save uploaded file", zap.Error(err))
// 			}
// 		}

// 		// --- Save item to DB ---
// 		if err := app.DB.Save(&item).Error; err != nil {
// 			app.Logger.Error("db update failed", zap.Error(err))
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update item"})
// 			return
// 		}

// 		// --- Index item in Bleve ---
// 		if err := searchIndex.IndexItem(&item); err != nil {
// 			app.Logger.Error("bleve index failed", zap.Error(err))
// 		}

//			c.JSON(http.StatusOK, gin.H{
//				"message":  "file uploaded",
//				"file_url": url,
//			})
//		}
//	}
func uploadFileHandler(app *App, searchIndex *search.SearchIndex) gin.HandlerFunc {
	return func(c *gin.Context) {
		// --- Parse item ID ---
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		// --- Load item from DB ---
		var item models.Item
		if err := app.DB.First(&item, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
			return
		}

		// --- Get uploaded file ---
		fh, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
			return
		}

		src, err := fh.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
			return
		}
		defer src.Close()

		// --- Increment version + object name ---
		item.Version++
		ext := strings.ToLower(filepath.Ext(fh.Filename))
		objectName := fmt.Sprintf("item-%d-v%d-%d%s", item.ID, item.Version, time.Now().Unix(), ext)

		// --- Upload to MinIO ---
		ctx := context.Background()
		if _, err := app.Minio.UploadStream(ctx, objectName, src, fh.Size, fh.Header.Get("Content-Type")); err != nil {
			app.Logger.Error("minio upload failed", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file"})
			return
		}

		url, err := app.Minio.PresignedURL(ctx, objectName, 24*time.Hour)
		if err != nil {
			app.Logger.Error("presign failed", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get file url"})
			return
		}

		item.FileURL = url
		item.Status = "SUBMITTED"

		// --- Handle PDF extraction ---
		if ext == ".pdf" {
			tempDir := "./temp"
			os.MkdirAll(tempDir, 0755)
			tempPath := filepath.Join(tempDir, fmt.Sprintf("%d_%s", time.Now().UnixNano(), fh.Filename))

			if err := c.SaveUploadedFile(fh, tempPath); err == nil {
				app.Logger.Info("Saved PDF to temp", zap.String("path", tempPath), zap.Int64("size", fh.Size))

				fullText := ExtractPDFTextWithTikaFallback(tempPath, app.Logger)
				app.Logger.Info("Extracted PDF text", zap.Int("length", len(fullText)))

				item.FullText = fullText

				// --- Parse structured legal document ---
				// legalDoc := search.ParseLegalDocument(fullText)
				// item.LegalJSON, _ = json.Marshal(legalDoc)

				legalDoc := search.ParseLegalDocument(fullText)
				b, _ := json.Marshal(legalDoc)
				item.LegalJSON = string(b) // store as plain JSON string in DB

				os.Remove(tempPath)
			} else {
				app.Logger.Error("Failed to save uploaded file", zap.Error(err))
			}
		}

		// --- Save item to DB ---
		if err := app.DB.Save(&item).Error; err != nil {
			app.Logger.Error("db update failed", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update item"})
			return
		}

		// --- Index item in Bleve ---
		if err := searchIndex.IndexItem(&item); err != nil {
			app.Logger.Error("bleve index failed", zap.Error(err))
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "file uploaded",
			"file_url": url,
		})
	}
}

// --- PDF Extraction helpers with logs ---

func ExtractPDFTextWithTikaFallback(filePath string, logger *zap.Logger) string {
	text := extractWithLedongthuc(filePath, logger)
	if len(strings.TrimSpace(text)) < 10 {
		logger.Info("Digital text extraction empty, falling back to Tika OCR")
		return extractWithTika(filePath, logger)
	}
	return text
}

func extractWithLedongthuc(filePath string, logger *zap.Logger) string {
	logger.Info("Starting digital extraction", zap.String("file", filePath))
	f, r, err := pdf.Open(filePath)
	if err != nil {
		logger.Error("Failed to open PDF", zap.Error(err))
		return ""
	}
	defer f.Close()

	var content strings.Builder
	for i := 1; i <= r.NumPage(); i++ {
		page := r.Page(i)
		if page.V.IsNull() {
			continue
		}
		text, err := page.GetPlainText(nil)
		if err != nil {
			logger.Warn("Failed page extraction", zap.Int("page", i), zap.Error(err))
			continue
		}
		content.WriteString(text + "\n")
	}

	lines := strings.Split(content.String(), "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		line = strings.Join(strings.Fields(line), " ")
		lines[i] = line
	}

	extracted := strings.Join(lines, "\n")
	logger.Info("Digital extraction completed", zap.Int("length", len(extracted)))
	return extracted
}

func extractWithTika(filePath string, logger *zap.Logger) string {
	logger.Info("Starting Tika OCR extraction", zap.String("file", filePath))

	file, err := os.Open(filePath)
	if err != nil {
		logger.Error("Failed to open PDF for Tika", zap.Error(err))
		return ""
	}
	defer file.Close()

	// Tika endpoint expects PUT or POST with Accept: text/plain
	req, _ := http.NewRequest("PUT", "http://localhost:9998/tika", file)
	req.Header.Set("Accept", "text/plain")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Tika request failed", zap.Error(err))
		return ""
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	text := strings.TrimSpace(string(body))
	logger.Info("Tika OCR completed", zap.Int("length", len(text)))

	if len(text) == 0 {
		logger.Warn("Tika OCR returned empty text, check Tesseract installation")
	}
	return text
}

// ---------------- Workflow ----------------

func publishItemHandler(app *App) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, _ := strconv.Atoi(idStr)
		var item models.Item
		if err := app.DB.First(&item, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
			return
		}
		item.Status = "PUBLISHED"
		if err := app.DB.Save(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish item"})
			return
		}
		c.JSON(http.StatusOK, item)
	}
}

func rejectItemHandler(app *App) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, _ := strconv.Atoi(idStr)
		var item models.Item
		if err := app.DB.First(&item, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
			return
		}
		item.Status = "REJECTED"
		if err := app.DB.Save(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reject item"})
			return
		}
		c.JSON(http.StatusOK, item)
	}
}
