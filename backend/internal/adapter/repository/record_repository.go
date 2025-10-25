package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/azuki774/mawinter/internal/domain"
	"gorm.io/gorm"
)

// RecordModel はRecordテーブルのGORMモデル
type RecordModel struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement"`
	CategoryID int       `gorm:"column:category_id;not null"`
	Datetime   time.Time `gorm:"column:datetime;not null;default:CURRENT_TIMESTAMP"`
	From       string    `gorm:"column:from;not null"`
	Type       string    `gorm:"column:type;not null"`
	Price      int       `gorm:"column:price;not null"`
	Memo       string    `gorm:"column:memo;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

// TableName はテーブル名を指定する
func (RecordModel) TableName() string {
	return "Record"
}

// ToDomain はGORMモデルをドメインエンティティに変換する
// CategoryNameは別途取得が必要
func (m *RecordModel) ToDomain(categoryName string) *domain.Record {
	return &domain.Record{
		ID:           m.ID,
		CategoryID:   m.CategoryID,
		CategoryName: categoryName,
		Datetime:     m.Datetime,
		From:         m.From,
		Type:         m.Type,
		Price:        m.Price,
		Memo:         m.Memo,
	}
}

// FromDomain はドメインエンティティからGORMモデルに変換する
func (m *RecordModel) FromDomain(record *domain.Record) {
	m.ID = record.ID
	m.CategoryID = record.CategoryID
	m.Datetime = record.Datetime
	m.From = record.From
	m.Type = record.Type
	m.Price = record.Price
	m.Memo = record.Memo
}

// RecordRepository はレコードリポジトリの実装
type RecordRepository struct {
	db *gorm.DB
}

// NewRecordRepository はRecordRepositoryを生成する
func NewRecordRepository(db *gorm.DB) *RecordRepository {
	return &RecordRepository{
		db: db,
	}
}

// Create は新しいレコードを作成する
func (r *RecordRepository) Create(ctx context.Context, record *domain.Record) (*domain.Record, error) {
	model := &RecordModel{}
	model.FromDomain(record)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}

	// カテゴリ名を取得
	var category CategoryModel
	if err := r.db.WithContext(ctx).Where("category_id = ?", model.CategoryID).First(&category).Error; err != nil {
		return nil, err
	}

	return model.ToDomain(category.Name), nil
}

// FindByID は指定されたIDのレコードを取得する
func (r *RecordRepository) FindByID(ctx context.Context, id int) (*domain.Record, error) {
	var model RecordModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}

	// カテゴリ名を取得
	var category CategoryModel
	if err := r.db.WithContext(ctx).Where("category_id = ?", model.CategoryID).First(&category).Error; err != nil {
		return nil, err
	}

	return model.ToDomain(category.Name), nil
}

// FindAll はレコードを取得する（ページネーション対応）
func (r *RecordRepository) FindAll(ctx context.Context, num, offset int, yyyymm string, categoryID int) ([]*domain.Record, error) {
	query := r.db.WithContext(ctx)

	// YYYYMMフィルタ
	if yyyymm != "" {
		// YYYYMMをYYYY-MM-01形式に変換
		if len(yyyymm) != 6 {
			return nil, fmt.Errorf("invalid yyyymm format: %s", yyyymm)
		}
		startDate := yyyymm[:4] + "-" + yyyymm[4:6] + "-01"
		// 次の月の1日を計算（その月の終わりまで）
		year := yyyymm[:4]
		month := yyyymm[4:6]
		endDate := fmt.Sprintf("%s-%s-01", year, month)
		query = query.Where("datetime >= ? AND datetime < DATE_ADD(?, INTERVAL 1 MONTH)", startDate, endDate)
	}

	// カテゴリIDフィルタ
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	// ID降順で取得
	query = query.Order("id DESC").Limit(num).Offset(offset)

	var models []*RecordModel
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	// カテゴリIDのマップを作成（一度に全カテゴリを取得）
	var categories []*CategoryModel
	if err := r.db.WithContext(ctx).Find(&categories).Error; err != nil {
		return nil, err
	}

	categoryMap := make(map[int]string)
	for _, cat := range categories {
		categoryMap[cat.CategoryID] = cat.Name
	}

	// ドメインエンティティに変換
	records := make([]*domain.Record, len(models))
	for i, model := range models {
		categoryName := categoryMap[model.CategoryID]
		records[i] = model.ToDomain(categoryName)
	}

	return records, nil
}

// Count は条件に一致するレコードの総数を取得する
func (r *RecordRepository) Count(ctx context.Context, yyyymm string, categoryID int) (int, error) {
	query := r.db.WithContext(ctx).Model(&RecordModel{})

	// YYYYMMフィルタ
	if yyyymm != "" {
		if len(yyyymm) != 6 {
			return 0, fmt.Errorf("invalid yyyymm format: %s", yyyymm)
		}
		startDate := yyyymm[:4] + "-" + yyyymm[4:6] + "-01"
		year := yyyymm[:4]
		month := yyyymm[4:6]
		endDate := fmt.Sprintf("%s-%s-01", year, month)
		query = query.Where("datetime >= ? AND datetime < DATE_ADD(?, INTERVAL 1 MONTH)", startDate, endDate)
	}

	// カテゴリIDフィルタ
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

// Delete は指定されたIDのレコードを削除する
func (r *RecordRepository) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&RecordModel{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetAvailablePeriods はDBに登録されているレコードのYYYYMMとFY(年度)の一覧を取得する
// 年度は4月始まり(4月〜翌年3月)で計算される
// 返される配列はいずれも新しい順にソートされている
func (r *RecordRepository) GetAvailablePeriods(ctx context.Context) ([]string, []string, error) {
	// YYYYMMを取得（重複なし、降順）
	type YYYYMMResult struct {
		YYYYMM string
	}
	var yyyymmResults []YYYYMMResult

	// DATE_FORMAT を使ってYYYYMMを抽出し、DISTINCTで重複排除
	if err := r.db.WithContext(ctx).
		Model(&RecordModel{}).
		Select("DISTINCT DATE_FORMAT(datetime, '%Y%m') as yyyymm").
		Order("yyyymm DESC").
		Scan(&yyyymmResults).Error; err != nil {
		return nil, nil, err
	}

	// YYYYMMのスライスを作成
	yyyymmSlice := make([]string, len(yyyymmResults))
	for i, result := range yyyymmResults {
		yyyymmSlice[i] = result.YYYYMM
	}

	// FY(年度)を計算
	// 年度は4月始まりなので、CASE文で計算
	// 1-3月は前年度、4-12月は当年度
	type FYResult struct {
		FY string
	}
	var fyResults []FYResult

	if err := r.db.WithContext(ctx).
		Model(&RecordModel{}).
		Select("DISTINCT CASE WHEN MONTH(datetime) BETWEEN 1 AND 3 THEN YEAR(datetime) - 1 ELSE YEAR(datetime) END as fy").
		Order("fy DESC").
		Scan(&fyResults).Error; err != nil {
		return nil, nil, err
	}

	// FYのスライスを作成
	fySlice := make([]string, len(fyResults))
	for i, result := range fyResults {
		fySlice[i] = result.FY
	}

	return yyyymmSlice, fySlice, nil
}

// GetYearSummary は指定された会計年度のカテゴリ別サマリーを取得する
// year: 会計年度（例: 2024 → 2024年4月〜2025年3月）
func (r *RecordRepository) GetYearSummary(ctx context.Context, year int) ([]*domain.CategoryYearSummary, error) {
	// 会計年度の開始日と終了日を計算
	// year=2024 → 2024-04-01 〜 2025-03-31
	startDate := fmt.Sprintf("%d-04-01", year)
	endDate := fmt.Sprintf("%d-04-01", year+1)

	// カテゴリ情報を全て取得
	var categories []*CategoryModel
	if err := r.db.WithContext(ctx).Find(&categories).Error; err != nil {
		return nil, err
	}

	// カテゴリIDからカテゴリ情報へのマップを作成
	categoryMap := make(map[int]*CategoryModel)
	for _, cat := range categories {
		categoryMap[cat.CategoryID] = cat
	}

	// 月別・カテゴリ別の集計を取得
	// 会計年度の4月を1とし、翌年3月を12とする
	type MonthlySum struct {
		CategoryID int
		FiscalMonth int  // 会計年度での月番号（1-12）
		TotalPrice  int
		Count       int
	}

	var monthlySums []MonthlySum

	// SQLクエリで月別・カテゴリ別に集計
	// 会計年度の月を計算: 4月=1, 5月=2, ..., 3月=12
	query := `
		SELECT
			category_id,
			CASE
				WHEN MONTH(datetime) >= 4 THEN MONTH(datetime) - 3
				ELSE MONTH(datetime) + 9
			END as fiscal_month,
			SUM(price) as total_price,
			COUNT(*) as count
		FROM Record
		WHERE datetime >= ? AND datetime < ?
		GROUP BY category_id, fiscal_month
		ORDER BY category_id, fiscal_month
	`

	if err := r.db.WithContext(ctx).Raw(query, startDate, endDate).Scan(&monthlySums).Error; err != nil {
		return nil, err
	}

	// カテゴリごとに集計データをまとめる
	summaryMap := make(map[int]*domain.CategoryYearSummary)

	for _, ms := range monthlySums {
		cat, exists := categoryMap[ms.CategoryID]
		if !exists {
			continue // カテゴリが見つからない場合はスキップ
		}

		// サマリーがまだ存在しない場合は作成
		if _, exists := summaryMap[ms.CategoryID]; !exists {
			summaryMap[ms.CategoryID] = &domain.CategoryYearSummary{
				CategoryID:   cat.CategoryID,
				CategoryName: cat.Name,
				CategoryType: domain.CategoryType(cat.CategoryType),
				Count:        0,
				Price:        [12]int{}, // ゼロ初期化
				Total:        0,
			}
		}

		summary := summaryMap[ms.CategoryID]

		// 会計年度の月番号（1-12）を配列インデックス（0-11）に変換
		monthIndex := ms.FiscalMonth - 1
		if monthIndex >= 0 && monthIndex < 12 {
			summary.Price[monthIndex] = ms.TotalPrice
			summary.Count += ms.Count
			summary.Total += ms.TotalPrice
		}
	}

	// マップからスライスに変換
	result := make([]*domain.CategoryYearSummary, 0, len(summaryMap))
	for _, summary := range summaryMap {
		result = append(result, summary)
	}

	return result, nil
}
