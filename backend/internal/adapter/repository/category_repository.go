package repository

import (
	"context"
	"time"

	"github.com/azuki774/mawinter/internal/domain"
	"gorm.io/gorm"
)

// CategoryModel はCategoryテーブルのGORMモデル
type CategoryModel struct {
	ID           int       `gorm:"column:id;primaryKey;autoIncrement"`
	CategoryID   int       `gorm:"column:category_id;not null;uniqueIndex"`
	Name         string    `gorm:"column:name"`
	CategoryType int       `gorm:"column:category_type;not null;default:0"`
	CreatedAt    time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

// TableName はテーブル名を指定する
func (CategoryModel) TableName() string {
	return "Category"
}

// ToDomain はGORMモデルをドメインエンティティに変換する
func (m *CategoryModel) ToDomain() *domain.Category {
	return &domain.Category{
		ID:           m.ID,
		CategoryID:   m.CategoryID,
		Name:         m.Name,
		CategoryType: domain.CategoryType(m.CategoryType),
	}
}

// CategoryRepository はカテゴリリポジトリの実装
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository はCategoryRepositoryを生成する
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

// FindAll は全てのカテゴリを取得する
func (r *CategoryRepository) FindAll(ctx context.Context) ([]*domain.Category, error) {
	var models []*CategoryModel
	if err := r.db.WithContext(ctx).Order("category_id").Find(&models).Error; err != nil {
		return nil, err
	}

	categories := make([]*domain.Category, len(models))
	for i, model := range models {
		categories[i] = model.ToDomain()
	}

	return categories, nil
}
