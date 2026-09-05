.PHONY: setup generate clean bin test build help migrate-up migrate-down migrate-status seed db-up db-down db-reset db-logs

# 変数定義
OPENAPI_SPEC := ../api/mawinter-api-v3.yaml
GENERATED_DIR := ./api
MODULE_NAME := github.com/azuki774/mawinter
BIN_DIR := ./bin
BINARY_NAME := mawinter
CMD_DIR := ./cmd/mawinter
IMAGE_NAME := ghcr.io/azuki774/mawinter-api
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo 'dev')
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')
COMPOSE_FILE := ../deployment/compose-for-apidev.yml
DBCONFIG := ../db/dbconfig.yml
MIGRATE_ENV ?= local
DB_HOST ?= 127.0.0.1
DB_PORT ?= 3306
DB_USER ?= root
DB_PASS ?= password
DB_NAME ?= mawinter
# docker / podman どちらでも動作するように自動判定 (環境変数 COMPOSE で上書き可能)
COMPOSE ?= $(shell if command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then echo "docker compose"; elif command -v podman >/dev/null 2>&1; then echo "podman compose"; else echo "docker compose"; fi)
BUILDER ?= $(shell if command -v podman >/dev/null 2>&1 && command -v docker >/dev/null 2>&1; then echo "docker"; elif command -v podman >/dev/null 2>&1; then echo "podman"; else echo "docker"; fi)

help: ## ヘルプメッセージを表示
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## 開発環境のセットアップ（ツールと依存関係のインストール）
	@echo "Setting up development environment..."
	@echo "Installing oapi-codegen..."
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	@echo ""
	@echo "Initializing Go module if not exists..."
	@if [ ! -f go.mod ]; then \
		go mod init $(MODULE_NAME); \
	fi
	@echo ""
	@echo "Installing dependencies..."
	go get -u github.com/gin-gonic/gin
	go get -u github.com/oapi-codegen/runtime
	go mod tidy
	@echo ""
	@echo "Setup complete!"

generate: ## OpenAPI仕様からGoコードを生成
	@echo "Generating Go code from OpenAPI spec..."
	@mkdir -p $(GENERATED_DIR)
	@echo "Generating types..."
	oapi-codegen -package api -generate types -o $(GENERATED_DIR)/types.gen.go $(OPENAPI_SPEC)
	@echo "Generating Gin server..."
	oapi-codegen -package api -generate gin -o $(GENERATED_DIR)/server.gen.go $(OPENAPI_SPEC)
	@echo "Generating spec embedding..."
	oapi-codegen -package api -generate spec -o $(GENERATED_DIR)/spec.gen.go $(OPENAPI_SPEC)
	go mod tidy
	@echo ""
	@echo "Code generation complete! Generated files are in $(GENERATED_DIR)/"
	@echo ""

bin: ## 静的バイナリをビルド
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 go build -a -tags netgo \
		-ldflags "-extldflags '-static' \
		-X main.version=$$(git describe --tag --abbrev=0 2>/dev/null || echo 'dev') \
		-X main.revision=$$(git rev-list -1 HEAD) \
		-X main.build=$$(git describe --tags 2>/dev/null || echo 'dev')" \
		-o $(BIN_DIR)/$(BINARY_NAME) $(CMD_DIR)

test: ## テストを実行
	go test -v ./...

build: ## Dockerイメージをビルド
	$(BUILDER) build -t $(IMAGE_NAME):$(VERSION) -t $(IMAGE_NAME):$(COMMIT) -t $(IMAGE_NAME):latest .

clean: ## 生成されたファイルとバイナリを削除
	rm -rf $(GENERATED_DIR)
	rm -rf $(BIN_DIR)

# マイグレーション関連
migrate-up: ## マイグレーションを実行
	sql-migrate up -config=$(DBCONFIG) -env=$(MIGRATE_ENV)

migrate-down: ## マイグレーションを差し戻し
	sql-migrate down -config=$(DBCONFIG) -env=$(MIGRATE_ENV)

migrate-status: ## マイグレーション状態を確認
	sql-migrate status -config=$(DBCONFIG) -env=$(MIGRATE_ENV)

seed: ## ダミーデータを投入
	@if command -v mysql >/dev/null 2>&1; then \
		mysql -h $(DB_HOST) -P $(DB_PORT) -u $(DB_USER) -p$(DB_PASS) $(DB_NAME) < ../db/seed/dummy.sql; \
	elif command -v podman >/dev/null 2>&1 && podman ps --format "{{.Names}}" | grep -q "mawinter-apidev"; then \
		podman exec -i mawinter-apidev_db_1 mysql -u $(DB_USER) -p$(DB_PASS) $(DB_NAME) < ../db/seed/dummy.sql; \
	elif command -v docker >/dev/null 2>&1 && docker ps --format "{{.Names}}" | grep -q "mawinter-apidev"; then \
		docker exec -i mawinter-apidev_db_1 mysql -u $(DB_USER) -p$(DB_PASS) $(DB_NAME) < ../db/seed/dummy.sql; \
	else \
		echo "mysql client not found and no running DB container. Install mysql-client or run 'make db-up' first."; exit 1; \
	fi

# ローカルDB操作 (docker / podman どちらでも可、COMPOSE 変数で上書き可能)
db-up: ## DBコンテナを起動
	$(COMPOSE) -f $(COMPOSE_FILE) up -d db

db-down: ## DBコンテナを停止
	$(COMPOSE) -f $(COMPOSE_FILE) down

db-logs: ## DBログを表示
	$(COMPOSE) -f $(COMPOSE_FILE) logs -f db

db-reset: ## DBをリセットしてマイグレーションとシードを実行
	$(COMPOSE) -f $(COMPOSE_FILE) down -v
	$(COMPOSE) -f $(COMPOSE_FILE) up -d db
	@echo "Waiting for DB to be healthy..."
	@for i in 1 2 3 4 5 6 7 8 9 10; do \
		if $(COMPOSE) -f $(COMPOSE_FILE) ps | grep -q "healthy"; then break; fi; \
		echo "  waiting... ($$i/10)"; sleep 2; \
	done
	$(MAKE) migrate-up MIGRATE_ENV=local
	$(MAKE) seed
