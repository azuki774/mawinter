package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/azuki774/mawinter/api"
	"github.com/azuki774/mawinter/internal/domain"
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
	// デフォルト値の設定
	num := 20
	if params.Num != nil {
		num = *params.Num
	}

	offset := 0
	if params.Offset != nil {
		offset = *params.Offset
	}

	yyyymm := ""
	if params.Yyyymm != nil {
		yyyymm = *params.Yyyymm
	}

	categoryID := 0
	if params.CategoryId != nil {
		categoryID = *params.CategoryId
	}

	// レコードを取得
	records, err := s.recordService.GetRecords(c.Request.Context(), num, offset, yyyymm, categoryID)
	if err != nil {
		slog.Error("Failed to get records", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get records"})
		return
	}

	// APIレスポンス型に変換
	response := make([]api.Record, len(records))
	for i, rec := range records {
		response[i] = api.Record{
			Id:           rec.ID,
			CategoryId:   rec.CategoryID,
			CategoryName: rec.CategoryName,
			Datetime:     rec.Datetime,
			From:         rec.From,
			Type:         rec.Type,
			Price:        rec.Price,
			Memo:         rec.Memo,
		}
	}

	c.JSON(http.StatusOK, response)
}

// PostV3Record - create record (POST /v3/record)
func (s *Server) PostV3Record(c *gin.Context) {
	var req api.ReqRecord
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Failed to bind request body", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// デフォルト値の設定
	datetime := "20060102" // デフォルトの日付フォーマット
	if req.Datetime != nil {
		datetime = *req.Datetime
	}

	from := ""
	if req.From != nil {
		from = *req.From
	}

	recordType := ""
	if req.Type != nil {
		recordType = *req.Type
	}

	memo := ""
	if req.Memo != nil {
		memo = *req.Memo
	}

	// datetimeをtime.Timeに変換（YYYYMMDD形式）
	parsedTime, err := parseDateTime(datetime)
	if err != nil {
		slog.Error("Failed to parse datetime", slog.String("datetime", datetime), slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid datetime format"})
		return
	}

	// ドメインエンティティを作成
	record := &domain.Record{
		CategoryID: req.CategoryId,
		Datetime:   parsedTime,
		From:       from,
		Type:       recordType,
		Price:      req.Price,
		Memo:       memo,
	}

	// レコードを作成
	createdRecord, err := s.recordService.CreateRecord(c.Request.Context(), record)
	if err != nil {
		slog.Error("Failed to create record", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create record"})
		return
	}

	// APIレスポンス型に変換
	response := api.Record{
		Id:           createdRecord.ID,
		CategoryId:   createdRecord.CategoryID,
		CategoryName: createdRecord.CategoryName,
		Datetime:     createdRecord.Datetime,
		From:         createdRecord.From,
		Type:         createdRecord.Type,
		Price:        createdRecord.Price,
		Memo:         createdRecord.Memo,
	}

	c.JSON(http.StatusCreated, response)
}

// GetV3RecordAvailable - record available (GET /v3/record/available)
func (s *Server) GetV3RecordAvailable(c *gin.Context) {
	// レコードの利用可能期間を取得
	yyyymm, fy, err := s.recordService.GetAvailablePeriods(c.Request.Context())
	if err != nil {
		slog.Error("Failed to get available periods", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get available periods"})
		return
	}

	// APIレスポンス型に変換
	response := gin.H{
		"fy":     fy,
		"yyyymm": yyyymm,
	}

	c.JSON(http.StatusOK, response)
}

// GetV3RecordCount - record count (GET /v3/record/count)
func (s *Server) GetV3RecordCount(c *gin.Context) {
	// クエリパラメータから条件を取得（オプション）
	yyyymm := c.Query("yyyymm")
	categoryID := 0
	if catIDStr := c.Query("category_id"); catIDStr != "" {
		// category_idが指定されている場合は変換
		var catID int
		if _, err := fmt.Sscanf(catIDStr, "%d", &catID); err == nil {
			categoryID = catID
		}
	}

	// レコード数を取得
	count, err := s.recordService.CountRecords(c.Request.Context(), yyyymm, categoryID)
	if err != nil {
		slog.Error("Failed to count records", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count records"})
		return
	}

	// APIレスポンス型に変換
	response := api.RecordCount{
		Num: &count,
	}

	c.JSON(http.StatusOK, response)
}

// GetV3RecordYear - get year summary (GET /v3/record/summary/{year})
func (s *Server) GetV3RecordYear(c *gin.Context, year int) {
	// year パラメータは自動的にパースされて渡される
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

// DeleteV3RecordId - delete record from id (DELETE /v3/record/{id})
func (s *Server) DeleteV3RecordId(c *gin.Context, id int) {
	// id パラメータは自動的にパースされて渡される
	err := s.recordService.DeleteRecord(c.Request.Context(), id)
	if err != nil {
		slog.Error("Failed to delete record", slog.Int("id", id), slog.String("error", err.Error()))
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}

	// 削除成功時は204 No Contentを返す
	c.Status(http.StatusNoContent)
}

// GetV3RecordId - get record from id (GET /v3/record/{id})
func (s *Server) GetV3RecordId(c *gin.Context, id int) {
	// id パラメータは自動的にパースされて渡される
	record, err := s.recordService.GetRecordByID(c.Request.Context(), id)
	if err != nil {
		slog.Error("Failed to get record", slog.Int("id", id), slog.String("error", err.Error()))
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}

	// APIレスポンス型に変換
	response := api.Record{
		Id:           record.ID,
		CategoryId:   record.CategoryID,
		CategoryName: record.CategoryName,
		Datetime:     record.Datetime,
		From:         record.From,
		Type:         record.Type,
		Price:        record.Price,
		Memo:         record.Memo,
	}

	c.JSON(http.StatusOK, response)
}

// GetV3Version - get version (GET /v3/version)
func (s *Server) GetV3Version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":   s.version,
		"reversion": s.revision,
		"build":     s.build,
	})
}

// parseDateTime はYYYYMMDD形式の文字列をtime.Timeに変換する
func parseDateTime(datetime string) (time.Time, error) {
	// YYYYMMDD形式をパース
	if len(datetime) == 8 {
		return time.Parse("20060102", datetime)
	}
	// それ以外の場合はRFC3339としてパース
	return time.Parse(time.RFC3339, datetime)
}
