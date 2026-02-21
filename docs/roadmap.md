# 開発ロードマップ

djou の実装順序と依存関係を整理したドキュメント。

---

## 開発プロセス: TDD

各Phaseは **テスト駆動開発 (TDD)** で進める。

```
┌─────────────────────────────────────────────────────────────┐
│                    Phase の進め方                            │
├─────────────────────────────────────────────────────────────┤
│  1. 📝 仕様書確認 (docs/feature-*.md)                       │
│  2. 🔴 テストを書く (失敗する状態)                           │
│  3. 🟢 実装する (テストが通る最小限)                         │
│  4. 🔵 リファクタリング                                      │
│  5. ✅ make ci で全チェック通過を確認                        │
│  6. 📤 PR作成 → CIパス → マージ                             │
└─────────────────────────────────────────────────────────────┘
```

### 各Phaseの作業手順

```bash
# 1. 仕様書を確認
cat docs/feature-*.md

# 2. テストを先に書く
# 3. テスト実行（Red: 失敗を確認）
go test ./...

# 4. 実装
# 5. テスト実行（Green: 成功を確認）
go test ./...

# 6. リファクタリング
# 7. CI確認
make ci

# 8. PRを作成（CI通過後にマージ）
```

---

## 依存関係図

```
                    ┌─────────────────────────────┐
                    │         共通基盤            │
                    │  Phase 1: DB操作・モデル    │
                    │  Phase 2: UI共通部品        │
                    └──────────────┬──────────────┘
                                   │
                                   ▼
                    ┌─────────────────────────────┐
                    │       記録機能 (djou)       │
                    │         Phase 3             │
                    └──────────────┬──────────────┘
                                   │
          ┌────────────────────────┼────────────────────────┐
          │                        │                        │
          ▼                        ▼                        ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  一覧 (list)    │    │  検索 (search)  │    │  集計 (stats)   │
│    Phase 4      │    │    Phase 5      │    │    Phase 6      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                   │
                                   ▼
                    ┌─────────────────────────────┐
                    │    CSV出力 (export)         │
                    │         Phase 7             │
                    └─────────────────────────────┘
```

---

## Phase 1: 共通基盤 - DB操作・モデル

### 概要

全機能が依存するデータベース操作とモデルを実装する。

### 実装内容

| ファイル | 内容 |
|----------|------|
| `internal/model/log.go` | Log構造体（task_name, estimate_hours, actual_hours, memo, tags） |
| `internal/model/validation.go` | 入力バリデーション |
| `internal/db/db.go` | DB接続・初期化 |
| `internal/db/log_repository.go` | CRUD操作（Create, Read, Search, Aggregate） |

### 詳細タスク

- [x] `LogRepository` 実装
- [x] `Create(log *LogInput) error` - 記録の挿入
- [x] `List(opts ListOptions) ([]Log, error)` - 一覧取得（limit, week, month, tag対応）
- [x] `Search(keywords []string, limit int) ([]Log, error)` - キーワード検索
- [x] `GetStats() (*Stats, error)` - 全体統計
- [x] `GetMonthlyStats() ([]MonthlyStats, error)` - 月別統計
- [x] `Export(from, to time.Time) ([]Log, error)` - エクスポート用データ取得
- [x] `GetById(id int64) (*Log, error)` - ID指定取得
- [x] `Update(log *Log) error` - 記録更新
- [x] `Delete(id int64) error` - 記録削除

### 完了条件

- 全てのDB操作関数が実装されている
- 単体テストが通る

---

## Phase 2: 共通基盤 - UI共通部品

### 概要

複数機能で共用するUI部品を実装する。

### 実装内容

| ファイル | 内容 |
|----------|------|
| `internal/ui/table.go` | テーブル出力（CJK文字対応） |
| `internal/ui/color.go` | カラー出力判定（TTY自動検出） |
| `internal/ui/i18n.go` | 多言語対応（日本語/英語切替） |
| `internal/ui/messages.go` | メッセージ定義 |

### 詳細タスク

- [x] `RenderTable(headers []string, rows [][]string)` - テーブル描画
- [x] `IsTTY()` - TTY判定（os.ModeCharDevice使用）
- [x] `Localizer` - 言語判定・メッセージ取得
- [x] メッセージ定数の定義（成功、エラー、確認）

### 完了条件

- テーブル出力が動作する
- 多言語メッセージが切り替わる

---

## Phase 3: 記録機能 (djou)

### 概要

対話形式で開発日誌を記録するメイン機能。

### 仕様書

[feature-record.md](./feature-record.md)

### 実装内容

| ファイル | 内容 |
|----------|------|
| `internal/cmd/root.go` | ルートコマンドの実装 |
| `internal/cmd/record.go` | 記録コマンド（Quick Record対応） |
| `internal/ui/form.go` | 対話フォーム（huhライブラリ使用） |

### 詳細タスク

- [x] 対話フォームの実装（5項目: タスク名, 見積, 実績, メモ, タグ）
- [x] バリデーション実装
- [x] Ctrl+C / Esc での中断処理
- [x] DB保存処理
- [x] 成功メッセージ表示
- [x] Quick Record モード（-t, -e, -a フラグ）

### 完了条件

- `djou` コマンドで記録ができる
- 任意項目をEnterでスキップできる
- Quick Record で即時記録ができる

### 依存

- Phase 1 (DB操作)
- Phase 2 (UI共通部品)

---

## Phase 4: 一覧表示機能 (djou list)

### 概要

記録した開発日誌を一覧表示する。

### 仕様書

[feature-list.md](./feature-list.md)

### 実装内容

| ファイル | 内容 |
|----------|------|
| `internal/cmd/list.go` | listサブコマンド |
| `internal/ui/interactive.go` | インタラクティブUI（選択・編集・削除） |

### 詳細タスク

- [x] `list` サブコマンドの追加
- [x] `--week`, `--month`, `--limit` オプション実装
- [x] `--tag` オプション（タグフィルター）
- [x] `--no-interactive` オプション
- [x] オプション排他チェック
- [x] テーブル形式での出力（日付, タスク, 見積, 実績）
- [x] インタラクティブモード（編集・削除・キャンセル）

### 完了条件

- `djou list` で一覧が表示される
- 各オプションが正しく動作する
- `--week` と `--month` の同時指定でエラー
- 対話モードで編集・削除ができる

### 依存

- Phase 1 (DB操作)
- Phase 2 (テーブル出力)

---

## Phase 5: 検索機能 (djou search)

### 概要

キーワードで開発日誌を検索する。

### 仕様書

[feature-search.md](./feature-search.md)

### 実装内容

| ファイル | 内容 |
|----------|------|
| `internal/cmd/search.go` | searchサブコマンド |

### 詳細タスク

- [x] `search` サブコマンドの追加
- [x] 複数キーワードのAND検索（task_name, memo対象）
- [x] `--limit` オプション実装
- [x] `--tag` オプション（タグフィルター）
- [x] `--detail` オプション（詳細表示モード）
- [x] 検索結果のテーブル表示

### 完了条件

- `djou search "キーワード"` で検索できる
- 複数キーワードでAND検索される
- 詳細表示が動作する

### 依存

- Phase 1 (DB操作)
- Phase 2 (テーブル出力)

---

## Phase 6: 集計機能 (djou stats)

### 概要

記録の統計情報を表示する。

### 仕様書

[feature-stats.md](./feature-stats.md)

### 実装内容

| ファイル | 内容 |
|----------|------|
| `internal/cmd/stats.go` | statsサブコマンド |

### 詳細タスク

- [x] `stats` サブコマンドの追加
- [x] 全体統計の計算・表示（件数, 期間, 見積/実績合計, 精度）
- [x] `--month` オプション（月別一覧 or 特定月指定）
- [x] ゼロ除算時は0%表示

### 完了条件

- `djou stats` で統計が表示される
- `djou stats --month` で月別一覧が表示される
- `djou stats --month 2025-01` で特定月の統計が表示される

### 依存

- Phase 1 (DB操作)
- Phase 2 (テーブル出力)

---

## Phase 7: CSV出力機能 (djou export)

### 概要

記録をCSVファイルに出力する。

### 仕様書

[feature-export.md](./feature-export.md)

### 実装内容

| ファイル | 内容 |
|----------|------|
| `internal/cmd/export.go` | exportサブコマンド |

### 詳細タスク

- [x] `export` サブコマンドの追加
- [x] `--output` オプション（出力先指定、設定ファイル対応）
- [x] `--from`, `--to` オプション（期間指定）
- [x] ファイル名自動生成（連番付与）
- [x] UTF-8 BOM付き出力
- [x] CSVエスケープ処理

### 完了条件

- `djou export` でCSVが出力される
- 同名ファイル存在時に連番が付与される
- 期間指定が動作する

### 依存

- Phase 1 (DB操作)

---

## 進捗管理

| Phase | 状態 |
|-------|------|
| Phase 1: DB操作・モデル | ✅ 完了 |
| Phase 2: UI共通部品 | ✅ 完了 |
| Phase 3: 記録機能 | ✅ 完了 |
| Phase 4: 一覧表示機能 | ✅ 完了 |
| Phase 5: 検索機能 | ✅ 完了 |
| Phase 6: 集計機能 | ✅ 完了 |
| Phase 7: CSV出力機能 | ✅ 完了 |

---

## 備考

- Phase 4〜6 は依存関係がないため、並行して実装可能
- Phase 7 は他の機能と独立しているが、テスト時にデータが必要なため後半に配置
- 各Phaseの完了後、動作確認を行ってから次のPhaseに進む
