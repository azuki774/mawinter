package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger は、log/slog を使った構造化ログ出力を行う Gin ミドルウェアです。
// JSON形式でHTTPリクエストとレスポンスの情報を記録します。
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// リクエスト開始時刻を記録
		start := time.Now()

		// リクエストパスとメソッドを取得
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 次のミドルウェア/ハンドラを実行
		c.Next()

		// レスポンス情報を取得
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		userAgent := c.Request.UserAgent()

		// クエリパラメータがある場合はパスに追加
		if raw != "" {
			path = path + "?" + raw
		}

		// ログレベルをステータスコードに応じて決定
		logLevel := slog.LevelInfo
		if statusCode >= 500 {
			logLevel = slog.LevelError
		} else if statusCode >= 400 {
			logLevel = slog.LevelWarn
		}

		// 構造化ログを出力
		slog.Log(c.Request.Context(), logLevel, "HTTP request",
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status", statusCode),
			slog.Float64("latency_ms", latency.Seconds()*1000),
			slog.String("client_ip", clientIP),
			slog.String("user_agent", userAgent),
		)

		// エラーがある場合は追加でログ出力
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				slog.ErrorContext(c.Request.Context(), "Request error",
					slog.String("error", e.Error()),
					slog.String("path", path),
				)
			}
		}
	}
}
