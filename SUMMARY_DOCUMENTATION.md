# djou - ドキュメントまとめ

## プロジェクト概要

**djou** (Dev Journal CLI) は、開発者の日々の作業を記録・管理・分析するためのCLIツール。
Go言語で実装され、SQLiteをデータストアとして使用する。現在のバージョンは **v0.1.3**。

---

## 1. README.md

ユーザー向けのメインドキュメント。以下を網羅している。

### 主要機能一覧

| コマンド | 概要 |
|---------|------|
| `djou` | インタラクティブフォームで作業を記録 |
| `djou list` | ログ一覧表示（`--week`, `--month`, `--limit`, `--tag`, `--no-interactive`） |
| `djou search` | キーワード検索（`--limit`, `--tag`, `--detail`） |
| `djou stats` | 統計表示（総件数、見積もり精度、月別内訳） |
| `djou export` | CSV出力（`--output`, `--from`, `--to`） |
| `djou config` | 設定管理（`init`, `show`, `edit`） |
| `djou version` | バージョン情報 |
| `djou mcp` | MCP (Model Context Protocol) サーバー起動 |

### データベーススキーマ (logs テーブル)

| カラム | 型 | 説明 |
|--------|-----|------|
| id | INTEGER PRIMARY KEY | 自動採番 |
| created_at | DATETIME | 自動設定 |
| task_name | TEXT NOT NULL | タスク名 |
| estimate_hours | REAL NOT NULL | 見積もり時間 |
| actual_hours | REAL NOT NULL | 実績時間 |
| memo | TEXT | メモ（任意） |
| tags | TEXT | カンマ区切りタグ（任意） |

### インストール方法

- `go install github.com/kca2250/djou/cmd/djou@latest`
- ソースビルド: `git clone && go build -o djou ./cmd/djou`

---

## 2. docs/djou-spec.md - 全体仕様書

v2スキーマの設計仕様。現在の実装と一部異なる拡張フィールド（`ai_minutes`, `problem`, `solution`, `learning`）を含む将来仕様を定義。

### 仕様上のスキーマ（未実装フィールドあり）

- task_name, estimate_hours, actual_hours（実装済み）
- ai_minutes, problem, solution, learning（未実装 - 将来拡張）

### 技術スタック

- CLI: spf13/cobra
- UI: charmbracelet/huh
- DB: SQLite (modernc.org/sqlite)
- MCP: mark3labs/mcp-go

---

## 3. docs/common-spec.md - 共通仕様

### 国際化 (i18n)

- `LANG` 環境変数で言語検出
- `ja_JP` プレフィックス → 日本語、その他 → 英語
- 対象: ヘルプ、エラーメッセージ、確認メッセージ

### カラー出力

- TTY自動検出（`go-isatty` 使用）
- ターミナル出力時のみカラー有効
- パイプ/ファイルリダイレクト時は無効

### 終了コード

- 0: 正常終了
- 1: 全エラー（区別なし）

### データベース

- パス: `~/.djou/djou.db`
- ディレクトリ・テーブルは自動作成
- 手動初期化不要

---

## 4. 各機能仕様 (docs/feature-*.md)

### 4.1 feature-record.md（Phase 3: 記録機能）

- **コマンド**: `djou`（引数なしで実行）
- **入力フィールド**: タスク名(必須), 見積もり(必須), 実績(必須), AIの利用時間(任意), ハマったこと(任意), 解決方法(任意), 学び(任意)
- **キャンセル**: Ctrl+C / Esc で確認ダイアログ表示
- **UIライブラリ**: charmbracelet/huh

> **注意**: 仕様では7フィールドだが、現在の実装は5フィールド（task_name, estimate_hours, actual_hours, memo, tags）

### 4.2 feature-list.md（Phase 4: 一覧表示）

- **コマンド**: `djou list`
- **フラグ**: `--week`/`-w`, `--month`/`-m`, `--limit`/`-l` (デフォルト10), `--tag`, `--no-interactive`
- **排他オプション**: `--week` と `--month` の同時指定不可
- **テーブル表示**: 日付、タスク名(最大20文字で切り詰め)、見積もり、実績
- **インタラクティブモード**: ログ選択 → 編集/削除/キャンセル

### 4.3 feature-search.md（Phase 5: 検索）

- **コマンド**: `djou search "keyword"`
- **検索対象**: task_name, problem, solution, learning（仕様）/ task_name, memo（実装）
- **AND検索**: 複数キーワード指定時は全一致
- **大文字小文字**: 区別なし、部分一致
- **フラグ**: `--limit`/`-l`, `--detail`/`-d`, `--tag`

### 4.4 feature-stats.md（Phase 6: 統計）

- **コマンド**: `djou stats`
- **統計項目**: 記録件数、見積もり合計、実績合計、見積もり精度、AI活用時間/率（仕様）
- **月別表示**: `--month` フラグで月別テーブル表示
- **特定月**: `--month 2025-01` で指定月のみ

> **注意**: AI関連の統計は仕様のみで実装では省略されている

### 4.5 feature-export.md（Phase 7: CSVエクスポート）

- **コマンド**: `djou export`
- **フラグ**: `--output`/`-o` (出力先), `--from` (開始日), `--to` (終了日)
- **ファイル名**: `djou_YYYY-MM.csv` を自動生成（衝突時は連番付与）
- **エンコーディング**: UTF-8 BOM付き（Excel互換）

### 4.6 feature-mcp.md（Phase 8: MCPサーバー）

- **コマンド**: `djou mcp`
- **プロトコル**: MCP (Model Context Protocol) / JSON-RPC 2.0 / stdio
- **提供ツール**: djou_record, djou_list, djou_search, djou_stats, djou_export, djou_update, djou_delete（計7つ）
- **用途**: Claude Code等のMCPクライアントからdjouを操作

---

## 5. docs/go-guidelines.md - Go開発ガイドライン

### プロジェクト構造

```
djou/
├── cmd/djou/main.go       # エントリポイント
├── internal/
│   ├── cmd/               # CLIコマンド
│   ├── db/                # データベース操作
│   ├── model/             # データモデル
│   ├── ui/                # UI部品
│   ├── mcp/               # MCPサーバー
│   ├── config/            # 設定管理
│   ├── version/           # バージョン情報
│   └── testutil/          # テストヘルパー
```

### 主要規約

- **命名**: パッケージ名は小文字・短縮形、インターフェースは Verb+er パターン
- **エラー処理**: `fmt.Errorf("failed to X: %w", err)` でラップ、終了コードは常に1
- **テスト**: テーブル駆動テスト必須、カバレッジ目標80%以上、`:memory:` SQLite使用
- **TDDサイクル**: Red → Green → Blue (仕様書読む → テスト書く → 実装 → リファクタ)
- **品質**: `go fmt`, `go vet`, import順序(stdlib → 外部 → プロジェクト)

---

## 6. docs/roadmap.md / roadmap-v2.md - 開発ロードマップ

### フェーズ一覧と進捗

| Phase | 機能 | 状態 | 依存 |
|-------|------|------|------|
| 1 | DB/Model (CRUD, リポジトリ) | ✅ 完了 | なし |
| 2 | UI共通 (テーブル, カラー, i18n) | ✅ 完了 | なし |
| 3 | Record (インタラクティブ記録) | ✅ 完了 | Phase 1, 2 |
| 4 | List (一覧表示) | ✅ 完了 | Phase 1, 2 |
| 5 | Search (キーワード検索) | ✅ 完了 | Phase 1, 2 |
| 6 | Stats (統計表示) | ✅ 完了 | Phase 1, 2 |
| 7 | Export (CSVエクスポート) | ✅ 完了 | Phase 1 |
| 8 | MCP Server | ✅ 完了 | Phase 1-7 全て |

### 開発フロー

```
仕様書読む → テスト書く(Red) → 実装(Green) → リファクタ(Blue) → make ci → PR → マージ
```

---

## 7. docs/ci-cd.md - CI/CDパイプライン

### CI (GitHub Actions)

| ジョブ | 内容 |
|--------|------|
| lint | `go vet` + `gofmt` チェック |
| build | `go build -v ./...` |
| test | `go test -v -race -coverprofile` (80%カバレッジ警告) |

- **トリガー**: main/develop への push・PR
- **カバレッジ**: 80%を閾値とする（現在は警告のみ）

### CD (GoReleaser + release-please)

- release-please がCHANGELOG自動生成・リリースPR作成
- マージ後にGoReleaserがマルチプラットフォームビルド
- Homebrew tap への公開（`kca2250/homebrew-tap`）
- 対象: linux/darwin/windows × amd64/arm64

---

## 8. CHANGELOG.md - 変更履歴

### リリース履歴

| バージョン | 日付 | 主な変更 |
|-----------|------|---------|
| v0.1.3 | 2026-01-21 | Homebrew tap 対応 |
| v0.1.2 | 2026-01-21 | release-please + GoReleaser 統合 |
| v0.1.1 | 2026-01-21 | 全機能実装（record〜mcp）、タグ、i18n、インタラクティブモード |
| v0.1.0 | 2026-01-21 | 初期リリース |

---

## 9. CI/CD設定ファイル

### .github/workflows/ci.yml

3つのジョブ（lint, build, test）をmain/developブランチで実行。カバレッジレポートをアーティファクトとして7日間保持。

### .github/workflows/release-please.yml

mainへのpushでrelease-pleaseが起動。リリース作成時にGoReleaserがビルド・Homebrew tap公開を実行。

### .github/workflows/release.yml

`v*` タグプッシュでGoReleaserを実行するレガシー/手動リリース用ワークフロー。

### .goreleaser.yml

- CGO_ENABLED=0 でクロスコンパイル
- ldflags でバージョン情報注入
- tar.gz (Unix) / zip (Windows)
- SHA256チェックサム生成
- Homebrew Formula 自動生成

### release-please-config.json / .release-please-manifest.json

Conventional Commitsに基づく自動バージョニング・CHANGELOG生成の設定。

---

## 10. 仕様と実装の差異

| 項目 | 仕様 (djou-spec.md) | 実装 |
|------|---------------------|------|
| DB フィールド | ai_minutes, problem, solution, learning あり | memo, tags に簡略化 |
| 検索対象 | task_name, problem, solution, learning | task_name, memo, tags |
| 統計 | AI活用時間/率を表示 | 見積もり精度のみ |
| 記録フォーム | 7フィールド | 5フィールド（task, estimate, actual, memo, tags） |
| Quick Record | 仕様なし | `-t`, `-e`, `-a` フラグで実装済み |
| MCP update/delete | 仕様なし | 実装済み（djou_update, djou_delete） |
