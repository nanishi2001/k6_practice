---
trigger: always_on
---

## Fundamental Rules

- **Package Manager**: プロジェクトのパッケージマネージャは **`bun`** です。npm/yarn/pnpmは使用しないでください。
- **Response Language**: 回答は常に**日本語**で行ってください。
- **Development Environment**: `nix develop` で開発環境に入ってから作業してください。

## Artifact Strategy

- 作業（コーディング、ドキュメント作成、プランニング等）を開始する前に、必ず **Artifacts (task.md, implementation_plan.md)** を作成・提示してください。
- チャット本文のみで完結させず、成果物を明確に切り分けてください。

## Hallucination Prevention

- 正確な情報に基づき回答してください。不明な場合は「わかりません」と伝えてください。
- 推測で回答する場合は、その旨を明示してください。

## Security Rules

### コミット前チェック（必須）

- ハードコードされた秘密情報がないこと（APIキー、パスワード、トークン）
- すべてのユーザー入力が検証されていること
- SQLインジェクション対策（パラメータ化クエリ）
- XSS対策（HTMLのサニタイズ）
- エラーメッセージで機密情報を漏洩しないこと

### 秘密情報の管理

```typescript
// NG: ハードコード
const apiKey = 'sk-xxxx';

// OK: 環境変数
const apiKey = process.env.API_KEY;
if (!apiKey) {
  throw new Error('API_KEY not configured');
}
```

## Coding Style

### 基本原則

- **不変性（Immutability）**: 既存オブジェクトを変更せず、新しいインスタンスを作成する
- **ファイル構成**: 1ファイル200-400行を目安、最大800行まで
- **関数サイズ**: 50行以内に収める
- **ネスト深度**: 最大4レベルまで

### TypeScript固有

- `any`型の使用禁止
- 関数の戻り値型を明示
- `@types/k6`を必ず使用
- ES2020+構文を使用

### k6テストスクリプト

- 各テストは単一の目的を持つ
- `options`オブジェクトで設定を明示的に定義
- `thresholds`（しきい値）を必ず設定
- `check()`で成功条件を定義
- マジックナンバーを避け、定数を使用

## Testing Rules

### パフォーマンス基準

| メトリクス | 基準値 |
|-----------|--------|
| 応答時間 p(95) | < 500ms |
| エラー率 | < 1% |
| 可用性 | > 99% |

### テストパターン

| パターン | VUs | 期間 | 用途 |
|---------|-----|------|------|
| 負荷テスト | 10-50 | 5-10分 | 通常負荷の性能測定 |
| ストレステスト | 50-200 | 10-20分 | 限界性能の特定 |
| スパイクテスト | 1→100→1 | 5分 | 急激な負荷変動への耐性 |

### テスト失敗時の対応

1. エラーメッセージを分析
2. しきい値が適切か確認
3. 対象システムの状態を確認
4. 段階的に修正・検証

## Agent Orchestration

### 利用可能なエージェント

| エージェント | 用途 | トリガー |
|-------------|------|----------|
| `planner` | 実装計画の作成 | 複雑な機能追加時 |
| `k6-reviewer` | k6テストのレビュー | テストスクリプト作成後 |

### エージェント呼び出しルール

- 複雑なタスクには`planner`を使用してから実装
- 新しいk6テストを書いたら`k6-reviewer`でレビュー
- 独立したタスクは並列実行で効率化

## Performance (Model Selection)

### モデル選択ガイドライン

| モデル | 用途 |
|--------|------|
| Haiku | 軽量タスク、頻繁な呼び出し |
| Sonnet | 通常のコーディング作業 |
| Opus | 深い推論、アーキテクチャ決定 |

### コンテキスト管理

- 大規模リファクタリングではコンテキストの最後1/5を避ける
- 複数ファイルにまたがる変更は段階的に実施

## Git Workflow

### ブランチ運用（必須）

**mainブランチへの直接コミットは禁止です。** 以下のフローに従ってください：

1. **ブランチを作成**: `git checkout -b <type>/<description>`
2. **作業・コミット**: 変更をコミット
3. **PRを作成**: `gh pr create` でPull Requestを作成
4. **マージ**: レビュー後にマージ

### ブランチ命名規則

```
<type>/<short-description>
```

例：
- `feat/add-soak-test`
- `fix/threshold-calculation`
- `refactor/handler-structure`

### コミットメッセージ形式

```
<type>: <description>

[optional body]
```

### タイプ

- `feat`: 新機能
- `fix`: バグ修正
- `test`: テスト追加・修正
- `docs`: ドキュメント
- `refactor`: リファクタリング
- `chore`: その他

### PR作成時の注意

- タイトルは70文字以内
- 変更内容のサマリーを記載
- テスト計画を含める

## Browser Automation

- ブラウザ操作には必ず **Playwright MCP tools** を使用してください。
- Node.jsやPythonによる自前スクリプトの実行は禁止します。