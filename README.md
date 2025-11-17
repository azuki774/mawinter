# mawinter

Go/Nuxt 製の家計簿サーバ

## frontend

- Nuxt3 製の Web 画面
- https://github.com/azuki774/mawinter-front をベースに作成

## backend

- Go 製のバックエンド API サーバ
- https://github.com/azuki774/mawinter-server をベースに作成

## トレース出力

- Go バックエンド: 環境変数 `OTLP_SERVER` に `host:port` を指定すると `http(s)://<host:port>/v1/traces` へトレースを送信します。
- 変数が未設定の場合はトレース送信を自動的に無効化します。
