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
	tests := []struct {
		name       string
		mockRepo   *mockCategoryRepository
		wantErr    bool
		wantLength int
		checkFunc  func(t *testing.T, categories []*domain.Category)
	}{
		{
			name: "正常系: カテゴリ一覧を取得できる",
			mockRepo: &mockCategoryRepository{
				categories: []*domain.Category{
					{ID: 1, CategoryID: 100, Name: "月給", CategoryType: domain.CategoryTypeIncome},
					{ID: 2, CategoryID: 200, Name: "家賃", CategoryType: domain.CategoryTypeOutgoing},
				},
			},
			wantErr:    false,
			wantLength: 2,
			checkFunc: func(t *testing.T, categories []*domain.Category) {
				if categories[0].CategoryID != 100 {
					t.Errorf("expected category_id 100, got %d", categories[0].CategoryID)
				}
				if categories[1].CategoryID != 200 {
					t.Errorf("expected category_id 200, got %d", categories[1].CategoryID)
				}
			},
		},
		{
			name: "異常系: リポジトリエラー",
			mockRepo: &mockCategoryRepository{
				err: errors.New("database error"),
			},
			wantErr:    true,
			wantLength: 0,
			checkFunc:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			service := NewCategoryService(tt.mockRepo)

			categories, err := service.GetAllCategories(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(categories) != tt.wantLength {
				t.Errorf("expected %d categories, got %d", tt.wantLength, len(categories))
			}

			if tt.checkFunc != nil {
				tt.checkFunc(t, categories)
			}
		})
	}
}
