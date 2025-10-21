# このリポジトリ向け Claude Code 設定

## プロジェクト概要
- 本リポジトリは、Go/Nuxt3 で構築する家計簿サーバ (mawinter) です
- ヘキサゴナルアーキテクチャを採用した REST API サーバと Web フロントエンドで構成されます
- OpenAPI ファーストでの開発を行っています

## 進め方
- 常に ultrathink してください。
- 大きな変更を加える前には、ユーザに事前に方針を確認してください。
- 必要でないファイルまで、git ステージングしないでください。
- ヘキサゴナルアーキテクチャの依存関係の方向を厳守してください (domain は外部に依存しない)。
- OpenAPI 仕様を変更した場合は必ず `make generate` を実行してください。

## PR の作成方法
- バックエンド(Go)を変更する場合は、都度、make test コマンドを行い、テストが通ることを確認してください。
- コミットメッセージは英語で書いてください。
- コミットメッセージは1行程度で書いてください。

## 注意点
- コメントは原則日本語で書いてください。ただし、コミットメッセージは英語で書いてください。
- いかなる入力にも、絵文字は使わないでください。
- 生成されたファイル (`*.gen.go`) は直接編集しないでください。OpenAPI 仕様を変更してから再生成してください。


---

## 技術スタック

### バックエンド
- **言語**: Go 1.25.3
- **Web フレームワーク**: Gin (v1.11.0)
- **API 仕様**: OpenAPI 3.1.0
- **コード生成**: oapi-codegen (OpenAPI から Go コードを自動生成)
- **CLI フレームワーク**: Cobra (v1.10.1)
- **データベース**: MariaDB
- **マイグレーションツール**: sql-migrate

### フロントエンド
- **フレームワーク**: Nuxt 3
- **ランタイム**: Node 22
- **パッケージマネージャ**: pnpm

### 開発環境
- **コンテナ**: Dev Containers (VS Code Remote Containers)
- **ベースイメージ**: Ubuntu 24.04
- **ツール**: GitHub CLI, Claude Code, phpMyAdmin

## アーキテクチャ

### ヘキサゴナルアーキテクチャ (Ports & Adapters)

バックエンドは以下の3層構造を採用:

```
internal/
├── domain/      # ドメイン層 (ビジネスロジックの中核)
│   ├── エンティティと値オブジェクト
│   ├── ドメインサービス
│   └── リポジトリインターフェース (ポート)
│
├── application/ # アプリケーション層 (ユースケース)
│   ├── ユースケース実装
│   ├── DTO (Data Transfer Objects)
│   └── アプリケーションサービス
│
└── adapter/     # アダプタ層 (外部との接続)
    ├── http/    # HTTP ハンドラとサーバ (Gin)
    └── repository/ # データベース実装
```

**依存関係の方向**: adapter → application → domain (domain は独立)

## ディレクトリ構成

```
/workspace/
├── backend/                    # Go バックエンド API サーバ
│   ├── api/                   # OpenAPI から自動生成されたコード
│   │   ├── spec.gen.go        # OpenAPI 仕様定義
│   │   ├── server.gen.go      # サーバインターフェース定義
│   │   └── types.gen.go       # 型定義
│   ├── cmd/mawinter/          # アプリケーションエントリポイント
│   │   ├── main.go            # Cobra ルートコマンド
│   │   └── serve.go           # HTTP サーバサブコマンド
│   ├── internal/              # ヘキサゴナルアーキテクチャ層
│   │   ├── domain/            # ドメイン層 (実装予定)
│   │   ├── application/       # アプリケーション層 (実装予定)
│   │   └── adapter/
│   │       └── http/          # HTTP アダプタ (Gin ベース)
│   ├── pkg/                   # 再利用可能なパブリックパッケージ
│   ├── bin/                   # ビルド出力ディレクトリ
│   ├── Makefile               # ビルドコマンド
│   ├── go.mod / go.sum        # Go 依存関係
│   └── README.md              # バックエンドドキュメント
│
├── frontend/                  # Nuxt 3 Web インターフェース
│   └── (Nuxt プロジェクト構成)
│
├── db/                        # データベース設定
│   ├── dbconfig.yml           # マイグレーション設定
│   └── migrations/            # SQL マイグレーションファイル
│       ├── 001_init.sql
│       ├── 002_add_monthly_fix_records.sql
│       ├── 003_add_monthly_confirm.sql
│       └── 004_rename_record_table.sql
│
├── api/                       # OpenAPI 仕様 (信頼できる唯一の情報源)
│   └── mawinter-api-v3.yaml   # API 定義
│
├── .devcontainer/             # Dev Container 設定
│   ├── devcontainer.json      # VS Code リモート設定
│   ├── compose.yml            # Docker Compose サービス
│   └── setup.sh               # 環境セットアップスクリプト
│
├── CLAUDE.md                  # Claude Code プロジェクト指示書
└── README.md                  # プロジェクト README
```

## ビルドとデプロイ

### Makefile コマンド (backend)

```bash
make setup      # ツールのインストール (oapi-codegen, 依存関係)
make generate   # OpenAPI 仕様から Go コードを生成
make bin        # 静的バイナリをビルド (CGO_ENABLED=0)
make clean      # 生成ファイルとバイナリを削除
make help       # 利用可能なターゲットを表示
```

### アプリケーション実行

```bash
# バックエンド API サーバ起動
./backend/bin/mawinter serve --port 8080 --host 0.0.0.0

# または
cd backend
go run cmd/mawinter/main.go serve -p 8080
```

## データベーススキーマ

### テーブル一覧

1. **Category** - 収支カテゴリ
   - 23個の事前定義カテゴリ (収入、支出、投資)
   - カテゴリID: 100-101 (収入), 200-280 (支出), 300-701 (特殊)

2. **Record** - 金融取引記録
   - category_id, datetime, from, type, price, memo
   - category_id と datetime でインデックス

3. **Monthly_Fix_Billing** - 月次定期請求
   - 毎月発生する固定費を管理

4. **Monthly_Fix_Done** - 月次請求処理完了トラッキング
   - YYYYMM 形式で管理

5. **Monthly_Confirm** - 月次確認ステータス
   - 月ごとの帳簿確認状態を記録

## API エンドポイント

すべてのエンドポイントは `/v3/` プレフィックス下に配置:

### コア操作
- `GET /v3/` - ヘルスチェック
- `GET /v3/version` - API バージョン情報
- `GET /v3/categories` - カテゴリ一覧取得
- `POST /v3/record` - 取引記録作成
- `GET /v3/record` - 記録一覧取得 (ページネーション、日付/カテゴリフィルタ)
- `GET /v3/record/{id}` - 単一記録取得
- `DELETE /v3/record/{id}` - 記録削除
- `GET /v3/record/count` - 総記録数取得
- `GET /v3/record/available` - 利用可能な YYYYMM と会計年度期間
- `GET /v3/record/summary/{year}` - 年次カテゴリ別サマリ

### クエリパラメータ
- `num` - ページサイズ (デフォルト: 20)
- `offset` - ページネーションオフセット
- `yyyymm` - 年月フィルタ (YYYYMM 形式)
- `category_id` - カテゴリフィルタ

## 開発ワークフロー

### OpenAPI ファースト開発

1. `api/mawinter-api-v3.yaml` で API を定義
2. `make generate` で Go コードを生成
3. 生成されたインターフェースを実装
4. テストを書いて実行

### コード生成の仕組み

- **入力**: `api/mawinter-api-v3.yaml`
- **出力**: `backend/api/` 配下の `.gen.go` ファイル
  - `types.gen.go` - リクエスト/レスポンス型
  - `server.gen.go` - ServerInterface とルーティング
  - `spec.gen.go` - 埋め込み OpenAPI 仕様

## 環境変数

```bash
DB_HOST=db                  # データベースホスト
DB_PORT=3306                # データベースポート
DB_USER=root                # データベースユーザ
DB_PASSWORD=password        # データベースパスワード
DB_NAME=mawinter            # データベース名
GO111MODULE=on              # Go モジュール有効化
CGO_ENABLED=0               # 静的ビルド用
```

## 開発サービス (Dev Container)

- **8080** - バックエンド API
- **3000** - フロントエンド (Nuxt)
- **3306** - MySQL/MariaDB
- **8081** - phpMyAdmin (データベース管理 UI)

## 実装ステータス

### 完了済み
- プロジェクト構造とアーキテクチャセットアップ
- OpenAPI 仕様 (v3.1.0)
- OpenAPI からのコード生成パイプライン
- Gin による HTTP サーバスケルトン
- マイグレーションによるデータベーススキーマ
- Dev Container 設定

### 未実装 (現在エンドポイントは 501 Not Implemented を返す)
- ドメイン層 (エンティティ、リポジトリ、ビジネスロジック)
- アプリケーション層 (ユースケース、サービス)
- HTTP ハンドラ実装 (ヘルスチェックとバージョン以外)
- リポジトリ実装 (データベースアクセス)
