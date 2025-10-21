package application

import (
	"context"
	"errors"
	"testing"

	"github.com/azuki774/mawinter/internal/domain"
)

// mockCategoryRepository はテスト用のモックリポジトリ
type mockCategoryRepository struct {
	categories []*domain.Category
	err        error
}

func (m *mockCategoryRepository) FindAll(ctx context.Context) ([]*domain.Category, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.categories, nil
}

func TestCategoryService_GetAllCategories(t *testing.T) {
	ctx := context.Background()

	t.Run("正常系: カテゴリ一覧を取得できる", func(t *testing.T) {
		mockRepo := &mockCategoryRepository{
			categories: []*domain.Category{
				{ID: 1, CategoryID: 100, Name: "月給", CategoryType: domain.CategoryTypeIncome},
				{ID: 2, CategoryID: 200, Name: "家賃", CategoryType: domain.CategoryTypeOutgoing},
			},
		}
		service := NewCategoryService(mockRepo)

		categories, err := service.GetAllCategories(ctx)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if len(categories) != 2 {
			t.Errorf("expected 2 categories, got %d", len(categories))
		}

		if categories[0].CategoryID != 100 {
			t.Errorf("expected category_id 100, got %d", categories[0].CategoryID)
		}
		if categories[1].CategoryID != 200 {
			t.Errorf("expected category_id 200, got %d", categories[1].CategoryID)
		}
	})

	t.Run("異常系: リポジトリエラー", func(t *testing.T) {
		mockRepo := &mockCategoryRepository{
			err: errors.New("database error"),
		}
		service := NewCategoryService(mockRepo)

		_, err := service.GetAllCategories(ctx)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
