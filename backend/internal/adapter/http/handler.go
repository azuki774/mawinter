package http

import (
	"log/slog"
	"net/http"

	"github.com/azuki774/mawinter/api"
	"github.com/gin-gonic/gin"
)

// 以下、api.ServerInterface の実装

// Get - health check (GET /v3/)
func (s *Server) Get(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GetV3Categories - get categories (GET /v3/categories)
func (s *Server) GetV3Categories(c *gin.Context) {
	categories, err := s.categoryService.GetAllCategories(c.Request.Context())
	if err != nil {
		slog.Error("Failed to get categories", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get categories"})
		return
	}

	// ドメインエンティティをAPIレスポンス型に変換
	response := make([]api.Category, len(categories))
	for i, cat := range categories {
		response[i] = api.Category{
			CategoryId:   cat.CategoryID,
			CategoryName: cat.Name,
			CategoryType: api.CategoryType(cat.CategoryType.String()),
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetV3Record - get records (GET /v3/record)
func (s *Server) GetV3Record(c *gin.Context, params api.GetV3RecordParams) {
	// params には num, offset, yyyymm, category_id が含まれる
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

// PostV3Record - create record (POST /v3/record)
func (s *Server) PostV3Record(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

// GetV3RecordAvailable - record available (GET /v3/record/available)
func (s *Server) GetV3RecordAvailable(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

// GetV3RecordCount - record count (GET /v3/record/count)
func (s *Server) GetV3RecordCount(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

// GetV3RecordYear - get year summary (GET /v3/record/summary/{year})
func (s *Server) GetV3RecordYear(c *gin.Context, year int) {
	// year パラメータは自動的にパースされて渡される
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

// DeleteV3RecordId - delete record from id (DELETE /v3/record/{id})
func (s *Server) DeleteV3RecordId(c *gin.Context, id int) {
	// id パラメータは自動的にパースされて渡される
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

// GetV3RecordId - get record from id (GET /v3/record/{id})
func (s *Server) GetV3RecordId(c *gin.Context, id int) {
	// id パラメータは自動的にパースされて渡される
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

// GetV3Version - get version (GET /v3/version)
func (s *Server) GetV3Version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":   s.version,
		"reversion": s.revision,
		"build":     s.build,
	})
}
