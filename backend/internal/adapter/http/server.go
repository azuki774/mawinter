package http

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/azuki774/mawinter/api"
	"github.com/azuki774/mawinter/pkg/config"
	"github.com/gin-gonic/gin"
)

// Server は HTTP サーバの構造体
// api.ServerInterface を実装する
type Server struct {
	router   *gin.Engine
	host     string
	port     int
	dbInfo   *config.DBInfo
	version  string
	revision string
	build    string
}

// NewServer は新しい HTTP サーバを作成
func NewServer(host string, port int, version, revision, build string, dbInfo *config.DBInfo) *Server {
	router := gin.Default()

	// プロキシを使わない設定
	router.SetTrustedProxies(nil)

	s := &Server{
		router:   router,
		host:     host,
		port:     port,
		dbInfo:   dbInfo,
		version:  version,
		revision: revision,
		build:    build,
	}

	// OpenAPI生成のRegisterHandlersを使用してルーティングを設定
	api.RegisterHandlers(s.router, s)

	return s
}

// Start はサーバを起動
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	slog.Info("HTTP server starting", slog.String("address", addr))

	if err := s.router.Run(addr); err != nil {
		slog.Error("Failed to start HTTP server",
			slog.String("address", addr),
			slog.String("error", err.Error()),
		)
		return err
	}
	return nil
}

// 以下、api.ServerInterface の実装

// Get - health check (GET /v3/)
func (s *Server) Get(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GetV3Categories - get categories (GET /v3/categories)
func (s *Server) GetV3Categories(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
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
