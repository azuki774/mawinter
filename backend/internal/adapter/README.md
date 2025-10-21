# internal/adapter/

アダプター層 - 外部との接続を担当するディレクトリです。

## 役割

- 外部システムとの接続実装
- ドメインインターフェースの具体的な実装
- プロトコル変換（HTTP、データベースなど）

## ディレクトリ構成

- [http/](http/) - HTTP サーバとハンドラーの実装
- [repository/](repository/) - データベースリポジトリの実装

## 配置するファイルの種類

### http/
- `server.go` - HTTP サーバの設定
- `handler.go` - HTTP ハンドラーの実装
- `middleware.go` - ミドルウェアの実装

### repository/
- `record_repository.go` - レコードリポジトリの実装
- `category_repository.go` - カテゴリリポジトリの実装

## 依存関係

- **domain** - ドメインインターフェースを実装
- **application** - ユースケースを呼び出す

## 重要な原則

- **ポートの実装** - ドメイン層で定義されたインターフェース（ポート）を実装
- **技術詳細はここに** - HTTP、DB、外部 API などの具体的な実装
- **変換責務** - 外部の形式とドメインモデルの相互変換

## 例

```go
// リポジトリの実装例
type recordRepository struct {
    db *sql.DB
}

func (r *recordRepository) Save(ctx context.Context, record *domain.Record) error {
    // DB への保存実装
}
```
