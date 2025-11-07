# Repository Guidelines

## Project Structure & Module Organization
The monorepo hosts the Go API in `backend/` (entrypoint `cmd/mawinter`), generated interfaces under `backend/api/`, domain logic under `backend/internal/`, and build artifacts in `backend/bin/`. The Nuxt client lives in `frontend/mawinter-web`, with runtime assets in `public/` and reusable UI in `components/`. Database migrations reside in `db/migrations` with connection settings in `db/dbconfig.yml`. The shared OpenAPI contract is `api/mawinter-api-v3.yaml`; regenerate backend stubs whenever the spec changes.

## Build, Test, and Development Commands
Backend tasks (run inside `backend/`):
- `make setup` installs oapi-codegen and syncs Go modules.
- `make generate` refreshes Go types and Gin handlers from the OpenAPI spec.
- `make test` runs all Go unit tests (`go test ./...`).
- `make bin` produces a static binary at `backend/bin/mawinter`.
Frontend tasks (run inside `frontend/mawinter-web`):
- `pnpm dev` starts the Nuxt dev server with hot reload.
- `pnpm build` creates the production bundle.
- `pnpm lint` executes ESLint against the Vue codebase.

## Coding Style & Naming Conventions
Go code must stay `gofmt`-clean; keep packages lowercase and give handlers explicit names such as `GetMonthlyRecords`. Generated files belong in `backend/api` and should not be edited manually. Vue components follow PascalCase filenames in `components/`, with shared logic in `composables/`. Use two-space indentation in Vue/TypeScript files and honor the repo ESLint configuration (`eslint.config.mjs`).

## Testing Guidelines
Backend tests live beside their packages; add suites under `backend/internal/<pkg>/` using the `_test.go` suffix. Cover new handlers with happy-path and validation cases and run `make test` before pushing. The frontend currently lacks automated tests—run `pnpm lint` and manually exercise new flows through `pnpm dev`. When introducing Vitest or end-to-end checks, keep specs near the feature directory for easier review.

## Commit & Pull Request Guidelines
Commit messages follow short imperative phrasing (`Add graph page`, `Fix category display`). Reference issues with `Refs #<id>` when relevant and keep each commit focused. Pull requests should include a concise summary, testing notes, and UI screenshots for visible changes. Call out OpenAPI updates explicitly and confirm regenerated files under `backend/api`.


---
## 注意点
- 大きな変更を加える前には、ユーザに事前に方針を確認してください。
- 必要でないファイルまで、git ステージングしないでください。
- ヘキサゴナルアーキテクチャの依存関係の方向を厳守してください (domain は外部に依存しない)。
- OpenAPI 仕様を変更した場合は必ず `make generate` を実行してください。
- バックエンド(Go)を変更する場合は、都度、`make test` コマンドを行い、テストが通ることを確認してください。
- **フロントエンド(Nuxt)を変更する場合は、必ずpush前に `cd frontend/mawinter-web && pnpm run lint` を実行し、ESLintエラーがないことを確認してください。**
- PR の本文は日本語で書いてください。
- コミットメッセージは英語で書いてください。
- **コミットメッセージは1行程度で書いてください。**
- GitHub への変更は GitHub CLI を利用してください。`gh` コマンドは認証済にしています。
