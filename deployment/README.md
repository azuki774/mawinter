# deployment

API 開発用のローカル環境

## 構成

- `compose-for-apidev.yml` - DB (MariaDB) と周辺サービスの compose 定義
- `../db/migrations/` - マイグレーションファイル (sql-migrate)
- `../db/seed/dummy.sql` - ダミーデータ

## 前提

- Docker / Podman + compose (docker compose / podman compose) がインストール済み
- Go 1.25 以上、sql-migrate (`go install github.com/rubenv/sql-migrate/...@latest`)
- MySQL クライアント (`mysql` コマンド、seed 投入時に使用)
- パスワード等はテスト用のベタ書き (`root` / `password` / `mawinter`)

## 使い方

### 1. DB 起動

```bash
docker compose -f deployment/compose-for-apidev.yml up -d db
# podman の場合
podman compose -f deployment/compose-for-apidev.yml up -d db
# または Makefile 経由 (docker/podman 自動判定)
make -C backend db-up
```

### 2. マイグレーション

```bash
# ホストから実行 (推奨)
make -C backend migrate-up

# 状態確認
make -C backend migrate-status

# docker 環境から実行する場合
sql-migrate up -config=db/dbconfig.yml -env=docker
```

### 3. ダミーデータ投入

```bash
make -C backend seed
```

冪等性: `memo LIKE 'seed-%'` のデータのみを削除して再投入する

### 4. API 起動

ホストで実行する場合:

```bash
go run ./backend/cmd/mawinter serve --port 8080
# または
make -C backend bin && ./backend/bin/mawinter serve --port 8080
```

コンテナで実行する場合:

```bash
docker compose -f deployment/compose-for-apidev.yml --profile api up -d --build
# podman の場合
podman compose -f deployment/compose-for-apidev.yml --profile api up -d --build
```

### 5. 周辺ツール (任意)

```bash
# phpMyAdmin
docker compose -f deployment/compose-for-apidev.yml --profile tools up -d
# podman の場合
podman compose -f deployment/compose-for-apidev.yml --profile tools up -d
# http://localhost:8081 でアクセス (user: root, password: password)
```

### 6. クリーンアップ

```bash
# コンテナ停止
docker compose -f deployment/compose-for-apidev.yml down
# podman の場合
podman compose -f deployment/compose-for-apidev.yml down

# ボリューム含めて削除 (DB データ全削除)
docker compose -f deployment/compose-for-apidev.yml down -v

# ワンコマンドでリセット (DB再作成 + マイグレーション + シード)
make -C backend db-reset
# podman 環境でも `make -C backend db-reset` は自動で podman compose を使用
```

## 環境変数

`backend/pkg/config` は以下を参照する:

```
DB_HOST=127.0.0.1  # ホストから接続時は 127.0.0.1, コンテナ間は db
DB_PORT=3306
DB_USER=root
DB_PASS=password
DB_NAME=mawinter
```

compose 内の `api` サービスは `DB_HOST=db` で接続する
