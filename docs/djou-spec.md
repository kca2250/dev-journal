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
| インストール | `go install`（将来: brew tap） |

## 目的

- 日々の開発作業を記録してデータを蓄積
- 定量データ（作業時間、AI活用時間）と定性データ（学び、ハマりポイント）を両方記録
- 蓄積したデータをベンチマークとして組織に展開

## コマンド一覧

### 記録

```bash
djou
```

対話形式で開発日誌を記録する。

### 一覧表示

```bash
djou list              # 直近の記録を表示
djou list --week       # 今週の記録
djou list --month      # 今月の記録
```

### 検索

```bash
djou search "キーワード"
```

タスク名、ハマったこと、解決方法、学びからキーワード検索。

### 集計

```bash
djou stats             # 全体の集計
djou stats --month     # 月別集計
```

表示項目：
- 合計作業時間（見積もり / 実績）
- 見積もり精度（実績 / 見積もり）
- AI活用時間・AI活用率
- 記録件数

### CSV出力

```bash
djou export                    # カレントディレクトリに出力
djou export --output ./path    # 出力先を指定
```

## 入力項目

### 対話形式の流れ

```bash
$ djou

? タスク名: ログイン画面実装
? 見積もり(h): 2
? 実績(h): 3
? AI活用(min): 30
? ハマったこと: CORSエラー
? 解決方法: プロキシ設定を追加
? 学び: API連携は早めに確認する

✅ 記録しました！
```

### 項目詳細

| 項目 | 必須 | 型 | 説明 |
|------|------|-----|------|
| タスク名 | ⭕ | string | 作業したタスクの名前 |
| 見積もり(h) | ⭕ | float | 見積もり時間（時間単位） |
| 実績(h) | ⭕ | float | 実際にかかった時間（時間単位） |
| AI活用(min) | ❌ | int | AI（Claude Code等）を活用した時間（分単位） |
| ハマったこと | ❌ | string | 作業中にハマった問題 |
| 解決方法 | ❌ | string | どう解決したか |
| 学び | ❌ | string | 得られた知見・気づき |

※ 日付は記録時に自動挿入

## データベース設計

### テーブル: logs

| カラム | 型 | 説明 |
|--------|-----|------|
| id | INTEGER PRIMARY KEY | 自動採番 |
| created_at | DATETIME | 記録日時（自動） |
| task_name | TEXT NOT NULL | タスク名 |
| estimate_hours | REAL NOT NULL | 見積もり時間 |
| actual_hours | REAL NOT NULL | 実績時間 |
| ai_minutes | INTEGER | AI活用時間（分） |
| problem | TEXT | ハマったこと |
| solution | TEXT | 解決方法 |
| learning | TEXT | 学び |

## インストール

### 現在（自分用）

```bash
go install github.com/[username]/djou@latest
```

### 将来（組織展開）

```bash
brew tap [username]/tools
brew install djou
```

## 技術スタック

- Go
- SQLite（github.com/mattn/go-sqlite3 または modernc.org/sqlite）
- 対話UI（github.com/AlecAivazis/survey/v2 または github.com/charmbracelet/huh）

## 出力例

### djou list

```
┌────────────┬──────────────────┬──────────┬────────┬─────────┐
│ 日付       │ タスク           │ 見積もり │ 実績   │ AI活用  │
├────────────┼──────────────────┼──────────┼────────┼─────────┤
│ 2025-01-10 │ ログイン画面実装 │ 2.0h     │ 3.0h   │ 30min   │
│ 2025-01-10 │ API連携          │ 1.5h     │ 1.5h   │ 20min   │
│ 2025-01-09 │ 環境構築         │ 1.0h     │ 2.0h   │ 45min   │
└────────────┴──────────────────┴──────────┴────────┴─────────┘
```

### djou stats

```
📊 Dev Journal 統計

期間: 2025-01-01 〜 2025-01-10
記録件数: 15件

⏱️  作業時間
  見積もり合計: 25.0h
  実績合計: 32.5h
  見積もり精度: 77%

🤖 AI活用
  AI活用時間: 6.5h
  AI活用率: 20%

📝 よくハマるポイント TOP3
  1. CORS関連 (3件)
  2. 型エラー (2件)
  3. 環境変数 (2件)
```

### djou export

```bash
$ djou export
✅ djou_2025-01.csv に出力しました
```

出力CSV形式：
```csv
date,task_name,estimate_hours,actual_hours,ai_minutes,problem,solution,learning
2025-01-10,ログイン画面実装,2.0,3.0,30,CORSエラー,プロキシ設定を追加,API連携は早めに確認する
```
