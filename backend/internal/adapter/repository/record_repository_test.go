package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/azuki774/mawinter/internal/domain"
)

func TestRecordModel_ToDomain(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name         string
		model        *RecordModel
		categoryName string
		want         *domain.Record
	}{
		{
			name: "正常系: モデルからドメインへの変換",
			model: &RecordModel{
				ID:         1,
				CategoryID: 210,
				Datetime:   now,
				From:       "test-from",
				Type:       "test-type",
				Price:      1234,
				Memo:       "test-memo",
			},
			categoryName: "食費",
			want: &domain.Record{
				ID:           1,
				CategoryID:   210,
				CategoryName: "食費",
				Datetime:     now,
				From:         "test-from",
				Type:         "test-type",
				Price:        1234,
				Memo:         "test-memo",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record := tt.model.ToDomain(tt.categoryName)

			if record.ID != tt.want.ID {
				t.Errorf("expected ID %d, got %d", tt.want.ID, record.ID)
			}
			if record.CategoryID != tt.want.CategoryID {
				t.Errorf("expected CategoryID %d, got %d", tt.want.CategoryID, record.CategoryID)
			}
			if record.CategoryName != tt.want.CategoryName {
				t.Errorf("expected CategoryName '%s', got '%s'", tt.want.CategoryName, record.CategoryName)
			}
			if !record.Datetime.Equal(tt.want.Datetime) {
				t.Errorf("expected Datetime %v, got %v", tt.want.Datetime, record.Datetime)
			}
			if record.From != tt.want.From {
				t.Errorf("expected From '%s', got '%s'", tt.want.From, record.From)
			}
			if record.Type != tt.want.Type {
				t.Errorf("expected Type '%s', got '%s'", tt.want.Type, record.Type)
			}
			if record.Price != tt.want.Price {
				t.Errorf("expected Price %d, got %d", tt.want.Price, record.Price)
			}
			if record.Memo != tt.want.Memo {
				t.Errorf("expected Memo '%s', got '%s'", tt.want.Memo, record.Memo)
			}
		})
	}
}

func TestRecordModel_FromDomain(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name   string
		record *domain.Record
		want   *RecordModel
	}{
		{
			name: "正常系: ドメインからモデルへの変換",
			record: &domain.Record{
				ID:           1,
				CategoryID:   210,
				CategoryName: "食費",
				Datetime:     now,
				From:         "test-from",
				Type:         "test-type",
				Price:        1234,
				Memo:         "test-memo",
			},
			want: &RecordModel{
				ID:         1,
				CategoryID: 210,
				Datetime:   now,
				From:       "test-from",
				Type:       "test-type",
				Price:      1234,
				Memo:       "test-memo",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := &RecordModel{}
			model.FromDomain(tt.record)

			if model.ID != tt.want.ID {
				t.Errorf("expected ID %d, got %d", tt.want.ID, model.ID)
			}
			if model.CategoryID != tt.want.CategoryID {
				t.Errorf("expected CategoryID %d, got %d", tt.want.CategoryID, model.CategoryID)
			}
			if !model.Datetime.Equal(tt.want.Datetime) {
				t.Errorf("expected Datetime %v, got %v", tt.want.Datetime, model.Datetime)
			}
			if model.From != tt.want.From {
				t.Errorf("expected From '%s', got '%s'", tt.want.From, model.From)
			}
			if model.Type != tt.want.Type {
				t.Errorf("expected Type '%s', got '%s'", tt.want.Type, model.Type)
			}
			if model.Price != tt.want.Price {
				t.Errorf("expected Price %d, got %d", tt.want.Price, model.Price)
			}
			if model.Memo != tt.want.Memo {
				t.Errorf("expected Memo '%s', got '%s'", tt.want.Memo, model.Memo)
			}
		})
	}
}

func TestRecordModel_TableName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "正常系: テーブル名を取得",
			want: "Record",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := RecordModel{}
			if got := model.TableName(); got != tt.want {
				t.Errorf("TableName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRecordRepository(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "正常系: リポジトリが生成される",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewRecordRepository(nil)
			if repo == nil {
				t.Error("expected non-nil repository")
			}
		})
	}
}

func TestRecordRepository_Create(t *testing.T) {
	gormDB, mock := setupMockDB(t)

	now := time.Date(2025, 10, 21, 0, 0, 0, 0, time.UTC)
	record := &domain.Record{
		CategoryID: 210,
		Datetime:   now,
		From:       "test-from",
		Type:       "test-type",
		Price:      1234,
		Memo:       "test-memo",
	}

	// INSERT クエリのモック（created_at, updated_atは自動追加されない）
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `Record`")).
		WithArgs(210, "test-from", "test-type", 1234, "test-memo", now).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// カテゴリ名取得のSELECTクエリのモック
	categoryRows := sqlmock.NewRows([]string{"id", "category_id", "name", "category_type", "created_at", "updated_at"}).
		AddRow(1, 210, "食費", 2, time.Now(), time.Now())
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `Category` WHERE category_id = ? ORDER BY `Category`.`id` LIMIT ?")).
		WithArgs(210, 1).
		WillReturnRows(categoryRows)

	// テスト実行
	repo := NewRecordRepository(gormDB)
	result, err := repo.Create(context.Background(), record)

	// 検証
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != 1 {
		t.Errorf("expected ID 1, got %d", result.ID)
	}
	if result.CategoryName != "食費" {
		t.Errorf("expected CategoryName '食費', got '%s'", result.CategoryName)
	}

	// 全ての期待が満たされたか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestRecordRepository_FindByID(t *testing.T) {
	gormDB, mock := setupMockDB(t)

	now := time.Date(2025, 10, 21, 0, 0, 0, 0, time.UTC)

	// SELECT クエリのモック
	recordRows := sqlmock.NewRows([]string{"id", "category_id", "datetime", "from", "type", "price", "memo", "created_at", "updated_at"}).
		AddRow(1, 210, now, "test-from", "test-type", 1234, "test-memo", time.Now(), time.Now())
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `Record` WHERE id = ? ORDER BY `Record`.`id` LIMIT ?")).
		WithArgs(1, 1).
		WillReturnRows(recordRows)

	// カテゴリ名取得のSELECTクエリのモック
	categoryRows := sqlmock.NewRows([]string{"id", "category_id", "name", "category_type", "created_at", "updated_at"}).
		AddRow(1, 210, "食費", 2, time.Now(), time.Now())
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `Category` WHERE category_id = ? ORDER BY `Category`.`id` LIMIT ?")).
		WithArgs(210, 1).
		WillReturnRows(categoryRows)

	// テスト実行
	repo := NewRecordRepository(gormDB)
	result, err := repo.FindByID(context.Background(), 1)

	// 検証
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != 1 {
		t.Errorf("expected ID 1, got %d", result.ID)
	}
	if result.CategoryID != 210 {
		t.Errorf("expected CategoryID 210, got %d", result.CategoryID)
	}
	if result.CategoryName != "食費" {
		t.Errorf("expected CategoryName '食費', got '%s'", result.CategoryName)
	}

	// 全ての期待が満たされたか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestRecordRepository_FindAll(t *testing.T) {
	gormDB, mock := setupMockDB(t)

	now := time.Date(2025, 10, 21, 0, 0, 0, 0, time.UTC)

	// SELECT クエリのモック
	recordRows := sqlmock.NewRows([]string{"id", "category_id", "datetime", "from", "type", "price", "memo", "created_at", "updated_at"}).
		AddRow(2, 210, now, "test-from-2", "test-type", 5678, "test-memo-2", time.Now(), time.Now()).
		AddRow(1, 210, now, "test-from-1", "test-type", 1234, "test-memo-1", time.Now(), time.Now())
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `Record` ORDER BY id DESC LIMIT ?")).
		WithArgs(20).
		WillReturnRows(recordRows)

	// カテゴリ一覧取得のSELECTクエリのモック
	categoryRows := sqlmock.NewRows([]string{"id", "category_id", "name", "category_type", "created_at", "updated_at"}).
		AddRow(1, 210, "食費", 2, time.Now(), time.Now())
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `Category`")).
		WillReturnRows(categoryRows)

	// テスト実行
	repo := NewRecordRepository(gormDB)
	results, err := repo.FindAll(context.Background(), 20, 0, "", 0)

	// 検証
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 records, got %d", len(results))
	}

	// ID降順で取得されることを確認
	if results[0].ID != 2 {
		t.Errorf("expected first record ID 2, got %d", results[0].ID)
	}
	if results[1].ID != 1 {
		t.Errorf("expected second record ID 1, got %d", results[1].ID)
	}

	// カテゴリ名が取得されることを確認
	if results[0].CategoryName != "食費" {
		t.Errorf("expected CategoryName '食費', got '%s'", results[0].CategoryName)
	}

	// 全ての期待が満たされたか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestRecordRepository_FindAll_WithFilters(t *testing.T) {
	gormDB, mock := setupMockDB(t)

	now := time.Date(2025, 10, 21, 0, 0, 0, 0, time.UTC)

	// YYYYMMとcategory_idでフィルタするSELECTクエリのモック
	recordRows := sqlmock.NewRows([]string{"id", "category_id", "datetime", "from", "type", "price", "memo", "created_at", "updated_at"}).
		AddRow(1, 210, now, "test-from", "test-type", 1234, "test-memo", time.Now(), time.Now())
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `Record` WHERE (datetime >= ? AND datetime < DATE_ADD(?, INTERVAL 1 MONTH)) AND category_id = ? ORDER BY id DESC LIMIT ? OFFSET ?")).
		WithArgs("2025-10-01", "2025-10-01", 210, 10, 5).
		WillReturnRows(recordRows)

	// カテゴリ一覧取得のSELECTクエリのモック
	categoryRows := sqlmock.NewRows([]string{"id", "category_id", "name", "category_type", "created_at", "updated_at"}).
		AddRow(1, 210, "食費", 2, time.Now(), time.Now())
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `Category`")).
		WillReturnRows(categoryRows)

	// テスト実行
	repo := NewRecordRepository(gormDB)
	results, err := repo.FindAll(context.Background(), 10, 5, "202510", 210)

	// 検証
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 record, got %d", len(results))
	}

	// 全ての期待が満たされたか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestRecordRepository_Count(t *testing.T) {
	gormDB, mock := setupMockDB(t)

	// COUNT クエリのモック
	countRows := sqlmock.NewRows([]string{"count"}).
		AddRow(42)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `Record`")).
		WillReturnRows(countRows)

	// テスト実行
	repo := NewRecordRepository(gormDB)
	count, err := repo.Count(context.Background(), "", 0)

	// 検証
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if count != 42 {
		t.Errorf("expected count 42, got %d", count)
	}

	// 全ての期待が満たされたか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestRecordRepository_Count_WithFilters(t *testing.T) {
	gormDB, mock := setupMockDB(t)

	// フィルタ付きCOUNT クエリのモック
	countRows := sqlmock.NewRows([]string{"count"}).
		AddRow(10)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `Record` WHERE (datetime >= ? AND datetime < DATE_ADD(?, INTERVAL 1 MONTH)) AND category_id = ?")).
		WithArgs("2025-10-01", "2025-10-01", 210).
		WillReturnRows(countRows)

	// テスト実行
	repo := NewRecordRepository(gormDB)
	count, err := repo.Count(context.Background(), "202510", 210)

	// 検証
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if count != 10 {
		t.Errorf("expected count 10, got %d", count)
	}

	// 全ての期待が満たされたか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestRecordRepository_Delete(t *testing.T) {
	gormDB, mock := setupMockDB(t)

	// DELETE クエリのモック
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `Record` WHERE `Record`.`id` = ?")).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// テスト実行
	repo := NewRecordRepository(gormDB)
	err := repo.Delete(context.Background(), 1)

	// 検証
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// 全ての期待が満たされたか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestRecordRepository_Delete_NotFound(t *testing.T) {
	gormDB, mock := setupMockDB(t)

	// DELETE クエリのモック（削除対象が0件）
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `Record` WHERE `Record`.`id` = ?")).
		WithArgs(999).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	// テスト実行
	repo := NewRecordRepository(gormDB)
	err := repo.Delete(context.Background(), 999)

	// 検証
	if err == nil {
		t.Error("expected error for not found record, got nil")
	}

	// 全ての期待が満たされたか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
