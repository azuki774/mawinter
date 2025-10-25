# mawinter-web

mawinter の Web フロントエンドプロジェクトです。Nuxt 3 を使用して構築されています。

## 必要な環境

- Node.js 22 以降
- pnpm (パッケージマネージャー)

## 初回セットアップ

プロジェクトのディレクトリに移動してから、依存パッケージをインストールします。

```bash
cd frontend/mawinter-web
pnpm install
```

## よく使うコマンド一覧

### 開発サーバーの起動

開発用のサーバーを起動します。`http://localhost:3000` でアクセスできます。
ファイルを編集すると自動的に反映されます (ホットリロード)。

```bash
pnpm dev
```

### 本番用ビルド

本番環境用に最適化されたファイルを生成します。

```bash
pnpm build
```

### 本番ビルドのプレビュー

ビルドした本番用ファイルをローカルで確認します。

```bash
pnpm preview
```

### 静的サイト生成

静的な HTML ファイルとして出力します (SSG: Static Site Generation)。

```bash
pnpm generate
```

## プロジェクト構成

```
mawinter-web/
├── app/              # アプリケーションのメインディレクトリ
│   └── app.vue       # ルートコンポーネント
├── public/           # 静的ファイル置き場 (画像など)
├── nuxt.config.ts    # Nuxt の設定ファイル
├── package.json      # npm パッケージの設定
└── tsconfig.json     # TypeScript の設定
```

## 参考ドキュメント

- [Nuxt 3 公式ドキュメント](https://nuxt.com/docs/getting-started/introduction)
- [Vue 3 公式ドキュメント](https://ja.vuejs.org/guide/introduction.html)

## トラブルシューティング

### ポートが既に使用されている

開発サーバーのデフォルトポート (3000) が使用中の場合は、別のポートを指定できます。

```bash
pnpm dev --port 3001
```

### キャッシュのクリア

動作がおかしい場合は、`.nuxt` ディレクトリと `node_modules` を削除して再インストールしてください。

```bash
rm -rf .nuxt node_modules
pnpm install
```
