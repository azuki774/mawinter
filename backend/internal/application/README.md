# internal/application/

アプリケーション層 - ユースケースを配置するディレクトリです。

## 役割

- アプリケーションのユースケース実装
- ドメインオブジェクトの調整
- トランザクション管理
- 入出力の変換（DTO の定義）

## 配置するファイル

- `usecase.go` - ユースケースの実装
- `dto.go` - Data Transfer Object（入出力用のデータ構造）
- `service.go` - アプリケーションサービスの実装

## 依存関係

- **domain** - ドメイン層に依存します
- ドメインのリポジトリインターフェースを使用します

## 重要な原則

- **ユースケース単位で実装** - 1つの機能を1つのユースケースとして実装
- **技術詳細から独立** - HTTP や DB などの実装詳細は含めません
- **ドメインロジックは呼び出すのみ** - ビジネスロジックはドメイン層に委譲

## 例

```go
// ユースケースの実装例
type RecordUseCase struct {
    repo domain.RecordRepository
}

func (u *RecordUseCase) CreateRecord(ctx context.Context, input CreateRecordInput) (*RecordDTO, error) {
    // ドメインオブジェクトの生成
    // リポジトリの呼び出し
    // DTO への変換
}
```
