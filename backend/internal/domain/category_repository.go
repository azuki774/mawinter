package domain

import "context"

// CategoryRepository はカテゴリリポジトリのインターフェース
type CategoryRepository interface {
	// FindAll は全てのカテゴリを取得する
	FindAll(ctx context.Context) ([]*Category, error)
}
