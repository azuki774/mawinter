package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Server は HTTP サーバの構造体
type Server struct {
	router *gin.Engine
	host   string
	port   int
}

// NewServer は新しい HTTP サーバを作成
func NewServer(host string, port int) *Server {
	router := gin.Default()

	s := &Server{
		router: router,
		host:   host,
		port:   port,
	}

	s.setupRoutes()
	return s
}

// setupRoutes はルーティングを設定
func (s *Server) setupRoutes() {
	// ヘルスチェックエンドポイント
	s.router.GET("/v3/", s.healthCheck)

	// 他のエンドポイントは空の実装
	s.router.GET("/v3/categories", s.getCategories)
	s.router.GET("/v3/record", s.getRecord)
	s.router.POST("/v3/record", s.postRecord)
	s.router.GET("/v3/record/available", s.getRecordAvailable)
	s.router.GET("/v3/record/count", s.getRecordCount)
	s.router.GET("/v3/record/summary/:year", s.getRecordYear)
	s.router.DELETE("/v3/record/:id", s.deleteRecordId)
	s.router.GET("/v3/record/:id", s.getRecordId)
	s.router.GET("/v3/version", s.getVersion)
}

// Start はサーバを起動
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	return s.router.Run(addr)
}

// ハンドラー関数（すべて空の実装）

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) getCategories(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

func (s *Server) getRecord(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

func (s *Server) postRecord(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

func (s *Server) getRecordAvailable(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

func (s *Server) getRecordCount(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

func (s *Server) getRecordYear(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

func (s *Server) deleteRecordId(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

func (s *Server) getRecordId(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

func (s *Server) getVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": "v3.0.0"})
}
