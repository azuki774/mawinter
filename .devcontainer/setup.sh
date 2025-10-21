#!/bin/bash
# .devcontainer/setup.sh

set -e

echo "Setting up mawinter development environment..."
mkdir -p /home/vscode/go/{bin,pkg,src}

# MySQL クライアント
if ! command -v mysql &> /dev/null; then
    echo "Installing MySQL client..."
    sudo apt-get update
    sudo apt-get install -y mysql-client
    echo "MySQL client installed"
fi

# pnpm
if ! command -v pnpm &> /dev/null; then
    echo "Installing pnpm..."
    npm install -g pnpm
    echo "pnpm installed"
fi

# Claude Code
if ! command -v claude-code &> /dev/null; then
    echo "Installing Claude Code..."
    npm install -g @anthropic-ai/claude-code
    echo "Claude Code installed"
fi

# Go ツール
echo "Installing Go tools..."
go install github.com/rubenv/sql-migrate/...@latest
echo "Go tools installed"

# プロジェクトのセットアップ
if [ -f Makefile ]; then
    echo "Running make setup..."
    make setup
fi

# GitHub CLI 認証
if [ -n "$GITHUB_TOKEN" ]; then
    echo "Authenticating GitHub CLI..."
    echo "$GITHUB_TOKEN" | gh auth login --with-token 2>/dev/null || true
    gh auth status
    echo "GitHub CLI authenticated"
else
    echo "WARNING: GITHUB_TOKEN not set. Run 'gh auth login' manually if needed"
fi

# Claude Code の確認
if [ -n "$ANTHROPIC_API_KEY" ]; then
    echo "ANTHROPIC_API_KEY detected"
else
    echo "WARNING: ANTHROPIC_API_KEY not set. Set it to use Claude Code"
fi

# バージョン確認
echo ""
echo "Installed versions:"
echo "  Go:      $(go version | awk '{print $3}')"
echo "  Node:    $(node --version)"
echo "  pnpm:    $(pnpm --version)"
echo "  gh:      $(gh --version | head -n1)"
echo "  MySQL:   $(mysql --version | awk '{print $3}')"

echo ""
echo "mawinter development environment ready!"
echo ""
echo "Next steps:"
echo "  1. Run 'make dev' to start development servers"
echo "  2. Run 'claude-code \"your task\"' to use Claude Code"
echo "  3. Visit http://localhost:3000 for frontend"
echo "  4. Visit http://localhost:8080 for backend API"
