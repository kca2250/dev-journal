# djou - 実装まとめ

## アーキテクチャ概要

```
cmd/djou/main.go           # エントリポイント（cmd.Execute() を呼ぶだけ）
internal/
├── cmd/                   # CLIコマンド層（Cobra）
│   ├── root.go            # ルートコマンド + デフォルト記録処理
│   ├── record.go          # 明示的記録コマンド（quick record対応）
│   ├── list.go            # 一覧表示 + インタラクティブ操作
│   ├── search.go          # キーワード検索
│   ├── stats.go           # 統計表示
│   ├── export.go          # CSVエクスポート
│   ├── config.go          # 設定管理（init/show/edit）
│   ├── mcp.go             # MCPサーバー起動
│   └── version.go         # バージョン表示
├── db/                    # データアクセス層
│   ├── db.go              # DB接続・スキーマ初期化
│   └── log_repository.go  # CRUD + 検索 + 統計
├── model/                 # データモデル
│   ├── log.go             # Log / LogInput 構造体
│   └── validation.go      # 入力バリデーション
├── config/                # 設定管理
│   └── config.go          # YAML設定ロード/保存
├── mcp/                   # MCPサーバー
│   ├── server.go          # サーバー初期化 + stdio起動
│   └── tools.go           # 7ツール定義 + ハンドラ
├── ui/                    # UI部品
│   ├── color.go           # ANSIカラー + TTY検出
│   ├── form.go            # 記録/編集フォーム（huh）
│   ├── i18n.go            # 国際化（ja/en）
│   ├── interactive.go     # 選択/確認ダイアログ
│   ├── messages.go        # メッセージ定数 + 翻訳マップ
│   └── table.go           # ASCIIテーブル描画（CJK対応）
├── version/               # バージョン情報
│   └── version.go         # ビルド時ldflags注入
└── testutil/              # テストユーティリティ
    └── db.go              # インメモリSQLiteセットアップ
```

### データフロー

```
ユーザー入力 → Cobra Command → バリデーション → LogRepository → SQLite
SQLite → LogRepository → Model → UI フォーマット → ターミナル出力
```

---

## 1. エントリポイント

### cmd/djou/main.go

```go
func main() { cmd.Execute() }
```

`cmd.Execute()` を呼ぶだけの最小実装。

---

## 2. コマンド層 (internal/cmd/)

### 2.1 root.go - ルートコマンド

- **コマンド名**: `djou`
- **デフォルト動作**: 引数なし実行でインタラクティブ記録フォーム起動
- **処理フロー**: DB接続 → ローカライザー初期化 → フォーム表示 → バリデーション → DB保存
- **`Execute()`**: ルートコマンド実行、エラー時 `os.Exit(1)`

### 2.2 record.go - 記録コマンド

- **2つのモード**:
  - **Quick Record**: `-t`, `-e`, `-a` フラグで全必須項目を指定 → 即時保存
  - **Interactive**: フラグ不足時にフォーム表示
- **フラグ**: `-t/--task`, `-e/--estimate`, `-a/--actual`, `-m/--memo`, `--tag`
- **Quick判定**: task, estimate, actual が全て指定されているか

### 2.3 list.go - 一覧表示

- **フラグ**: `--week`/`-w`, `--month`/`-m`, `--limit`/`-l`(デフォルト10), `--tag`, `--no-interactive`
- **排他チェック**: `--week` と `--month` の同時指定で `ErrExclusiveOptions` エラー
- **インタラクティブモード**: ログ選択 → アクション選択（編集/削除/キャンセル）のループ
- **`handleEdit()`**: 既存値をプリセットした編集フォーム → DB更新
- **`handleDelete()`**: 確認ダイアログ → DB削除
- **テーブル描画**: `renderLogTable()` で日付・タスク名・見積もり・実績を表形式出力

### 2.4 search.go - キーワード検索

- **引数**: 1つ以上のキーワード（必須）
- **検索方式**: AND検索（全キーワードが task_name または memo に含まれる）
- **フラグ**: `--limit`/`-l`, `--detail`/`-d`, `--tag`
- **`getMatchFields()`**: 各キーワードがどのフィールドにマッチしたかを特定
- **表示**: コンパクトテーブル or 詳細表示（`--detail`）

### 2.5 stats.go - 統計表示

- **3つのモード**:
  1. `djou stats` → 全体統計（期間、件数、見積もり/実績合計、精度）
  2. `djou stats --month` → 月別テーブル表示
  3. `djou stats --month 2025-01` → 特定月の統計
- **精度計算**: `(estimate / actual) * 100`（actual=0 の場合は 0%）
- **月フォーマット検証**: `YYYY-MM` 形式を `validateMonthFormat()` で検証

### 2.6 export.go - CSVエクスポート

- **出力先決定順**: `--output` フラグ > config設定 > カレントディレクトリ
- **ファイル名生成**: `djou_YYYY-MM.csv`、衝突時は `_1`, `_2` ... と連番付与
- **UTF-8 BOM**: Excel互換性のため `0xEF 0xBB 0xBF` を先頭に書き込み
- **日付フィルタ**: `--from`, `--to` で `YYYY-MM-DD` 形式の期間指定
- **エラー型**: `ErrInvalidDateFormat`, `ErrDirectoryNotFound`

### 2.7 config.go - 設定管理

- **3つのサブコマンド**:
  - `config init`: `~/.djou/config.yaml` にテンプレート作成
  - `config show`: 設定ファイルの内容表示
  - `config edit`: `EDITOR`/`VISUAL`/`vi` で外部エディタ起動
- **テンプレート**: タグ分類（task_type, project, tech_area）とexport設定を含むYAML

### 2.8 mcp.go - MCPサーバー

- `mcp.Run()` に委譲するだけの薄いラッパー

### 2.9 version.go - バージョン表示

- `version.GetFullVersion()` を出力

---

## 3. データアクセス層 (internal/db/)

### 3.1 db.go - DB接続管理

- **DB構造体**: `*sql.DB` のラッパー
- **`Open()`**: `~/.djou/djou.db` を開き、ディレクトリ・テーブルを自動作成
- **スキーマ**: logs テーブル（id, created_at, task_name, estimate_hours, actual_hours, memo, tags）

### 3.2 log_repository.go - リポジトリ

**構造体**:

```go
type ListOptions struct { Limit int; Week bool; Month bool; Tag string }
type SearchOptions struct { Keywords []string; Limit int; Tag string }
type Stats struct { Count int; MinDate, MaxDate string; TotalEstimate, TotalActual float64 }
type MonthlyStats struct { Month string; Count int; TotalEstimate, TotalActual float64 }
```

**主要メソッド**:

| メソッド | SQL操作 | 概要 |
|---------|---------|------|
| `Create(LogInput)` | INSERT | 新規記録作成（memo/tagsはNULL対応） |
| `List(ListOptions)` | SELECT | 一覧取得（week/month/tag/limitフィルタ） |
| `Search(keywords, limit)` | SELECT + LIKE | AND検索（task_name, memo対象） |
| `SearchWithOptions(SearchOptions)` | SELECT + LIKE | タグ付き検索 |
| `GetStats()` | COUNT, SUM, MIN, MAX | 全体統計 |
| `GetMonthlyStats()` | GROUP BY strftime | 月別統計 |
| `GetStatsByMonth(month)` | WHERE strftime | 特定月統計 |
| `Export(from, to)` | SELECT + WHERE | 期間指定エクスポート |
| `GetById(id)` | SELECT | ID指定取得 |
| `Update(Log)` | UPDATE | 既存記録更新 |
| `Delete(id)` | DELETE | 記録削除 |

- **NULL処理**: memo, tags は `sql.NullString` で取得し、Valid ならポインタにセット
- **日付フィルタ**: weekは月曜〜日曜、monthは月初〜翌月初
- **ソート**: `ORDER BY created_at DESC`（新しい順）

---

## 4. データモデル (internal/model/)

### 4.1 log.go

```go
type Log struct {
    ID            int
    CreatedAt     time.Time
    TaskName      string
    EstimateHours float64
    ActualHours   float64
    Memo          *string  // nilで未入力
    Tags          *string  // nilで未入力
}

type LogInput struct {
    TaskName      string
    EstimateHours float64
    ActualHours   float64
    Memo          *string
    Tags          *string
}
```

### 4.2 validation.go

- **`ValidateInput(LogInput)`**: 3つのバリデーションルール
  - `ErrEmptyTaskName`: タスク名が空（空白トリム後）
  - `ErrInvalidEstimate`: 見積もりが0以下
  - `ErrInvalidActual`: 実績が0以下

---

## 5. 設定管理 (internal/config/)

### config.go

```go
type Config struct {
    Tags   TagConfig    `yaml:"tags"`
    Export ExportConfig `yaml:"export"`
}
type TagConfig struct {
    TaskType []string `yaml:"task_type"`
    Project  []string `yaml:"project"`
    TechArea []string `yaml:"tech_area"`
}
type ExportConfig struct {
    OutputDir string `yaml:"output_dir"`
}
```

- **`Load()`**: `~/.djou/config.yaml` を読み込み（不在時はデフォルト値）
- **`GetExportOutputDir()`**: `DJOU_EXPORT_DIR` 環境変数 > config設定 > 空文字
- **`OutputDir()`**: `~` をホームディレクトリに展開
- **`TemplateYAML()`**: コメント付きYAMLテンプレート文字列を返す

---

## 6. MCPサーバー (internal/mcp/)

### 6.1 server.go

- **`NewServer()`**: mcp-goでサーバー作成、7ツールを登録
- **`Run()`**: stdioトランスポートでサーバー起動

### 6.2 tools.go - 7つのMCPツール

| ツール名 | 必須パラメータ | 概要 |
|---------|--------------|------|
| `djou_record` | task_name, estimate_hours, actual_hours | ログ記録 |
| `djou_list` | なし | 一覧取得（filter, limit） |
| `djou_search` | keywords | キーワード検索 |
| `djou_stats` | なし | 統計表示（month, monthly_list） |
| `djou_export` | なし | CSV出力（output, from, to） |
| `djou_update` | id | ログ更新（既存値とマージ） |
| `djou_delete` | id | ログ削除 |

- **updateの実装**: GetById → 差分マージ → Update のパターン
- **出力形式**: Markdownテーブル（list, search）、フォーマットテキスト（stats）
- **ヘルパー**: `truncateString()`, `getMatchFields()`, `generateFilename()`, `writeCSV()`, `escapeCSV()`

---

## 7. UI層 (internal/ui/)

### 7.1 color.go - カラー出力

- **ANSI定数**: Red, Green, Yellow, Blue, Magenta, Cyan, White, Bold, Reset
- **`ColorWriter`**: io.Writer ラッパー、`colorize` フラグで色有無を制御
- **`IsTTY()`**: `os.ModeCharDevice` でターミナル判定（go-isatty不使用、独自実装）
- **便利メソッド**: `WriteSuccess()`, `WriteError()`, `WriteWarning()`, `WriteInfo()`, `WriteBold()`

### 7.2 form.go - フォーム

- **`FormInput`**: 文字列ベースのフォーム入力（パース前）
- **`RecordForm`**: 新規記録フォーム（huhライブラリ）
- **`EditForm`**: 既存ログのプリセット編集フォーム
- **バリデーション**: `ValidateTaskName()`, `ValidatePositiveFloat()`
- **変換**: `FormInput.ToLogInput()` で文字列→型変換
- **Escキー対応**: `NewKeyMapWithEsc()` でEscを終了キーに設定

### 7.3 i18n.go - 国際化

- **`Localizer`**: 言語設定とメッセージ辞書を管理
- **言語検出**: `LANG` 環境変数を解析（`ja_JP.UTF-8` → `ja`）
- **サポート言語**: ja, en（デフォルト: en）
- **メッセージ取得**: `Get(key)` / `Getf(key, args...)` でフォーマット付き取得
- **グローバル**: `SetGlobalLocalizer()` / `T()` / `Tf()` でシングルトンアクセス

### 7.4 interactive.go - インタラクティブUI

- **`LogSelector`**: ログ一覧から選択するセレクトメニュー
- **`ActionSelector`**: 操作選択（Edit/Delete/Cancel）
- **`DeleteConfirm`**: 削除確認ダイアログ
- **アクション定数**: `ActionEdit`, `ActionDelete`, `ActionCancel`

### 7.5 messages.go - メッセージ定義

- **約50のメッセージキー**: 各コマンドのUI文字列を定数定義
- **翻訳マップ**: `messages["ja"]`, `messages["en"]` の2言語マップ
- **カバー範囲**: record, list, search, stats, export, config, テーブルヘッダー, インタラクティブ操作

### 7.6 table.go - テーブル描画

- **`RenderTable(headers, rows)`**: 罫線付きASCIIテーブル描画
- **CJK対応**: `isWideRune()` で全角文字判定、`displayWidth()` で表示幅計算
- **ユーティリティ**: `FormatHours()` (X.Xh形式), `TruncateString()`, `TruncateStringDisplay()`, `StringWidth()`
- **罫線文字**: `┌─┬┐ │ ├─┼┤ └─┴┘` (Unicode Box Drawing)

---

## 8. バージョン管理 (internal/version/)

### version.go

- **ビルド時注入変数**: `Version`, `Commit`, `Date`（ldflags経由）
- **`GoVersion`**: `runtime.Version()` で自動取得
- **デフォルト値**: Version="dev", Commit="unknown", Date="unknown"

---

## 9. テストユーティリティ (internal/testutil/)

### db.go

- **`SetupTestDB(t)`**: インメモリSQLite (`:memory:`) を開き、スキーマ適用
- **クリーンアップ**: `t.Cleanup()` で自動クローズ

---

## 10. テストカバレッジ

### テストファイル一覧と概要

| テストファイル | テスト数 | 主な検証内容 |
|--------------|---------|------------|
| root_test.go | 3 | コマンドメタデータ（名前、説明、ハンドラ存在） |
| record_test.go | - | （テストファイルなし） |
| list_test.go | 4 | コマンド構造、排他エラー、テーブル描画、長いタスク名の切り詰め |
| search_test.go | 4 | コマンド構造、マッチフィールド検出(5ケース)、テーブル描画 |
| stats_test.go | 5 | 精度計算(5ケース)、月別テーブル、空データ、月フォーマット検証(7ケース) |
| export_test.go | 7 | ファイル名生成(3ケース)、日付検証(7ケース)、CSV行フォーマット、BOM書き込み |
| config_test.go | 4 | コマンド名、サブコマンド存在、ファイル作成 |
| mcp_test.go | 3 | コマンドメタデータ |
| log_repository_test.go | 6 | CRUD操作(2ケース)、一覧(3ケース)、検索(5ケース)、統計(2ケース) |
| validation_test.go | 1 | 8バリデーションシナリオ（テーブル駆動テスト） |
| color_test.go | 6 | カラー出力、TTY検出 |
| form_test.go | 5 | 正数パース(7ケース)、タスク名検証(3ケース)、float検証(4ケース)、変換(5ケース) |
| i18n_test.go | 7 | 言語初期化(6ケース)、メッセージ取得、フォーマット、正規化(10ケース) |
| table_test.go | 4 | テーブル描画(3ケース)、時間フォーマット、切り詰め(4ケース) |
| tools_test.go | 4 | マッチ検出(5ケース)、文字列検索、CSVエスケープ(4ケース)、ファイル名生成 |
| server_test.go | 2 | ツール登録数(7)、定数検証 |

### テスト方針

- **テーブル駆動テスト**: 全テストで統一的に使用
- **インメモリDB**: `:memory:` SQLiteで高速テスト
- **UI層**: フォーム実行はテスト対象外、バリデーション関数のみテスト
- **MCP層**: ユーティリティ関数とサーバー初期化のみテスト

---

## 11. ビルド・リリース

### Makefile ターゲット

| ターゲット | 内容 |
|-----------|------|
| `build` | `go build ./...` |
| `test` | `go test -race ./...` |
| `coverage` | カバレッジレポート生成 |
| `lint` | `go vet` + `gofmt` チェック |
| `fmt` | `go fmt ./...` |
| `ci` | lint + build + test |
| `dev` | `./djou` バイナリをビルド |

### .goreleaser.yml

- **ビルド設定**: CGO_ENABLED=0、linux/darwin/windows × amd64/arm64
- **バージョン注入**: `-ldflags` で Version, Commit, Date を埋め込み
- **アーカイブ**: tar.gz (Unix) / zip (Windows)
- **配布**: GitHub Release + Homebrew tap (`kca2250/homebrew-tap`)

### go.mod

- **モジュール**: `github.com/kca2250/djou`
- **Go バージョン**: 1.25.6
- **主要依存**: cobra, huh, mcp-go, modernc.org/sqlite

---

## 12. 主要な設計パターン

| パターン | 適用箇所 |
|---------|---------|
| **Repository** | LogRepository がDB操作をカプセル化 |
| **Command** | Cobra フレームワークでサブコマンド分離 |
| **Factory** | NewServer, NewLocalizer, NewColorWriter |
| **Singleton** | グローバル Localizer インスタンス |
| **Layered Architecture** | cmd → db/model → SQLite の階層構造 |
| **Dependency Injection** | *sql.DB をリポジトリに注入 |
| **Null Object** | ポインタ型でオプショナルフィールドを表現 |
