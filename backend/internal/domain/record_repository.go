package domain

import "context"

// RecordRepository はレコードリポジトリのインターフェース
type RecordRepository interface {
	// Create は新しいレコードを作成する
	Create(ctx context.Context, record *Record) (*Record, error)

	// FindByID は指定されたIDのレコードを取得する
	FindByID(ctx context.Context, id int) (*Record, error)

	// FindAll はレコードを取得する（ページネーション対応）
	// num: 取得件数
	// offset: オフセット
	// yyyymm: 年月フィルタ（例: "202501"）、空文字列の場合はフィルタしない
	// categoryID: カテゴリIDフィルタ、0の場合はフィルタしない
	FindAll(ctx context.Context, num, offset int, yyyymm string, categoryID int) ([]*Record, error)

	// Count は条件に一致するレコードの総数を取得する
	// yyyymm, categoryID が指定された場合はそれらでフィルタした件数を返す
	Count(ctx context.Context, yyyymm string, categoryID int) (int, error)

	// Delete は指定されたIDのレコードを削除する
	Delete(ctx context.Context, id int) error

	// GetAvailablePeriods はDBに登録されているレコードのYYYYMMとFY(年度)の一覧を取得する
	// 返される配列はいずれも新しい順にソートされている
	GetAvailablePeriods(ctx context.Context) (yyyymm []string, fy []string, err error)

	// GetYearSummary は指定された会計年度のカテゴリ別サマリーを取得する
	// year: 会計年度（例: 2024 → 2024年4月〜2025年3月）
	GetYearSummary(ctx context.Context, year int) ([]*CategoryYearSummary, error)
}
