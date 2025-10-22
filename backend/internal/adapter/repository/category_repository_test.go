package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/azuki774/mawinter/internal/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestCategoryModel_ToDomain(t *testing.T) {
	tests := []struct {
		name  string
		model *CategoryModel
		want  *domain.Category
	}{
		{
			name: "正常系: モデルからドメインへの変換",
			model: &CategoryModel{
				ID:           1,
				CategoryID:   100,
				Name:         "月給",
				CategoryType: 1,
			},
			want: &domain.Category{
				ID:           1,
				CategoryID:   100,
				Name:         "月給",
				CategoryType: domain.CategoryTypeIncome,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category := tt.model.ToDomain()

			if category.ID != tt.want.ID {
				t.Errorf("expected ID %d, got %d", tt.want.ID, category.ID)
			}
			if category.CategoryID != tt.want.CategoryID {
				t.Errorf("expected CategoryID %d, got %d", tt.want.CategoryID, category.CategoryID)
			}
			if category.Name != tt.want.Name {
				t.Errorf("expected Name '%s', got '%s'", tt.want.Name, category.Name)
			}
			if category.CategoryType != tt.want.CategoryType {
				t.Errorf("expected CategoryType %v, got %v", tt.want.CategoryType, category.CategoryType)
			}
		})
	}
}

func TestCategoryModel_TableName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "正常系: テーブル名を取得",
			want: "Category",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := CategoryModel{}
			if got := model.TableName(); got != tt.want {
				t.Errorf("TableName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCategoryRepository(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "正常系: リポジトリが生成される",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewCategoryRepository(nil)
			if repo == nil {
				t.Error("expected non-nil repository")
			}
		})
	}
}

// setupMockDB はテスト用のモックDBを作成する
func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm DB: %v", err)
	}

	return gormDB, mock
}

func TestCategoryRepository_FindAll(t *testing.T) {
	gormDB, mock := setupMockDB(t)

	// モックデータの準備
	rows := sqlmock.NewRows([]string{"id", "category_id", "name", "category_type", "created_at", "updated_at"}).
		AddRow(1, 100, "月給", 1, time.Now(), time.Now()).
		AddRow(2, 200, "家賃", 2, time.Now(), time.Now()).
		AddRow(3, 700, "NISA入出金", 4, time.Now(), time.Now())

	// 期待するクエリを設定
	mock.ExpectQuery("^SELECT \\* FROM `Category` ORDER BY category_id$").
		WillReturnRows(rows)

	// テスト実行
	repo := NewCategoryRepository(gormDB)
	categories, err := repo.FindAll(context.Background())

	// 検証
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(categories) != 3 {
		t.Errorf("expected 3 categories, got %d", len(categories))
	}

	// 最初のカテゴリの確認
	if categories[0].CategoryID != 100 {
		t.Errorf("expected category_id 100, got %d", categories[0].CategoryID)
	}
	if categories[0].Name != "月給" {
		t.Errorf("expected name '月給', got '%s'", categories[0].Name)
	}
	if categories[0].CategoryType != domain.CategoryTypeIncome {
		t.Errorf("expected category_type Income, got %v", categories[0].CategoryType)
	}

	// 2番目のカテゴリの確認
	if categories[1].CategoryType != domain.CategoryTypeOutgoing {
		t.Errorf("expected category_type Outgoing, got %v", categories[1].CategoryType)
	}

	// 3番目のカテゴリの確認
	if categories[2].CategoryType != domain.CategoryTypeInvesting {
		t.Errorf("expected category_type Investing, got %v", categories[2].CategoryType)
	}

	// 全ての期待が満たされたか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestCategoryRepository_FindAll_Error(t *testing.T) {
	gormDB, mock := setupMockDB(t)

	// エラーを返すモック
	mock.ExpectQuery("^SELECT \\* FROM `Category` ORDER BY category_id$").
		WillReturnError(context.DeadlineExceeded)

	// テスト実行
	repo := NewCategoryRepository(gormDB)
	categories, err := repo.FindAll(context.Background())

	// 検証
	if err == nil {
		t.Error("expected error, got nil")
	}

	if categories != nil {
		t.Errorf("expected nil categories, got %v", categories)
	}

	// 全ての期待が満たされたか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
