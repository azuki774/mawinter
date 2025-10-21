# Backend

Mawinter の Go 製バックエンド API サーバです。

## アーキテクチャ

ヘキサゴナルアーキテクチャ（ポート&アダプターアーキテクチャ）を採用しています。

## ディレクトリ構成

- [api/](api/) - OpenAPI から自動生成されたコード
- [bin/](bin/) - ビルド成果物の出力先
- [cmd/](cmd/) - アプリケーションのエントリーポイント
- [internal/](internal/) - 内部実装（ヘキサゴナルアーキテクチャの各層）
- [pkg/](pkg/) - 外部に公開可能な共有パッケージ

## ビルド方法

```bash
# 開発ビルド
make build

# 静的バイナリビルド
make bin
```

## 実行方法

```bash
# サーバ起動
./bin/mawinter serve
```
