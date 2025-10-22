package application

import (
	"context"

	"github.com/azuki774/mawinter/internal/domain"
)

// RecordService はレコードに関するアプリケーションサービス
type RecordService struct {
	repo domain.RecordRepository
}

// NewRecordService はRecordServiceを生成する
func NewRecordService(repo domain.RecordRepository) *RecordService {
	return &RecordService{
		repo: repo,
	}
}

// CreateRecord は新しいレコードを作成する
func (s *RecordService) CreateRecord(ctx context.Context, record *domain.Record) (*domain.Record, error) {
	return s.repo.Create(ctx, record)
}

// GetRecordByID は指定されたIDのレコードを取得する
func (s *RecordService) GetRecordByID(ctx context.Context, id int) (*domain.Record, error) {
	return s.repo.FindByID(ctx, id)
}

// GetRecords はレコードを取得する（ページネーション対応）
func (s *RecordService) GetRecords(ctx context.Context, num, offset int, yyyymm string, categoryID int) ([]*domain.Record, error) {
	return s.repo.FindAll(ctx, num, offset, yyyymm, categoryID)
}

// CountRecords は条件に一致するレコードの総数を取得する
func (s *RecordService) CountRecords(ctx context.Context, yyyymm string, categoryID int) (int, error) {
	return s.repo.Count(ctx, yyyymm, categoryID)
}

// DeleteRecord は指定されたIDのレコードを削除する
func (s *RecordService) DeleteRecord(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
