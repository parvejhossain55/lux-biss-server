package upload

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/parvej/luxbiss_server/internal/common"
)

type Handler struct {
	uploadDir string
}

func NewHandler(uploadDir string) *Handler {
	// Create upload dir if not exists
	os.MkdirAll(uploadDir, os.ModePerm)
	return &Handler{
		uploadDir: uploadDir,
	}
}

func (h *Handler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		common.BadRequest(c, "No image file provided", nil)
		return
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		common.BadRequest(c, "Invalid file type. Only JPG, PNG, and WebP are allowed", nil)
		return
	}

	// Validate file size (e.g., max 2MB)
	if file.Size > 2*1024*1024 {
		common.BadRequest(c, "File size exceeds 2MB limit", nil)
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s_%d%s", uuid.New().String()[:8], time.Now().Unix(), ext)
	filepath := filepath.Join(h.uploadDir, filename)

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, common.Response{
			Success:   false,
			Message:   "Failed to save image",
			RequestID: c.GetString("request_id"),
		})
		return
	}

	// Return public URL (assuming it's served statically at /uploads)
	url := fmt.Sprintf("/uploads/%s", filename)

	common.OKWithMeta(c, "Image uploaded successfully", map[string]string{
		"url": url,
	}, nil)
}
