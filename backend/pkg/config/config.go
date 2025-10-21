package config

import (
	"fmt"
	"os"
)

// DBInfo はデータベース接続情報を保持する構造体
type DBInfo struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

// LoadDBInfo は環境変数からデータベース接続情報を読み込む
func LoadDBInfo() (*DBInfo, error) {
	dbInfo := &DBInfo{
		Host: getEnv("DB_HOST", "localhost"),
		Port: getEnv("DB_PORT", "3306"),
		User: getEnv("DB_USER", "root"),
		Pass: getEnv("DB_PASS", ""),
		Name: getEnv("DB_NAME", "mawinter"),
	}

	// 必須項目のバリデーション
	if dbInfo.Host == "" {
		return nil, fmt.Errorf("DB_HOST is required")
	}
	if dbInfo.Port == "" {
		return nil, fmt.Errorf("DB_PORT is required")
	}
	if dbInfo.User == "" {
		return nil, fmt.Errorf("DB_USER is required")
	}
	if dbInfo.Name == "" {
		return nil, fmt.Errorf("DB_NAME is required")
	}

	return dbInfo, nil
}

// getEnv は環境変数を取得し、存在しない場合はデフォルト値を返す
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
