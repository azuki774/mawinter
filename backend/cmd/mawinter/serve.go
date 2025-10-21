package main

import (
	"fmt"

	"github.com/azuki774/mawinter/internal/adapter/http"
	"github.com/spf13/cobra"
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
	fmt.Printf("Starting server on %s:%d\n", host, port)
	fmt.Printf("Version: %s, Revision: %s, Build: %s\n", version, revision, build)

	// HTTPサーバの起動
	server := http.NewServer(host, port, version, revision, build)
	return server.Start()
}
