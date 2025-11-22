package http

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/azuki774/mawinter/api"
	"github.com/azuki774/mawinter/internal/adapter/http/middleware"
	"github.com/azuki774/mawinter/internal/application"
	"github.com/azuki774/mawinter/pkg/config"
	"github.com/azuki774/mawinter/pkg/telemetry"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// Server は HTTP サーバの構造体
// api.ServerInterface を実装する
type Server struct {
	router          *gin.Engine
	host            string
	port            int
	dbInfo          *config.DBInfo
	version         string
	revision        string
	build           string
	categoryService *application.CategoryService
	recordService   *application.RecordService
}

// NewServer は新しい HTTP サーバを作成
func NewServer(host string, port int, version, revision, build string, dbInfo *config.DBInfo, categoryService *application.CategoryService, recordService *application.RecordService) *Server {
	router := gin.New()

	// ミドルウェアを設定
	router.Use(gin.Recovery())                                // panicからの回復
	router.Use(middleware.Logger())                           // 構造化ログ（JSON形式）
	// ヘルスチェック(/api/v3)は高頻度アクセスでノイズになるためトレースを除外する
	router.Use(otelgin.Middleware(
		telemetry.ServiceNameAPI,
		otelgin.WithFilter(func(r *http.Request) bool {
			path := r.URL.Path
			return path != "/api/v3" && path != "/api/v3/"
		}),
	)) // OpenTelemetryトレーシング

	// プロキシを使わない設定
	router.SetTrustedProxies(nil)

	s := &Server{
		router:          router,
		host:            host,
		port:            port,
		dbInfo:          dbInfo,
		version:         version,
		revision:        revision,
		build:           build,
		categoryService: categoryService,
		recordService:   recordService,
	}

	// OpenAPI生成のRegisterHandlersを使用してルーティングを設定
	// /api プレフィックスを追加
	api.RegisterHandlersWithOptions(s.router, s, api.GinServerOptions{
		BaseURL: "/api",
	})

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
