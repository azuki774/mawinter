# seed

API 開発用のダミーデータ

## ファイル

- `dummy.sql` - `Record` 200件 + `Monthly_Fix_Billing` 5件 + `Monthly_Fix_Done` / `Monthly_Confirm` 各18ヶ月分

## 投入方法

```bash
make -C backend seed
```

内部では以下を実行している:

```bash
mysql -h 127.0.0.1 -P 3306 -u root -ppassword mawinter < db/seed/dummy.sql
```

## 仕様

- `Record`: 200件、カテゴリ23種をランダム分散、日付は 2024-04-01 〜 2025-09-30
  - 収入カテゴリ (100,101,110) は 50,000〜350,000 円
  - 家賃 (200) は 50,000〜90,000 円
  - その他は 100〜30,000 円
  - `from` は `seed-*`、`type` は `cash` / `credit` / `transfer` / 空文字
  - `memo` は `seed-XXX-カテゴリID` 形式で冪等性担保
- `Monthly_Fix_Billing`: 家賃・通信費等の固定費5件
- `Monthly_Fix_Done`: 202404〜202509 (202509 のみ未完了)
- `Monthly_Confirm`: 202404〜202509 (202508 以降は未確認)

## 冪等性

先頭で以下を実行して seed 由来データのみを削除:

```sql
DELETE FROM `Record` WHERE `memo` LIKE 'seed-%';
DELETE FROM `Monthly_Fix_Billing` WHERE `memo` LIKE 'seed-%';
```

`Category` はマイグレーションで投入済みのため触らない
