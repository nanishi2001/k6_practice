# k6 Performance Testing Project

k6とNix、TypeScriptを使用した非機能テストプロジェクト

## Tech Stack

- **k6**: 負荷テストツール (Grafana Labs)
- **TypeScript**: 型安全なテストスクリプト
- **Bun**: TypeScript → JavaScript ビルド & パッケージマネージャ
- **Nix Flakes**: 再現可能な開発環境

## Project Structure

```
.
├── flake.nix           # Nix開発環境定義
├── package.json        # 依存関係
├── tsconfig.json       # TypeScript設定
├── src/                # TypeScriptソース
│   ├── load-test.ts        # 負荷テスト
│   ├── stress-test.ts      # ストレステスト
│   └── spike-test.ts       # スパイクテスト
├── dist/               # ビルド出力 (gitignore)
└── results/            # テスト結果出力 (gitignore)
```

## Development Environment

```bash
# 開発環境に入る
nix develop

# 依存関係インストール
bun install

# ビルド
bun run build

# テスト実行
bun run test:load
bun run test:stress
bun run test:spike

# 全テスト実行
bun run test:all
```

## Critical Rules

### 1. TypeScript

- 型定義を必ず使用（`@types/k6`）
- `Options`型でk6オプションを定義
- `any`型の使用禁止
- 関数の戻り値型を明示

### 2. テストスクリプト作成

- 各テストは単一の目的を持つ
- オプションは`options`オブジェクトで明示的に定義
- しきい値(thresholds)を必ず設定
- チェック(checks)で成功条件を定義

### 3. パフォーマンス基準

- 応答時間 p(95) < 500ms
- エラー率 < 1%
- 可用性 > 99%

### 4. テストパターン

| パターン | VUs | 期間 | 用途 |
|---------|-----|------|------|
| 負荷テスト | 10-50 | 5-10分 | 通常負荷の性能測定 |
| ストレステスト | 50-200 | 10-20分 | 限界性能の特定 |
| スパイクテスト | 1-100-1 | 5分 | 急激な負荷変動への耐性 |

### 5. コードスタイル

- ES2020+構文を使用
- 関数は小さく保つ
- マジックナンバーを避け、定数を使用
- コメントは日本語で記述可

## k6 Options Template

```typescript
import { Options } from 'k6/options';

export const options: Options = {
  stages: [
    { duration: '1m', target: 10 },
    { duration: '3m', target: 10 },
    { duration: '1m', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],
    http_req_failed: ['rate<0.01'],
  },
};
```

## Git Flow

### ブランチ命名規則

| プレフィックス | 用途 | 例 |
|--------------|------|-----|
| `feat/` | 新機能追加 | `feat/api-unit-tests` |
| `fix/` | バグ修正 | `fix/auth-token-expiry` |
| `refactor/` | リファクタリング | `refactor/extract-test-credentials` |
| `docs/` | ドキュメント | `docs/update-readme` |
| `ci/` | CI/CD関連 | `ci/add-github-actions` |
| `test/` | テスト追加 | `test/add-integration-tests` |

### コミットメッセージ規則

```
<type>(<scope>): <subject>

<body>

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>
```

**type**: `feat`, `fix`, `refactor`, `docs`, `ci`, `test`, `chore`

### PR作成フロー（git worktree使用）

```bash
# 1. worktreeで作業ブランチを作成
git worktree add -b <type>/<description> ../<project>-<short-name> main

# 2. worktreeディレクトリで作業・コミット
cd ../<project>-<short-name>
# ... 作業 ...
git add <files>
git commit -m "<type>(<scope>): <subject>"

# 3. リモートにプッシュ
git push -u origin <type>/<description>

# 4. PR作成
gh pr create

# 5. マージ後、worktreeとブランチを削除
git worktree remove ../<project>-<short-name>
git branch -d <type>/<description>
```

### 注意事項

- `main` ブランチへの直接プッシュ・コミットは禁止
- PRは必ずCIを通過させる
- 1つのPRは1つの目的に絞る
- マージ済みのworktreeは速やかに削除する

## Available Commands

- `nix develop` - 開発環境に入る
- `bun install` - 依存関係インストール
- `bun run build` - TypeScriptビルド
- `bun run test:load` - 負荷テスト実行
- `bun run test:stress` - ストレステスト実行
- `bun run test:spike` - スパイクテスト実行
- `bun run test:all` - 全テスト実行
- `go test ./... -v` - Go APIテスト実行 (api/ディレクトリ内)
