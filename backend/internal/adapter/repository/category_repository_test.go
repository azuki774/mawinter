package repository

import (
	"testing"

	"github.com/azuki774/mawinter/internal/domain"
)

func TestCategoryModel_ToDomain(t *testing.T) {
	model := &CategoryModel{
		ID:           1,
		CategoryID:   100,
		Name:         "月給",
		CategoryType: 1,
	}

	category := model.ToDomain()

	if category.ID != 1 {
		t.Errorf("expected ID 1, got %d", category.ID)
	}
	if category.CategoryID != 100 {
		t.Errorf("expected CategoryID 100, got %d", category.CategoryID)
	}
	if category.Name != "月給" {
		t.Errorf("expected Name '月給', got '%s'", category.Name)
	}
	if category.CategoryType != domain.CategoryTypeIncome {
		t.Errorf("expected CategoryType Income, got %v", category.CategoryType)
	}
}

func TestCategoryModel_TableName(t *testing.T) {
	model := CategoryModel{}
	if model.TableName() != "Category" {
		t.Errorf("expected table name 'Category', got '%s'", model.TableName())
	}
}

// 注: 実際のデータベース接続を使った統合テストは、
// CI/CD環境でデータベースが利用可能な場合に実行されるべきです。
// ここでは基本的な単体テストのみを含めています。

func TestNewCategoryRepository(t *testing.T) {
	// nilチェック
	repo := NewCategoryRepository(nil)
	if repo == nil {
		t.Error("expected non-nil repository")
	}
}

// FindAllの統合テストの例（実際のDBが必要）
// この関数は実際のDBが利用可能な環境でのみ実行されます
func TestCategoryRepository_FindAll_Integration(t *testing.T) {
	// DB接続がない場合はスキップ
	t.Skip("Integration test requires database connection")

	// 実際の統合テストの実装例:
	// db, err := gorm.Open(...)
	// if err != nil {
	//     t.Fatal(err)
	// }
	//
	// repo := NewCategoryRepository(db)
	// categories, err := repo.FindAll(context.Background())
	// if err != nil {
	//     t.Fatal(err)
	// }
	//
	// if len(categories) == 0 {
	//     t.Error("expected at least one category")
	// }
}

func TestCategoryRepository_FindAll_Context(t *testing.T) {
	// この関数は実際のDBが必要なのでスキップ
	t.Skip("Integration test requires database connection")

	// 実際の実装例（コンテキストのキャンセルテスト）:
	// ctx, cancel := context.WithCancel(context.Background())
	// cancel() // すぐにキャンセル
	//
	// repo := NewCategoryRepository(db)
	// _, err := repo.FindAll(ctx)
	// if err == nil {
	//     t.Error("expected error due to cancelled context")
	// }
}
