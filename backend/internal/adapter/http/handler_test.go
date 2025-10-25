package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/azuki774/mawinter/api"
	"github.com/azuki774/mawinter/internal/application"
	"github.com/azuki774/mawinter/internal/domain"
	"github.com/azuki774/mawinter/pkg/config"
	"github.com/gin-gonic/gin"
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

// mockRecordRepository はテスト用のモックリポジトリ
type mockRecordRepository struct {
	records      []*domain.Record
	err          error
	findByIDFunc func(ctx context.Context, id int) (*domain.Record, error)
	deleteFunc   func(ctx context.Context, id int) error
}

func (m *mockRecordRepository) Create(ctx context.Context, record *domain.Record) (*domain.Record, error) {
	return nil, nil
}

func (m *mockRecordRepository) FindByID(ctx context.Context, id int) (*domain.Record, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockRecordRepository) FindAll(ctx context.Context, num, offset int, yyyymm string, categoryID int) ([]*domain.Record, error) {
	return nil, nil
}

func (m *mockRecordRepository) Count(ctx context.Context, yyyymm string, categoryID int) (int, error) {
	return 0, nil
}

func (m *mockRecordRepository) Delete(ctx context.Context, id int) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func TestGetV3Categories(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		mockRepo           *mockCategoryRepository
		wantStatusCode     int
		wantCategoryCount  int
		checkResponseFunc  func(t *testing.T, response []api.Category)
	}{
		{
			name: "正常系: カテゴリ一覧を取得できる",
			mockRepo: &mockCategoryRepository{
				categories: []*domain.Category{
					{ID: 1, CategoryID: 100, Name: "月給", CategoryType: domain.CategoryTypeIncome},
					{ID: 2, CategoryID: 200, Name: "家賃", CategoryType: domain.CategoryTypeOutgoing},
					{ID: 3, CategoryID: 700, Name: "NISA入出金", CategoryType: domain.CategoryTypeInvesting},
				},
			},
			wantStatusCode:    http.StatusOK,
			wantCategoryCount: 3,
			checkResponseFunc: func(t *testing.T, response []api.Category) {
				if response[0].CategoryId != 100 {
					t.Errorf("expected category_id 100, got %d", response[0].CategoryId)
				}
				if response[0].CategoryName != "月給" {
					t.Errorf("expected category_name '月給', got '%s'", response[0].CategoryName)
				}
				if response[0].CategoryType != "income" {
					t.Errorf("expected category_type 'income', got '%s'", response[0].CategoryType)
				}
				if response[1].CategoryType != "outgoing" {
					t.Errorf("expected category_type 'outgoing', got '%s'", response[1].CategoryType)
				}
				if response[2].CategoryType != "investing" {
					t.Errorf("expected category_type 'investing', got '%s'", response[2].CategoryType)
				}
			},
		},
		{
			name: "異常系: リポジトリエラー",
			mockRepo: &mockCategoryRepository{
				err: context.DeadlineExceeded,
			},
			wantStatusCode:    http.StatusInternalServerError,
			wantCategoryCount: 0,
			checkResponseFunc: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			categoryService := application.NewCategoryService(tt.mockRepo)
			recordService := application.NewRecordService(&mockRecordRepository{})
			dbInfo := &config.DBInfo{}
			server := NewServer("localhost", 8080, "v1.0.0", "abc123", "20250101", dbInfo, categoryService, recordService)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v3/categories", nil)
			server.router.ServeHTTP(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("expected status code %d, got %d", tt.wantStatusCode, w.Code)
			}

			if tt.wantStatusCode == http.StatusOK {
				var response []api.Category
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}

				if len(response) != tt.wantCategoryCount {
					t.Errorf("expected %d categories, got %d", tt.wantCategoryCount, len(response))
				}

				if tt.checkResponseFunc != nil {
					tt.checkResponseFunc(t, response)
				}
			}
		})
	}
}

func TestGetV3RecordId(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name              string
		recordID          int
		mockRepo          *mockRecordRepository
		wantStatusCode    int
		checkResponseFunc func(t *testing.T, response api.Record)
	}{
		{
			name:     "正常系: IDでレコードを取得できる",
			recordID: 1,
			mockRepo: &mockRecordRepository{
				findByIDFunc: func(ctx context.Context, id int) (*domain.Record, error) {
					return &domain.Record{
						ID:           1,
						CategoryID:   100,
						CategoryName: "月給",
						Price:        300000,
						From:         "会社A",
						Type:         "salary",
						Memo:         "テストメモ",
					}, nil
				},
			},
			wantStatusCode: http.StatusOK,
			checkResponseFunc: func(t *testing.T, response api.Record) {
				if response.Id != 1 {
					t.Errorf("expected id 1, got %d", response.Id)
				}
				if response.CategoryId != 100 {
					t.Errorf("expected category_id 100, got %d", response.CategoryId)
				}
				if response.Price != 300000 {
					t.Errorf("expected price 300000, got %d", response.Price)
				}
			},
		},
		{
			name:     "異常系: レコードが見つからない",
			recordID: 999,
			mockRepo: &mockRecordRepository{
				findByIDFunc: func(ctx context.Context, id int) (*domain.Record, error) {
					return nil, context.DeadlineExceeded
				},
			},
			wantStatusCode:    http.StatusNotFound,
			checkResponseFunc: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			categoryService := application.NewCategoryService(&mockCategoryRepository{})
			recordService := application.NewRecordService(tt.mockRepo)
			dbInfo := &config.DBInfo{}
			server := NewServer("localhost", 8080, "v1.0.0", "abc123", "20250101", dbInfo, categoryService, recordService)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v3/record/1", nil)
			server.router.ServeHTTP(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("expected status code %d, got %d", tt.wantStatusCode, w.Code)
			}

			if tt.wantStatusCode == http.StatusOK && tt.checkResponseFunc != nil {
				var response api.Record
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				tt.checkResponseFunc(t, response)
			}
		})
	}
}

func TestDeleteV3RecordId(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		recordID       int
		mockRepo       *mockRecordRepository
		wantStatusCode int
	}{
		{
			name:     "正常系: IDでレコードを削除できる",
			recordID: 1,
			mockRepo: &mockRecordRepository{
				deleteFunc: func(ctx context.Context, id int) error {
					return nil
				},
			},
			wantStatusCode: http.StatusNoContent,
		},
		{
			name:     "異常系: レコードが見つからない",
			recordID: 999,
			mockRepo: &mockRecordRepository{
				deleteFunc: func(ctx context.Context, id int) error {
					return context.DeadlineExceeded
				},
			},
			wantStatusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			categoryService := application.NewCategoryService(&mockCategoryRepository{})
			recordService := application.NewRecordService(tt.mockRepo)
			dbInfo := &config.DBInfo{}
			server := NewServer("localhost", 8080, "v1.0.0", "abc123", "20250101", dbInfo, categoryService, recordService)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/api/v3/record/1", nil)
			server.router.ServeHTTP(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("expected status code %d, got %d", tt.wantStatusCode, w.Code)
			}
		})
	}
}
