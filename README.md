# mawinter

Go/Nuxt 製の家計簿サーバ

## ローカル開発 (API 開発用)

```bash
# DB 起動
docker compose -f deployment/compose-for-apidev.yml up -d db

# マイグレーション
make -C backend migrate-up

# ダミーデータ投入
make -C backend seed

# API 起動 (ホスト)
go run ./backend/cmd/mawinter serve --port 8080
```

詳細は `deployment/README.md` と `AGENTS.md` を参照

## frontend

- Nuxt3 製の Web 画面
- https://github.com/azuki774/mawinter-front をベースに作成

## backend

- Go 製のバックエンド API サーバ
- https://github.com/azuki774/mawinter-server をベースに作成

## トレース出力

- Go バックエンド: 環境変数 `OTLP_SERVER` に `host:port` を指定すると `http(s)://<host:port>/v1/traces` へトレースを送信します。
- 変数が未設定の場合はトレース送信を自動的に無効化します。

## Renovate の更新方針

- 設定は [`renovate.jsonc`](renovate.jsonc)。通常の更新 PR は毎週金曜日の 20:00〜20:59（`Asia/Tokyo`）に作成します。20:00 ちょうどの実行保証ではなく、Renovate の実行タイミングに依存する時間枠です。
- 通常の更新 PR は同時に最大 7 件で、時間あたりの作成数制限はありません。
- 脆弱性アラートに基づくセキュリティ更新は、通常のスケジュールと同時 7 件の制限を迂回して即時 PR 作成の対象になります。GitHub の dependency graph と Dependabot alerts の有効化、および Renovate がアラートを読み取り PR を作成できる権限が必要です。
- GitHub Actions はコミット SHA（SHA-1）に固定し、Actions の更新だけを自動マージ対象にします。他の依存関係は自動マージしません。
