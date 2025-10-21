# internal/domain/

ドメイン層 - ビジネスロジックの中核を配置するディレクトリです。

## 役割

- ビジネスルールとロジックの実装
- エンティティとバリューオブジェクトの定義
- ドメインサービスの実装
- リポジトリインターフェースの定義

## 配置するファイル

- `entity.go` - ドメインエンティティの定義
- `value_object.go` - バリューオブジェクトの定義
- `repository.go` - リポジトリインターフェースの定義
- `service.go` - ドメインサービスの実装

## 重要な原則

- **他の層に依存しない** - このレイヤーは完全に独立している必要があります
- **ビジネスロジックに集中** - 技術的な詳細（DB、HTTP など）は含めません
- **インターフェースで抽象化** - 外部への依存はインターフェースで定義します

## 例

```go
// リポジトリインターフェースの定義例
type RecordRepository interface {
    Save(ctx context.Context, record *Record) error
    FindByID(ctx context.Context, id string) (*Record, error)
}
```
