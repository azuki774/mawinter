package repository

import (
	"context"
	"fmt"
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

func TestRecordRepository_GetAvailablePeriods(t *testing.T) {
	tests := []struct {
		name           string
		mockRecords    []RecordModel
		wantYYYYMM     []string
		wantFY         []string
		wantErr        bool
		checkYYYYMMLen int
		checkFYLen     int
	}{
		{
			name: "正常系: 複数年月のレコードから期間を取得",
			mockRecords: []RecordModel{
				{Datetime: time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC)},  // FY2024, 202405
				{Datetime: time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)},   // FY2024, 202404
				{Datetime: time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC)},  // FY2023, 202403
				{Datetime: time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)}, // FY2023, 202312
				{Datetime: time.Date(2023, 5, 10, 0, 0, 0, 0, time.UTC)},  // FY2023, 202305
			},
			wantYYYYMM:     []string{"202405", "202404", "202403", "202312", "202305"},
			wantFY:         []string{"2024", "2023"},
			wantErr:        false,
			checkYYYYMMLen: 5,
			checkFYLen:     2,
		},
		{
			name: "正常系: 同じ年月の重複レコードは1つにまとめる",
			mockRecords: []RecordModel{
				{Datetime: time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)},
				{Datetime: time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC)},
				{Datetime: time.Date(2024, 5, 31, 0, 0, 0, 0, time.UTC)},
			},
			wantYYYYMM:     []string{"202405"},
			wantFY:         []string{"2024"},
			wantErr:        false,
			checkYYYYMMLen: 1,
			checkFYLen:     1,
		},
		{
			name:           "正常系: レコードが0件の場合は空配列を返す",
			mockRecords:    []RecordModel{},
			wantYYYYMM:     []string{},
			wantFY:         []string{},
			wantErr:        false,
			checkYYYYMMLen: 0,
			checkFYLen:     0,
		},
		{
			name: "正常系: 年度境界のテスト（3月と4月）",
			mockRecords: []RecordModel{
				{Datetime: time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)},  // FY2024
				{Datetime: time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC)}, // FY2023
			},
			wantYYYYMM:     []string{"202404", "202403"},
			wantFY:         []string{"2024", "2023"},
			wantErr:        false,
			checkYYYYMMLen: 2,
			checkFYLen:     2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock := setupMockDB(t)

			// YYYYMM取得クエリのモック
			yyyymmRows := sqlmock.NewRows([]string{"yyyymm"})
			yyyymmMap := make(map[string]bool)
			for _, record := range tt.mockRecords {
				yyyymm := record.Datetime.Format("200601")
				if !yyyymmMap[yyyymm] {
					yyyymmMap[yyyymm] = true
				}
			}
			// 降順でソート
			yyyymmList := make([]string, 0, len(yyyymmMap))
			for k := range yyyymmMap {
				yyyymmList = append(yyyymmList, k)
			}
			for i := 0; i < len(yyyymmList)-1; i++ {
				for j := i + 1; j < len(yyyymmList); j++ {
					if yyyymmList[i] < yyyymmList[j] {
						yyyymmList[i], yyyymmList[j] = yyyymmList[j], yyyymmList[i]
					}
				}
			}
			for _, ym := range yyyymmList {
				yyyymmRows.AddRow(ym)
			}
			mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT DATE_FORMAT(datetime, '%Y%m') as yyyymm FROM `Record` ORDER BY yyyymm DESC")).
				WillReturnRows(yyyymmRows)

			// FY取得クエリのモック
			fyRows := sqlmock.NewRows([]string{"fy"})
			fyMap := make(map[string]bool)
			for _, record := range tt.mockRecords {
				year := record.Datetime.Year()
				month := record.Datetime.Month()
				fyYear := year
				if month >= 1 && month <= 3 {
					fyYear = year - 1
				}
				fy := fmt.Sprintf("%d", fyYear)
				if !fyMap[fy] {
					fyMap[fy] = true
				}
			}
			// 降順でソート
			fyList := make([]string, 0, len(fyMap))
			for k := range fyMap {
				fyList = append(fyList, k)
			}
			for i := 0; i < len(fyList)-1; i++ {
				for j := i + 1; j < len(fyList); j++ {
					if fyList[i] < fyList[j] {
						fyList[i], fyList[j] = fyList[j], fyList[i]
					}
				}
			}
			for _, f := range fyList {
				fyRows.AddRow(f)
			}
			mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT CASE WHEN MONTH(datetime) BETWEEN 1 AND 3 THEN YEAR(datetime) - 1 ELSE YEAR(datetime) END as fy FROM `Record` ORDER BY fy DESC")).
				WillReturnRows(fyRows)

			// テスト実行
			repo := NewRecordRepository(gormDB)
			yyyymm, fy, err := repo.GetAvailablePeriods(context.Background())

			// エラーチェック
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAvailablePeriods() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 配列長のチェック
			if len(yyyymm) != tt.checkYYYYMMLen {
				t.Errorf("expected yyyymm length %d, got %d", tt.checkYYYYMMLen, len(yyyymm))
			}
			if len(fy) != tt.checkFYLen {
				t.Errorf("expected fy length %d, got %d", tt.checkFYLen, len(fy))
			}

			// 結果の検証
			if !tt.wantErr && len(tt.wantYYYYMM) > 0 {
				for i, want := range tt.wantYYYYMM {
					if i < len(yyyymm) {
						if yyyymm[i] != want {
							t.Errorf("yyyymm[%d] = %s, want %s", i, yyyymm[i], want)
						}
					}
				}
			}

			if !tt.wantErr && len(tt.wantFY) > 0 {
				for i, want := range tt.wantFY {
					if i < len(fy) {
						if fy[i] != want {
							t.Errorf("fy[%d] = %s, want %s", i, fy[i], want)
						}
					}
				}
			}

			// ソート順のチェック（新しい順）
			if len(yyyymm) > 1 {
				for i := 0; i < len(yyyymm)-1; i++ {
					if yyyymm[i] < yyyymm[i+1] {
						t.Errorf("yyyymm is not sorted in descending order: %v", yyyymm)
						break
					}
				}
			}

			if len(fy) > 1 {
				for i := 0; i < len(fy)-1; i++ {
					if fy[i] < fy[i+1] {
						t.Errorf("fy is not sorted in descending order: %v", fy)
						break
					}
				}
			}

			// 全ての期待が満たされたか確認
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestRecordRepository_GetYearSummary(t *testing.T) {
	tests := []struct {
		name          string
		year          int
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		checkFunc     func(t *testing.T, summaries []*domain.CategoryYearSummary)
	}{
		{
			name: "正常系: 会計年度のサマリーを取得できる",
			year: 2024,
			mockSetup: func(mock sqlmock.Sqlmock) {
				// カテゴリ取得のモック
				categoryRows := sqlmock.NewRows([]string{"id", "category_id", "name", "category_type"}).
					AddRow(1, 201, "食費", 2).  // outgoing = 2
					AddRow(2, 100, "月給", 1)   // income = 1
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `Category`")).
					WillReturnRows(categoryRows)

				// 月別集計のモック
				// 2024年度（2024年4月〜2025年3月）のデータ
				summaryRows := sqlmock.NewRows([]string{"category_id", "fiscal_month", "total_price", "count"}).
					AddRow(201, 1, 50000, 10).  // 4月
					AddRow(201, 2, 60000, 12).  // 5月
					AddRow(201, 12, 55000, 11). // 3月
					AddRow(100, 1, 300000, 1).  // 4月
					AddRow(100, 12, 310000, 1)  // 3月

				querySQL := `
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
				mock.ExpectQuery(regexp.QuoteMeta(querySQL)).
					WithArgs("2024-04-01", "2025-04-01").
					WillReturnRows(summaryRows)
			},
			wantErr: false,
			checkFunc: func(t *testing.T, summaries []*domain.CategoryYearSummary) {
				if len(summaries) != 2 {
					t.Errorf("expected 2 summaries, got %d", len(summaries))
					return
				}

				// カテゴリID順が保証されていないのでマップで検証
				summaryMap := make(map[int]*domain.CategoryYearSummary)
				for _, s := range summaries {
					summaryMap[s.CategoryID] = s
				}

				// 食費（201）のチェック
				if s, ok := summaryMap[201]; ok {
					if s.CategoryName != "食費" {
						t.Errorf("expected category name '食費', got '%s'", s.CategoryName)
					}
					if s.Count != 33 {
						t.Errorf("expected count 33, got %d", s.Count)
					}
					if s.Total != 165000 {
						t.Errorf("expected total 165000, got %d", s.Total)
					}
					if s.Price[0] != 50000 {
						t.Errorf("expected price[0] 50000, got %d", s.Price[0])
					}
					if s.Price[1] != 60000 {
						t.Errorf("expected price[1] 60000, got %d", s.Price[1])
					}
					if s.Price[11] != 55000 {
						t.Errorf("expected price[11] 55000, got %d", s.Price[11])
					}
				} else {
					t.Errorf("category 201 not found in summaries")
				}

				// 月給（100）のチェック
				if s, ok := summaryMap[100]; ok {
					if s.CategoryName != "月給" {
						t.Errorf("expected category name '月給', got '%s'", s.CategoryName)
					}
					if s.Count != 2 {
						t.Errorf("expected count 2, got %d", s.Count)
					}
					if s.Total != 610000 {
						t.Errorf("expected total 610000, got %d", s.Total)
					}
				} else {
					t.Errorf("category 100 not found in summaries")
				}
			},
		},
		{
			name: "正常系: データが空の場合",
			year: 2024,
			mockSetup: func(mock sqlmock.Sqlmock) {
				// カテゴリ取得のモック
				categoryRows := sqlmock.NewRows([]string{"id", "category_id", "name", "category_type"})
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `Category`")).
					WillReturnRows(categoryRows)

				// 空の集計結果
				summaryRows := sqlmock.NewRows([]string{"category_id", "fiscal_month", "total_price", "count"})
				querySQL := `
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
				mock.ExpectQuery(regexp.QuoteMeta(querySQL)).
					WithArgs("2024-04-01", "2025-04-01").
					WillReturnRows(summaryRows)
			},
			wantErr: false,
			checkFunc: func(t *testing.T, summaries []*domain.CategoryYearSummary) {
				if len(summaries) != 0 {
					t.Errorf("expected 0 summaries, got %d", len(summaries))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			repo := NewRecordRepository(db)

			tt.mockSetup(mock)

			summaries, err := repo.GetYearSummary(context.Background(), tt.year)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetYearSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.checkFunc != nil {
				tt.checkFunc(t, summaries)
			}

			// 全ての期待が満たされたか確認
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}
