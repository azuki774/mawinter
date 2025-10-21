.PHONY: setup generate clean bin test help

# 変数定義
OPENAPI_SPEC := ../api/mawinter-api-v3.yaml
GENERATED_DIR := ./api
MODULE_NAME := github.com/azuki774/mawinter
BIN_DIR := ./bin
BINARY_NAME := mawinter
CMD_DIR := ./cmd/mawinter

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

clean: ## 生成されたファイルとバイナリを削除
	rm -rf $(GENERATED_DIR)
	rm -rf $(BIN_DIR)
