# 集計機能 (djou stats)

記録した開発日誌の統計を表示する機能。

## コマンド

```bash
djou stats             # 全体の集計
djou stats --month     # 月別集計
djou stats --month 2025-01  # 特定月の集計
```

## 機能要件

### FR-1: 基本統計の表示

以下の統計情報を計算・表示する。

| 項目 | 計算方法 |
|------|----------|
| 記録件数 | COUNT(*) |
| 見積もり合計 | SUM(estimate_hours) |
| 実績合計 | SUM(actual_hours) |
| 見積もり精度 | (見積もり合計 / 実績合計) × 100 |

**ゼロ除算の扱い**: 実績時間が0の場合、見積もり精度は **0%** として表示

### FR-2: 期間表示

- 集計対象期間（最古の記録日〜最新の記録日）を表示

### FR-3: 月別集計 (--month)

- 月ごとの統計をテーブル形式で表示
- **特定月の指定**: `--month 2025-01` のように特定の月を指定して、その月のみの統計を表示可能

## 出力例

### 全体集計

```
統計情報

期間: 2025-01-01 〜 2025-01-10
総ログ数: 15件

合計時間: 見積 25.0h / 実績 32.5h
見積精度: 77%
```

### 月別集計

```
月別統計

┌─────────┬───────┬──────┬──────┬────────┐
│ 月      │ 件数  │ 見積 │ 実績 │ 精度   │
├─────────┼───────┼──────┼──────┼────────┤
│ 2025-01 │ 15    │ 25.0h│ 32.5h│ 77%    │
│ 2024-12 │ 20    │ 30.0h│ 28.0h│ 107%   │
└─────────┴───────┴──────┴──────┴────────┘
```

## 技術仕様

### CLIオプション

| フラグ | 短縮形 | 説明 | デフォルト |
|--------|--------|------|------------|
| --month | -m | 月別集計を表示。値を指定すると特定月のみ表示（例: `--month 2025-01`） | - |

### データベースクエリ

```sql
-- 全体統計
SELECT
  COUNT(*) as count,
  MIN(created_at) as min_date,
  MAX(created_at) as max_date,
  SUM(estimate_hours) as total_estimate,
  SUM(actual_hours) as total_actual
FROM logs;

-- 月別統計
SELECT
  strftime('%Y-%m', created_at) as month,
  COUNT(*) as count,
  SUM(estimate_hours) as total_estimate,
  SUM(actual_hours) as total_actual
FROM logs
GROUP BY strftime('%Y-%m', created_at)
ORDER BY month DESC;
```

### エラーハンドリング

| エラー | 対応 |
|--------|------|
| DB接続失敗 | エラーメッセージを表示して終了（exit 1） |
| 記録なし | 「記録がありません」と表示 |
| 不正な月指定 | エラーメッセージを表示して終了（exit 1） |
