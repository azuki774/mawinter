package logger

import (
	"log/slog"
	"os"
	"runtime"
)

// New はJSON形式のslogロガーを作成する
// エラー時にスタックトレースと呼び出し元を表示する設定
func New() *slog.Logger {
	opts := &slog.HandlerOptions{
		AddSource: true, // 呼び出し元のファイル名と行番号を表示
		Level:     slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// エラーレベル以上の場合、スタックトレース情報を追加
			if a.Key == slog.LevelKey {
				// 型アサーションを安全に実行
				if level, ok := a.Value.Any().(slog.Level); ok {
					if level >= slog.LevelError {
						// スタックトレースを取得
						stackTrace := getStackTrace()
						return slog.Group("error_info",
							slog.String("level", level.String()),
							slog.String("stack_trace", stackTrace),
						)
					}
				}
			}
			return a
		},
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler)
}

// getStackTrace はスタックトレースを文字列として取得する
func getStackTrace() string {
	const maxDepth = 32
	var pcs [maxDepth]uintptr
	n := runtime.Callers(3, pcs[:])

	trace := ""
	frames := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := frames.Next()
		trace += frame.Function + "\n\t" + frame.File + ":" + string(rune(frame.Line)) + "\n"
		if !more {
			break
		}
	}
	return trace
}
