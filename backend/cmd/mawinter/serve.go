package main

import (
	"fmt"
	"log/slog"

	"github.com/azuki774/mawinter/internal/adapter/http"
	"github.com/azuki774/mawinter/internal/adapter/repository"
	"github.com/azuki774/mawinter/internal/application"
	"github.com/azuki774/mawinter/pkg/config"
	"github.com/azuki774/mawinter/pkg/logger"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	port int
	host string
)

func init() {
	// serve コマンドを root コマンドに追加
	rootCmd.AddCommand(serveCmd)

	// フラグの定義
	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "HTTPサーバのポート番号")
	serveCmd.Flags().StringVarP(&host, "host", "H", "0.0.0.0", "HTTPサーバのホスト")
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "HTTPサーバを起動",
	Long:  "Mawinter の HTTP API サーバを起動します。",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runServer(host, port)
	},
}

func runServer(host string, port int) error {
	// デフォルトロガーの初期化
	slog.SetDefault(logger.New())

	slog.Info("Starting Mawinter server",
		slog.String("host", host),
		slog.Int("port", port),
		slog.String("version", version),
		slog.String("revision", revision),
		slog.String("build", build),
	)

	// データベース設定の読み込み
	dbInfo, err := config.LoadDBInfo()
	if err != nil {
		slog.Error("Failed to load database configuration",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to load database configuration: %w", err)
	}

	slog.Info("Database configuration loaded",
		slog.String("host", dbInfo.Host),
		slog.String("port", dbInfo.Port),
		slog.String("user", dbInfo.User),
		slog.String("name", dbInfo.Name),
	)

	// データベース接続の初期化
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbInfo.User,
		dbInfo.Pass,
		dbInfo.Host,
		dbInfo.Port,
		dbInfo.Name,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Failed to connect to database",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	slog.Info("Database connection established")

	// 依存性の注入
	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := application.NewCategoryService(categoryRepo)

	// HTTPサーバの起動
	server := http.NewServer(host, port, version, revision, build, dbInfo, categoryService)
	return server.Start()
}
