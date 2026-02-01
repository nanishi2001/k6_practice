---
name: k6-run
description: k6テストを実行するクイックコマンド
user_invocable: true
---

# /k6-run

k6テストを素早く実行します。

## 使用方法

```
/k6-run [type]
```

- `load` - 負荷テスト
- `stress` - ストレステスト
- `spike` - スパイクテスト
- `all` - 全テスト

## 実行手順

1. `bun run build` でビルド
2. `bun run test:<type>` で指定テストを実行
3. 結果を要約して報告
