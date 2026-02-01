---
name: git-worktree
description: git worktreeを使用したブランチ作業の管理
user_invocable: true
---

# /git-worktree

git worktreeを使用してブランチ作業を管理します。

## 使用方法

```
/git-worktree <command> [args]
```

### コマンド

- `create <type>/<name>` - 新しいworktreeを作成
- `cleanup` - マージ済みworktreeを削除
- `list` - worktree一覧を表示

## 実行手順

### create

1. `git worktree add -b <branch> ../<project>-<short-name> main`
2. 作成したworktreeのパスを報告

### cleanup

1. `git worktree list` で一覧取得
2. 各ブランチの `gh pr list --head <branch> --state merged` を確認
3. マージ済みのworktreeを `git worktree remove` で削除
4. ローカルブランチを `git branch -d` で削除
5. 削除結果を報告

### list

1. `git worktree list` を実行
2. 各worktreeのPR状態を確認して報告
