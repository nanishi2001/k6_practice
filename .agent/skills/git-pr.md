---
name: git-pr
description: PRの作成と管理
user_invocable: true
---

# /git-pr

PRを作成・管理します。

## 使用方法

```
/git-pr [command]
```

### コマンド

- `create` - 現在のブランチでPRを作成（デフォルト）
- `status` - 現在のブランチのPR状態を確認

## 実行手順

### create

1. `git status` で変更を確認
2. `git diff` で差分を確認
3. `git log` でコミット履歴を確認
4. リモートにプッシュ: `git push -u origin <branch>`
5. `gh pr create` でPR作成
6. PR URLを報告

### status

1. `gh pr view` でPR情報を取得
2. CI状態とレビュー状態を報告
