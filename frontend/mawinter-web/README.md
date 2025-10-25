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

## 環境変数の設定

バックエンド API との接続に必要な環境変数を設定できます。

### ローカル開発環境 (Dev Container)

Dev Container 環境では、`.devcontainer/devcontainer.json` に設定済みです。
デフォルトで `http://localhost:8080` に接続します。

### 本番環境

本番環境では `.env` ファイルまたはシステムの環境変数で設定してください。

`.env.example` をコピーして `.env` ファイルを作成します。

```bash
cp .env.example .env
```

`.env` ファイルの内容:

```bash
# バックエンドAPIのベースURL
NUXT_PUBLIC_MAWINTER_API=https://your-api-server.com
```

### 環境変数の一覧

| 環境変数名                      | 説明                           | デフォルト値           |
| ------------------------------- | ------------------------------ | ---------------------- |
| `NUXT_PUBLIC_MAWINTER_API`      | バックエンド API のベース URL | `http://localhost:8080` |

API のベースエンドポイント (`/api`) は固定で、変更できません。

### API 呼び出しの使い方

`useApi` composable を使用してバックエンド API を呼び出します。

```typescript
// コンポーネント内で使用
const { fetchApi, callApi, getApiUrl } = useApi()

// 例1: useFetch でデータ取得 (SSR対応)
const { data: categories } = await fetchApi('/v3/categories')

// 例2: $fetch で直接呼び出し (クライアントサイド)
const result = await callApi('/v3/record', {
  method: 'POST',
  body: { category_id: 100, price: 1000 }
})

// 例3: URL を取得して別の方法で使用
const url = getApiUrl('/v3/record')
// => http://localhost:8080/api/v3/record
```

## プロジェクト構成

```
mawinter-web/
├── app/              # アプリケーションのメインディレクトリ
│   └── app.vue       # ルートコンポーネント
├── composables/      # 再利用可能な composable 関数
│   └── useApi.ts     # API 呼び出し用 composable
├── public/           # 静的ファイル置き場 (画像など)
├── .env.example      # 環境変数のサンプルファイル
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
