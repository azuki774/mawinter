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
	records []*domain.Record
	err     error
}

func (m *mockRecordRepository) Create(ctx context.Context, record *domain.Record) (*domain.Record, error) {
	return nil, nil
}

func (m *mockRecordRepository) FindByID(ctx context.Context, id int) (*domain.Record, error) {
	return nil, nil
}

func (m *mockRecordRepository) FindAll(ctx context.Context, num, offset int, yyyymm string, categoryID int) ([]*domain.Record, error) {
	return nil, nil
}

func (m *mockRecordRepository) Count(ctx context.Context, yyyymm string, categoryID int) (int, error) {
	return 0, nil
}

func (m *mockRecordRepository) Delete(ctx context.Context, id int) error {
	return nil
}

func TestGetV3Categories(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("正常系: カテゴリ一覧を取得できる", func(t *testing.T) {
		// モックデータの準備
		mockRepo := &mockCategoryRepository{
			categories: []*domain.Category{
				{ID: 1, CategoryID: 100, Name: "月給", CategoryType: domain.CategoryTypeIncome},
				{ID: 2, CategoryID: 200, Name: "家賃", CategoryType: domain.CategoryTypeOutgoing},
				{ID: 3, CategoryID: 700, Name: "NISA入出金", CategoryType: domain.CategoryTypeInvesting},
			},
		}

		categoryService := application.NewCategoryService(mockRepo)
		recordService := application.NewRecordService(&mockRecordRepository{})

		// サーバの作成
		dbInfo := &config.DBInfo{}
		server := NewServer("localhost", 8080, "v1.0.0", "abc123", "20250101", dbInfo, categoryService, recordService)

		// テストリクエスト
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v3/categories", nil)
		server.router.ServeHTTP(w, req)

		// ステータスコードの確認
		if w.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
		}

		// レスポンスボディの確認
		var response []api.Category
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if len(response) != 3 {
			t.Errorf("expected 3 categories, got %d", len(response))
		}

		// 最初のカテゴリの確認
		if response[0].CategoryId != 100 {
			t.Errorf("expected category_id 100, got %d", response[0].CategoryId)
		}
		if response[0].CategoryName != "月給" {
			t.Errorf("expected category_name '月給', got '%s'", response[0].CategoryName)
		}
		if response[0].CategoryType != "income" {
			t.Errorf("expected category_type 'income', got '%s'", response[0].CategoryType)
		}

		// 2番目のカテゴリの確認
		if response[1].CategoryType != "outgoing" {
			t.Errorf("expected category_type 'outgoing', got '%s'", response[1].CategoryType)
		}

		// 3番目のカテゴリの確認
		if response[2].CategoryType != "investing" {
			t.Errorf("expected category_type 'investing', got '%s'", response[2].CategoryType)
		}
	})

	t.Run("異常系: リポジトリエラー", func(t *testing.T) {
		// エラーを返すモック
		mockRepo := &mockCategoryRepository{
			err: context.DeadlineExceeded,
		}

		categoryService := application.NewCategoryService(mockRepo)
		recordService := application.NewRecordService(&mockRecordRepository{})
		dbInfo := &config.DBInfo{}
		server := NewServer("localhost", 8080, "v1.0.0", "abc123", "20250101", dbInfo, categoryService, recordService)

		// テストリクエスト
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v3/categories", nil)
		server.router.ServeHTTP(w, req)

		// エラー時は500を返す
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}
