# Go 実装ガイドライン

djou プロジェクトにおける Go の実装方針とベストプラクティス。

## 使用バージョン

- Go 1.25.x

---

## プロジェクト構造

### 標準レイアウト

```
djou/
├── cmd/
│   └── djou/
│       └── main.go       # エントリーポイント（最小限に保つ）
├── internal/             # 外部に公開しないパッケージ
│   ├── cmd/              # CLIコマンド定義
│   ├── db/               # データベース操作
│   ├── model/            # データモデル
│   └── ui/               # UI関連
├── docs/                 # ドキュメント
├── go.mod
└── go.sum
```

### ルール

- `cmd/` にはエントリーポイントのみ配置
- ビジネスロジックは `internal/` 配下に配置
- `internal/` は外部パッケージからインポート不可（Goの仕様）

---

## 命名規則

### パッケージ名

```go
// Good
package db
package model
package ui

// Bad
package database      // 冗長
package models        // 複数形は避ける
package utils         // 曖昧すぎる
```

### 変数・関数名

```go
// Good: キャメルケース、短く意味のある名前
func (r *LogRepository) Create(log *LogInput) error
var userID int
var db *sql.DB

// Bad: 冗長、アンダースコア
func (r *LogRepository) CreateNewLogEntry(logEntry *LogInput) error
var user_id int
```

### インターフェース名

```go
// Good: 動詞 + er（単一メソッドの場合）
type Reader interface {
    Read(p []byte) (n int, err error)
}

type LogRepository interface {
    Create(log *LogInput) error
    List(opts ListOptions) ([]Log, error)
}

// Bad
type ILogRepository interface {}  // I プレフィックスは不要
```

### 定数

```go
// Good: キャメルケース
const defaultLimit = 10
const MaxRetries = 3

// Bad
const DEFAULT_LIMIT = 10  // スネークケースは避ける
```

---

## エラーハンドリング

### 基本方針

```go
// Good: エラーを即座にチェック
result, err := doSomething()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// Bad: エラーを無視
result, _ := doSomething()
```

### エラーのラップ

```go
// Good: %w でラップして元のエラーを保持
if err != nil {
    return fmt.Errorf("failed to open database: %w", err)
}

// 呼び出し側でエラーを判定可能
if errors.Is(err, sql.ErrNoRows) {
    // 処理
}
```

### カスタムエラー

```go
// シンプルなエラー定義
var (
    ErrNotFound     = errors.New("record not found")
    ErrInvalidInput = errors.New("invalid input")
)

// 使用例
if len(logs) == 0 {
    return nil, ErrNotFound
}
```

### エラーメッセージ

```go
// Good: 小文字で始める、ピリオドなし
return fmt.Errorf("failed to connect database: %w", err)

// Bad: 大文字で始める、ピリオドあり
return fmt.Errorf("Failed to connect database.: %w", err)
```

---

## インターフェース設計

### 小さく保つ

```go
// Good: 必要最小限のメソッド
type LogCreator interface {
    Create(log *LogInput) error
}

type LogLister interface {
    List(opts ListOptions) ([]Log, error)
}

// 必要に応じて組み合わせ
type LogRepository interface {
    LogCreator
    LogLister
    Search(keywords []string, limit int) ([]Log, error)
}
```

### 消費側で定義

```go
// Good: インターフェースは使う側で定義
// internal/cmd/list.go
type logLister interface {
    List(opts db.ListOptions) ([]model.Log, error)
}

func NewListCommand(lister logLister) *cobra.Command {
    // ...
}
```

---

## 構造体設計

### コンストラクタ

```go
// Good: New 関数でコンストラクタを提供
func NewLogRepository(db *sql.DB) *LogRepository {
    return &LogRepository{db: db}
}

// オプションパターン（設定が多い場合）
type Option func(*Config)

func WithTimeout(d time.Duration) Option {
    return func(c *Config) {
        c.Timeout = d
    }
}

func NewClient(opts ...Option) *Client {
    cfg := &Config{Timeout: 30 * time.Second} // デフォルト値
    for _, opt := range opts {
        opt(cfg)
    }
    return &Client{config: cfg}
}
```

### ゼロ値の活用

```go
// Good: ゼロ値が有効なデフォルトになる設計
type ListOptions struct {
    Limit  int  // 0 の場合はデフォルト10件として扱う
    Week   bool
    Month  bool
}

func (o ListOptions) GetLimit() int {
    if o.Limit <= 0 {
        return 10
    }
    return o.Limit
}
```

---

## データベース操作

### コネクション管理

```go
// Good: アプリケーション起動時に1回だけOpen
func main() {
    db, err := db.Open()
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // dbを各コンポーネントに渡す
}
```

### クエリ

```go
// Good: プリペアドステートメントでSQLインジェクション対策
func (r *LogRepository) Search(keyword string) ([]Log, error) {
    rows, err := r.db.Query(
        `SELECT * FROM logs WHERE task_name LIKE ?`,
        "%"+keyword+"%",
    )
    // ...
}

// Bad: 文字列結合（SQLインジェクションの危険）
query := "SELECT * FROM logs WHERE task_name LIKE '%" + keyword + "%'"
```

### トランザクション

```go
// Good: deferでロールバックを保証
func (r *LogRepository) CreateWithTx(log *LogInput) error {
    tx, err := r.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback() // コミット後は無視される

    // 処理...

    return tx.Commit()
}
```

---

## テスト

### テスト方針

| 項目 | 決定 |
|------|------|
| 依存性注入 | **完全DI** - リポジトリ層はインターフェースで定義 |
| DBテスト | **インメモリSQLite** (`:memory:`) |
| UIテスト | **UIとロジックを分離** - ロジック部分のみテスト |
| カバレッジ目標 | **80%以上**（UI層を除く） |
| テスト種類 | **単体テスト + 統合テスト** |
| フレームワーク | **標準testingのみ** |

### ファイル配置

```
internal/db/
├── log_repository.go
├── log_repository_test.go      # 単体テスト
└── log_repository_integration_test.go  # 統合テスト（必要に応じて分離）
```

### テスト関数の命名

```go
// Good: Test + 関数名 + シナリオ
func TestLogRepository_Create(t *testing.T) { }
func TestLogRepository_Create_EmptyTaskName(t *testing.T) { }
func TestLogRepository_List_WithWeekFilter(t *testing.T) { }
```

### テーブル駆動テスト

全てのテストはテーブル駆動で書く。

```go
func TestValidateInput(t *testing.T) {
    tests := []struct {
        name    string
        input   LogInput
        wantErr bool
    }{
        {
            name:    "valid input",
            input:   LogInput{TaskName: "test", EstimateHours: 1, ActualHours: 1},
            wantErr: false,
        },
        {
            name:    "empty task name",
            input:   LogInput{TaskName: "", EstimateHours: 1, ActualHours: 1},
            wantErr: true,
        },
        {
            name:    "negative estimate hours",
            input:   LogInput{TaskName: "test", EstimateHours: -1, ActualHours: 1},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateInput(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateInput() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### インメモリSQLiteでのDBテスト

```go
// テスト用DBのセットアップヘルパー
func setupTestDB(t *testing.T) *DB {
    t.Helper()

    db, err := sql.Open("sqlite", ":memory:")
    if err != nil {
        t.Fatalf("failed to open test db: %v", err)
    }

    // スキーマ作成
    if _, err := db.Exec(schema); err != nil {
        t.Fatalf("failed to create schema: %v", err)
    }

    t.Cleanup(func() {
        db.Close()
    })

    return &DB{db}
}

func TestLogRepository_Create(t *testing.T) {
    db := setupTestDB(t)
    repo := NewLogRepository(db)

    tests := []struct {
        name    string
        input   LogInput
        wantErr bool
    }{
        {
            name: "valid input",
            input: LogInput{
                TaskName:      "テストタスク",
                EstimateHours: 2.0,
                ActualHours:   3.0,
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := repo.Create(context.Background(), &tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### モックを使った単体テスト

```go
// インターフェース定義（消費側で定義）
type logRepository interface {
    List(ctx context.Context, opts ListOptions) ([]model.Log, error)
}

// モック実装
type mockLogRepository struct {
    logs []model.Log
    err  error
}

func (m *mockLogRepository) List(ctx context.Context, opts ListOptions) ([]model.Log, error) {
    return m.logs, m.err
}

// コマンドのテスト
func TestListCommand_Execute(t *testing.T) {
    tests := []struct {
        name     string
        mockLogs []model.Log
        mockErr  error
        wantErr  bool
    }{
        {
            name: "success with logs",
            mockLogs: []model.Log{
                {ID: 1, TaskName: "タスク1"},
            },
            wantErr: false,
        },
        {
            name:    "db error",
            mockErr: errors.New("db connection failed"),
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mock := &mockLogRepository{
                logs: tt.mockLogs,
                err:  tt.mockErr,
            }

            err := executeList(mock, ListOptions{})
            if (err != nil) != tt.wantErr {
                t.Errorf("executeList() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### UIとロジックの分離

```go
// Bad: UIとロジックが混在
func RunRecord() error {
    // フォーム表示
    var input LogInput
    form := huh.NewForm(...)
    form.Run()

    // バリデーション
    if input.TaskName == "" {
        return errors.New("task name required")
    }

    // 保存
    return repo.Create(input)
}

// Good: UIとロジックを分離
// ui/form.go - UIのみ（テスト対象外）
func CollectInput() (*LogInput, error) {
    var input LogInput
    form := huh.NewForm(...)
    if err := form.Run(); err != nil {
        return nil, err
    }
    return &input, nil
}

// service/record.go - ロジック（テスト対象）
func (s *RecordService) Execute(input *LogInput) error {
    if err := ValidateInput(input); err != nil {
        return err
    }
    return s.repo.Create(context.Background(), input)
}

// validator.go - バリデーション（テスト対象）
func ValidateInput(input *LogInput) error {
    if input.TaskName == "" {
        return ErrEmptyTaskName
    }
    if input.EstimateHours <= 0 {
        return ErrInvalidEstimate
    }
    return nil
}
```

### テストカバレッジの確認

```bash
# カバレッジ付きでテスト実行
go test -coverprofile=coverage.out ./...

# カバレッジレポート表示
go tool cover -func=coverage.out

# HTMLレポート生成
go tool cover -html=coverage.out -o coverage.html
```

### テスト対象の優先度

| 優先度 | 対象 | カバレッジ目標 |
|--------|------|----------------|
| 高 | リポジトリ層（CRUD） | 90%以上 |
| 高 | バリデーション | 100% |
| 中 | コマンド層（フラグ解析、ロジック） | 80%以上 |
| 低 | UI層（対話フォーム） | テスト対象外 |

### テストヘルパー

```go
// 共通のテストヘルパーは testutil パッケージに配置
// internal/testutil/db.go

package testutil

import (
    "database/sql"
    "testing"
)

func SetupTestDB(t *testing.T) *sql.DB {
    t.Helper()

    db, err := sql.Open("sqlite", ":memory:")
    if err != nil {
        t.Fatalf("failed to open test db: %v", err)
    }

    t.Cleanup(func() {
        db.Close()
    })

    return db
}

// テストデータ作成ヘルパー
func CreateTestLog(t *testing.T, repo *LogRepository, taskName string) *model.Log {
    t.Helper()

    input := &model.LogInput{
        TaskName:      taskName,
        EstimateHours: 1.0,
        ActualHours:   1.0,
    }

    if err := repo.Create(context.Background(), input); err != nil {
        t.Fatalf("failed to create test log: %v", err)
    }

    // 作成したログを返す
    logs, _ := repo.List(context.Background(), ListOptions{Limit: 1})
    return &logs[0]
}
```

---

## テスト駆動開発 (TDD)

### 基本方針

**仕様ドキュメントを元にテストを先に書き、その後実装する。**

```
┌─────────────────────────────────────────────────────┐
│                    TDD サイクル                      │
├─────────────────────────────────────────────────────┤
│  1. 📝 仕様書を確認 (docs/feature-*.md)             │
│  2. 🔴 Red: 失敗するテストを書く                     │
│  3. 🟢 Green: テストが通る最小限の実装               │
│  4. 🔵 Refactor: コードを整理                        │
│  5. 繰り返し                                         │
└─────────────────────────────────────────────────────┘
```

### 実践例

#### 1. 仕様書を確認

`docs/feature-record.md` から要件を確認：

```markdown
| 項目 | 必須 | バリデーション |
|------|------|----------------|
| タスク名 | ⭕ | 空文字不可 |
| 見積もり(h) | ⭕ | 正の数値 |
```

#### 2. Red: 失敗するテストを書く

```go
func TestValidateInput(t *testing.T) {
    tests := []struct {
        name    string
        input   LogInput
        wantErr error
    }{
        {
            name:    "empty task name should fail",
            input:   LogInput{TaskName: "", EstimateHours: 1, ActualHours: 1},
            wantErr: ErrEmptyTaskName,
        },
        {
            name:    "negative estimate should fail",
            input:   LogInput{TaskName: "test", EstimateHours: -1, ActualHours: 1},
            wantErr: ErrInvalidEstimate,
        },
        {
            name:    "valid input should pass",
            input:   LogInput{TaskName: "test", EstimateHours: 1, ActualHours: 1},
            wantErr: nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateInput(&tt.input)
            if !errors.Is(err, tt.wantErr) {
                t.Errorf("ValidateInput() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

この時点でテストを実行 → **Red（失敗）**

```bash
go test ./...
# ValidateInput が未実装なのでコンパイルエラー or 失敗
```

#### 3. Green: 最小限の実装

```go
var (
    ErrEmptyTaskName   = errors.New("task name is required")
    ErrInvalidEstimate = errors.New("estimate hours must be positive")
)

func ValidateInput(input *LogInput) error {
    if input.TaskName == "" {
        return ErrEmptyTaskName
    }
    if input.EstimateHours <= 0 {
        return ErrInvalidEstimate
    }
    return nil
}
```

テストを実行 → **Green（成功）**

```bash
go test ./...
# PASS
```

#### 4. Refactor: 必要に応じて整理

- 重複コードの削除
- 命名の改善
- テストが通ることを確認

### TDD のルール

| ルール | 説明 |
|--------|------|
| テストなしにコードを書かない | 全ての本番コードはテストから始める |
| 失敗するテストを1つだけ書く | 一度に多くのテストを書かない |
| テストを通す最小限のコードを書く | 過剰な実装をしない |
| リファクタリングはテストが通ってから | Green の状態でのみ整理 |

### 仕様書とテストの対応

| 仕様書 | テストファイル |
|--------|----------------|
| `feature-record.md` | `internal/model/validation_test.go` |
| `feature-list.md` | `internal/db/log_repository_test.go` |
| `feature-search.md` | `internal/db/log_repository_test.go` |
| `feature-stats.md` | `internal/db/stats_test.go` |
| `feature-export.md` | `internal/cmd/export_test.go` |

### 開発フロー

```bash
# 1. 仕様書を確認
cat docs/feature-record.md

# 2. テストファイルを作成
touch internal/model/validation_test.go

# 3. テストを書く（Red）
code internal/model/validation_test.go

# 4. テスト実行（失敗を確認）
go test ./internal/model/...

# 5. 実装する（Green）
code internal/model/validation.go

# 6. テスト実行（成功を確認）
go test ./internal/model/...

# 7. リファクタリング
# 8. 全テスト実行
make test
```

---

## Context の使用

### 基本方針

```go
// Good: 第一引数に context.Context
func (r *LogRepository) List(ctx context.Context, opts ListOptions) ([]Log, error) {
    rows, err := r.db.QueryContext(ctx, query, args...)
    // ...
}

// CLIでは context.Background() を使用
ctx := context.Background()
```

### キャンセル処理

```go
// シグナルハンドリング
ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
defer cancel()
```

---

## 並行処理

### 基本方針

- CLIツールでは基本的に並行処理は不要
- 必要な場合は `sync.WaitGroup` や `errgroup` を使用

```go
import "golang.org/x/sync/errgroup"

func processMultiple(items []Item) error {
    g, ctx := errgroup.WithContext(context.Background())

    for _, item := range items {
        item := item // ループ変数のキャプチャ（Go 1.22以降は不要）
        g.Go(func() error {
            return process(ctx, item)
        })
    }

    return g.Wait()
}
```

---

## ログ出力

### 基本方針

- CLIツールでは `fmt` パッケージで十分
- 標準出力: 正常な出力
- 標準エラー出力: エラーメッセージ

```go
// 正常出力
fmt.Println("✅ 記録しました！")

// エラー出力
fmt.Fprintln(os.Stderr, "エラー: データベースに接続できません")
```

### 構造化ログが必要な場合

```go
import "log/slog"

slog.Info("record created", "task", taskName, "id", id)
```

---

## CLIフラグ

### Cobra の使用

```go
var listCmd = &cobra.Command{
    Use:   "list",
    Short: "記録を一覧表示",
    RunE:  runList,
}

func init() {
    listCmd.Flags().BoolP("week", "w", false, "今週の記録を表示")
    listCmd.Flags().BoolP("month", "m", false, "今月の記録を表示")
    listCmd.Flags().IntP("limit", "l", 10, "表示件数")

    // 排他オプションのチェックは RunE 内で行う
}

func runList(cmd *cobra.Command, args []string) error {
    week, _ := cmd.Flags().GetBool("week")
    month, _ := cmd.Flags().GetBool("month")

    if week && month {
        return errors.New("--week と --month は同時に指定できません")
    }
    // ...
}
```

---

## コードフォーマット

### ツール

```bash
# フォーマット
go fmt ./...

# インポート整理
goimports -w .

# 静的解析
go vet ./...
```

### インポート順序

```go
import (
    // 標準ライブラリ
    "context"
    "fmt"
    "time"

    // サードパーティ
    "github.com/spf13/cobra"

    // 自プロジェクト
    "github.com/kca2250/djou/internal/db"
    "github.com/kca2250/djou/internal/model"
)
```

---

## 避けるべきパターン

### init() の乱用

```go
// Bad: init() でグローバル状態を変更
func init() {
    db, _ = sql.Open("sqlite", "path") // エラー処理できない
}

// Good: 明示的に初期化
func main() {
    db, err := sql.Open("sqlite", "path")
    if err != nil {
        log.Fatal(err)
    }
}
```

### グローバル変数

```go
// Bad: グローバル変数
var globalDB *sql.DB

// Good: 依存性注入
type App struct {
    db *sql.DB
}
```

### panic の使用

```go
// Bad: panic でエラー処理
if err != nil {
    panic(err)
}

// Good: エラーを返す
if err != nil {
    return fmt.Errorf("failed to process: %w", err)
}
```

---

## 参考資料

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
