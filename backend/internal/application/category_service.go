package application

import (
	"context"

	"github.com/azuki774/mawinter/internal/domain"
)

// CategoryService はカテゴリに関するアプリケーションサービス
type CategoryService struct {
	repo domain.CategoryRepository
}

// NewCategoryService はCategoryServiceを生成する
func NewCategoryService(repo domain.CategoryRepository) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

// GetAllCategories は全てのカテゴリを取得する
func (s *CategoryService) GetAllCategories(ctx context.Context) ([]*domain.Category, error) {
	return s.repo.FindAll(ctx)
}
