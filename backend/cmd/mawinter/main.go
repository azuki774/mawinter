package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// ビルド時に埋め込まれるバージョン情報
var (
	version  = "dev"
	revision = "unknown"
	build    = "dev"
)

var rootCmd = &cobra.Command{
	Use:   "mawinter",
	Short: "Mawinter - 家計簿サーバ",
	Long:  "Mawinter は Go/Nuxt3 で構築された家計簿サーバです。",
}

func main() {
	fmt.Printf("Version: %s, Revision: %s, Build: %s\n", version, revision, build)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
