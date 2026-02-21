# djou - Dev Journal CLI 仕様書

開発日誌（Dev Journal）をCLIで記録・管理するツール

## 概要

| 項目 | 内容 |
|------|------|
| ツール名 | djou |
| 正式名称 | Dev Journal |
| 言語 | Go |
| データ保存 | SQLite |
| 保存先 | `~/.djou/djou.db` |
| インストール | `go install` / `brew tap` |

## 目的

- 日々の開発作業を記録してデータを蓄積
- 定量データ（見積もり時間、実績時間）と定性データ（メモ）を記録
- タグによる分類・フィルタリング
- 蓄積したデータをベンチマークとして組織に展開

## コマンド一覧

### 記録

```bash
djou                                    # 対話形式で記録
djou record -t "タスク" -e 2 -a 3       # Quick Record（即時記録）
```

対話形式で開発日誌を記録する。`-t`, `-e`, `-a` フラグで必須項目を指定すると即時記録も可能。

### 一覧表示

```bash
djou list              # 直近の記録を表示（対話モード）
djou list --week       # 今週の記録
djou list --month      # 今月の記録
djou list --tag "project:案件A"  # タグでフィルタリング
djou list --no-interactive       # 非対話モード
```

### 検索

```bash
djou search "キーワード"
```

タスク名、メモからキーワード検索。

### 集計

```bash
djou stats             # 全体の集計
djou stats --month     # 月別集計
```

表示項目：
- 合計作業時間（見積もり / 実績）
- 見積もり精度（見積もり / 実績）
- 記録件数

### CSV出力

```bash
djou export                    # カレントディレクトリに出力
djou export --output ./path    # 出力先を指定
```

### 設定

```bash
djou config init    # 設定ファイルを初期化
djou config show    # 設定を表示
djou config edit    # エディタで編集
```

### バージョン

```bash
djou version
```

### MCP サーバー

```bash
djou mcp    # MCP サーバーとして起動（stdio）
```

Claude Code などの MCP クライアントから djou の機能を利用可能にする。

## 入力項目

### 対話形式の流れ

```bash
$ djou

? タスク名: ログイン画面実装
? 見積時間 (h): 2
? 実績時間 (h): 3
? メモ: CORSエラーでハマった。プロキシ設定で解決。
? タグ: task_type:新機能, project:案件A

✅ ログを記録しました
```

### 項目詳細

| 項目 | 必須 | 型 | 説明 |
|------|------|-----|------|
| タスク名 | ⭕ | string | 作業したタスクの名前 |
| 見積時間(h) | ⭕ | float | 見積もり時間（時間単位） |
| 実績時間(h) | ⭕ | float | 実際にかかった時間（時間単位） |
| メモ | ❌ | string | 作業に関するメモ（複数行可） |
| タグ | ❌ | string | カンマ区切りのタグ（例: `task_type:新機能, project:案件A`） |

※ 日付は記録時に自動挿入
※ 任意項目はEnterキーでスキップ可能

## データベース設計

### テーブル: logs

| カラム | 型 | 説明 |
|--------|-----|------|
| id | INTEGER PRIMARY KEY | 自動採番 |
| created_at | DATETIME | 記録日時（自動） |
| task_name | TEXT NOT NULL | タスク名 |
| estimate_hours | REAL NOT NULL | 見積もり時間 |
| actual_hours | REAL NOT NULL | 実績時間 |
| memo | TEXT | メモ |
| tags | TEXT | タグ（カンマ区切り） |

## インストール

### go install

```bash
go install github.com/kca2250/djou/cmd/djou@latest
```

### Homebrew

```bash
brew tap kca2250/tap
brew install djou
```

### ソースからビルド

```bash
git clone https://github.com/kca2250/dev-journal.git
cd dev-journal
go build -o djou ./cmd/djou
```

## 技術スタック

- Go
- SQLite（modernc.org/sqlite）
- 対話UI（github.com/charmbracelet/huh）
- CLI（github.com/spf13/cobra）
- MCP（github.com/mark3labs/mcp-go）

## 出力例

### djou list

```
┌────────────┬──────────────────┬──────┬──────┐
│ 日付       │ タスク           │ 見積 │ 実績 │
├────────────┼──────────────────┼──────┼──────┤
│ 2025-01-10 │ ログイン画面実装 │ 2.0h │ 3.0h │
│ 2025-01-10 │ API連携          │ 1.5h │ 1.5h │
│ 2025-01-09 │ 環境構築         │ 1.0h │ 2.0h │
└────────────┴──────────────────┴──────┴──────┘
```

### djou stats

```
統計情報

期間: 2025-01-01 〜 2025-01-10
総ログ数: 15件

合計時間: 見積 25.0h / 実績 32.5h
見積精度: 77%
```

### djou export

```bash
$ djou export
✅ djou_2025-01.csv に出力しました
```

出力CSV形式：
```csv
date,task_name,estimate_hours,actual_hours,memo,tags
2025-01-10,ログイン画面実装,2.0,3.0,CORSエラーでハマった,task_type:新機能
2025-01-10,API連携,1.5,1.5,,
```
