---
name: GitHub 操作スキル
description: gitコミット作成、gitブランチの作成や、GitHub への push や GitHub の Pull Request 作成に使用します。
---

# GitHub 操作スキル

## いつ使うか

- git ブランチを作成するとき
- git commit でコミットを作成するとき
- Pull Request を作成するとき

## 指示

- git ブランチの作成時は、ブランチ名を下記の通りにすること。
    - `<種別>/<説明>`
    - 種別は、`feat`, `chore`, `docs`, `fix` のうち適切なものにすること
    - 説明は、単語を `-` 区切りにした、開発内容を簡単に説明したものにすること
        - 例: `feat/update-dockerfile-baseimage`
- コミットは下記のように作成すること
    - 不要なファイルを git add しないこと
    - `git diff` などを行い、差分を適切に理解すること
    - コミットメッセージは、英語で1文だけのシンプルなものにすること
- Pull Request は下記のように作成すること
    - PRのタイトルは英語で作成すること
    - PRの本文は日本語で入力すること
