# CSV出力機能 (djou export)

記録した開発日誌をCSVファイルに出力する機能。

## コマンド

```bash
djou export                              # カレントディレクトリに出力
djou export --output ./path              # 出力先を指定
djou export --from 2025-01-01 --to 2025-01-31  # 期間を指定
```

## 機能要件

### FR-1: CSV形式での出力

全ての記録をCSV形式でファイルに出力する。

### FR-1.1: 期間指定フィルター

- `--from` / `--to` オプションで出力期間を指定可能
- 指定しない場合は全件出力

### FR-2: ファイル名の自動生成

- デフォルトファイル名: `djou_YYYY-MM.csv`
- 出力時の年月を使用
- **同名ファイルが存在する場合**: 連番を付与（`djou_2025-01_1.csv`, `djou_2025-01_2.csv` ...）

### FR-3: 出力先の指定

- `--output` オプションで出力先ディレクトリを指定可能
- デフォルトはカレントディレクトリ
- 設定ファイル（`~/.djou/config.yaml`）の `export.output_dir` でデフォルト出力先を変更可能
- 環境変数 `DJOU_EXPORT_DIR` でもオーバーライド可能
- 優先順位: `--output` フラグ > 環境変数 > 設定ファイル > カレントディレクトリ

### FR-4: CSVフォーマット

| カラム | 説明 |
|--------|------|
| date | 記録日 (YYYY-MM-DD) |
| task_name | タスク名 |
| estimate_hours | 見積もり時間 |
| actual_hours | 実績時間 |
| memo | メモ |
| tags | タグ |

### FR-5: 文字コード

- UTF-8（BOM付き）で出力
- Excelでの日本語文字化けを防ぐ

### FR-6: エスケープ処理

- カンマ、改行、ダブルクォートを含む値は適切にエスケープ

## 出力例

### コマンド実行

```bash
$ djou export
✅ djou_2025-01.csv に出力しました

$ djou export --output ./reports
✅ ./reports/djou_2025-01.csv に出力しました
```

### CSV内容

```csv
date,task_name,estimate_hours,actual_hours,memo,tags
2025-01-10,ログイン画面実装,2.0,3.0,CORSエラーでハマった,task_type:新機能
2025-01-10,API連携,1.5,1.5,,
2025-01-09,環境構築,1.0,2.0,Docker起動しない。ポート競合を解消,project:案件A
```

## 技術仕様

### CLIオプション

| フラグ | 短縮形 | 説明 | デフォルト |
|--------|--------|------|------------|
| --output | -o | 出力先ディレクトリ | カレントディレクトリ（設定ファイルで変更可） |
| --from | - | 出力期間の開始日（YYYY-MM-DD） | - |
| --to | - | 出力期間の終了日（YYYY-MM-DD） | - |

### データベースクエリ

```sql
SELECT
  date(created_at) as date,
  task_name,
  estimate_hours,
  actual_hours,
  memo,
  tags
FROM logs
ORDER BY created_at ASC;
```

### エラーハンドリング

| エラー | 対応 |
|--------|------|
| DB接続失敗 | エラーメッセージを表示して終了（exit 1） |
| 出力先ディレクトリが存在しない | エラーメッセージを表示して終了（exit 1） |
| ファイル書き込み権限なし | エラーメッセージを表示して終了（exit 1） |
| 記録なし | 「出力する記録がありません」と表示 |
| 不正な日付指定 | エラーメッセージを表示して終了（exit 1） |
