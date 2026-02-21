# MCP Server 機能 (djou mcp)

Claude Code などの MCP クライアントから djou の機能を利用できるようにする。

## 概要

| 項目 | 内容 |
|------|------|
| 目的 | AI（Claude Code等）が djou を直接ツールとして利用可能にする |
| プロトコル | Model Context Protocol (MCP) |
| トランスポート | stdio（標準入出力） |
| ライブラリ | github.com/mark3labs/mcp-go |

## コマンド

```bash
djou mcp    # MCP サーバーとして起動（stdio）
```

## 機能要件

### FR-1: MCP サーバーとしての起動

- `djou mcp` コマンドで MCP サーバーモードとして起動
- stdio トランスポートで MCP クライアントと通信
- JSON-RPC 2.0 プロトコルに準拠

### FR-2: 提供する Tools

既存の CLI 機能を MCP Tools として提供する。

| Tool名 | 説明 | 対応するCLI |
|--------|------|-------------|
| `djou_record` | 開発日誌を記録する | `djou` |
| `djou_list` | 記録一覧を取得する | `djou list` |
| `djou_search` | キーワードで検索する | `djou search` |
| `djou_stats` | 統計情報を取得する | `djou stats` |
| `djou_export` | CSV出力する | `djou export` |
| `djou_update` | 記録を更新する | （対話モードの編集に対応） |
| `djou_delete` | 記録を削除する | （対話モードの削除に対応） |

### FR-3: 各 Tool の仕様

#### djou_record

開発日誌を記録する。

**入力パラメータ:**

| パラメータ | 型 | 必須 | 説明 |
|------------|-----|------|------|
| task_name | string | ⭕ | タスク名 |
| estimate_hours | number | ⭕ | 見積もり時間（時間） |
| actual_hours | number | ⭕ | 実績時間（時間） |
| memo | string | ❌ | メモ |
| tags | string | ❌ | タグ（カンマ区切り） |

**出力:**
```json
{
  "content": [
    {
      "type": "text",
      "text": "✅ 記録しました: ログイン画面実装"
    }
  ]
}
```

#### djou_list

記録一覧を取得する。

**入力パラメータ:**

| パラメータ | 型 | 必須 | 説明 |
|------------|-----|------|------|
| filter | string | ❌ | "week" または "month"、指定なしで直近 |
| limit | integer | ❌ | 表示件数、デフォルト: 10 |

**出力:**
```json
{
  "content": [
    {
      "type": "text",
      "text": "| 日付 | タスク | 見積 | 実績 |\n|---|---|---|---|\n| 2025-01-10 | ログイン画面実装 | 2.0h | 3.0h |"
    }
  ]
}
```

#### djou_search

キーワードで検索する。

**入力パラメータ:**

| パラメータ | 型 | 必須 | 説明 |
|------------|-----|------|------|
| keywords | string | ⭕ | 検索キーワード（スペース区切りでAND検索） |
| limit | integer | ❌ | 表示件数、デフォルト: 無制限 |

**出力:**
```json
{
  "content": [
    {
      "type": "text",
      "text": "🔍 検索結果: 2件\n\n| 日付 | タスク | マッチ箇所 |\n|---|---|---|\n| 2025-01-10 | ログイン画面実装 | メモ |"
    }
  ]
}
```

#### djou_stats

統計情報を取得する。

**入力パラメータ:**

| パラメータ | 型 | 必須 | 説明 |
|------------|-----|------|------|
| month | string | ❌ | 特定月を指定（例: "2025-01"）、指定なしで全体統計 |
| monthly_list | boolean | ❌ | true で月別一覧を表示 |

**出力:**
```json
{
  "content": [
    {
      "type": "text",
      "text": "統計情報\n\n期間: 2025-01-01 〜 2025-01-10\n総ログ数: 15件\n\n合計時間: 見積 25.0h / 実績 32.5h\n見積精度: 77%"
    }
  ]
}
```

#### djou_export

CSV出力する。

**入力パラメータ:**

| パラメータ | 型 | 必須 | 説明 |
|------------|-----|------|------|
| output | string | ❌ | 出力先ディレクトリ、デフォルト: カレントディレクトリ |
| from | string | ❌ | 期間開始日（YYYY-MM-DD） |
| to | string | ❌ | 期間終了日（YYYY-MM-DD） |

**出力:**
```json
{
  "content": [
    {
      "type": "text",
      "text": "✅ djou_2025-01.csv に出力しました"
    }
  ]
}
```

#### djou_update

既存の記録を更新する。

**入力パラメータ:**

| パラメータ | 型 | 必須 | 説明 |
|------------|-----|------|------|
| id | integer | ⭕ | 更新対象のログID |
| task_name | string | ❌ | タスク名（指定時のみ更新） |
| estimate_hours | number | ❌ | 見積もり時間（指定時のみ更新） |
| actual_hours | number | ❌ | 実績時間（指定時のみ更新） |
| memo | string | ❌ | メモ（指定時のみ更新） |
| tags | string | ❌ | タグ（指定時のみ更新） |

**動作:** 指定されたフィールドのみ更新し、未指定フィールドは既存値を保持する。

**出力:**
```json
{
  "content": [
    {
      "type": "text",
      "text": "✅ ログを更新しました (ID: 1)"
    }
  ]
}
```

#### djou_delete

記録を削除する。

**入力パラメータ:**

| パラメータ | 型 | 必須 | 説明 |
|------------|-----|------|------|
| id | integer | ⭕ | 削除対象のログID |

**出力:**
```json
{
  "content": [
    {
      "type": "text",
      "text": "✅ ログを削除しました (ID: 1)"
    }
  ]
}
```

## Claude Code での利用例

### CLAUDE.md への記述例

```markdown
# 開発ログの記録

作業の区切りごとに djou_record ツールで開発ログを記録してください。

## 記録タイミング
- 機能実装が完了したとき
- 1時間以上の作業が終わったとき
- 問題を解決したとき

## 記録内容
- task_name: 作業内容を簡潔に
- estimate_hours: 事前見積もり（なければ実績と同じ）
- actual_hours: 実際にかかった時間
- memo: ハマったことや学びがあれば記載
- tags: 分類タグ（例: task_type:新機能, project:案件A）
```

### Claude Code からの呼び出し例

```
Claude Code が djou_record を呼び出し:
{
  "task_name": "認証機能の実装",
  "estimate_hours": 2,
  "actual_hours": 2.5,
  "memo": "JWTトークンの有効期限切れハンドリングでハマった。リフレッシュトークンで解決。",
  "tags": "task_type:新機能, project:認証基盤"
}
```

## 技術仕様

### ディレクトリ構成

```
djou/
├── cmd/
│   └── djou/
│       └── main.go          # エントリーポイント
├── internal/
│   ├── cmd/
│   │   ├── root.go          # ルートコマンド（記録）
│   │   ├── record.go        # 記録コマンド（Quick Record）
│   │   ├── list.go          # 一覧コマンド
│   │   ├── search.go        # 検索コマンド
│   │   ├── stats.go         # 統計コマンド
│   │   ├── export.go        # エクスポートコマンド
│   │   ├── config.go        # 設定コマンド
│   │   ├── version.go       # バージョンコマンド
│   │   └── mcp.go           # MCP サブコマンド
│   ├── mcp/                  # MCP サーバー実装
│   │   ├── server.go        # サーバー初期化
│   │   └── tools.go         # Tools 定義・ハンドラー（7ツール）
│   ├── db/
│   │   └── log_repository.go # データアクセス層
│   └── model/
│       └── log.go           # データモデル
└── go.mod
```

### 依存関係

```go
require (
    github.com/mark3labs/mcp-go v1.x.x
    // 既存の依存関係
)
```

### MCP サーバー初期化

```go
// internal/mcp/server.go
package mcp

func NewServer() *server.MCPServer {
    s := server.NewMCPServer(
        "djou",
        "1.0.0",
        server.WithToolCapabilities(false),
    )

    // Tools 登録（7ツール）
    registerRecordTool(s)
    registerListTool(s)
    registerSearchTool(s)
    registerStatsTool(s)
    registerExportTool(s)
    registerUpdateTool(s)
    registerDeleteTool(s)

    return s
}

func Run() error {
    s := NewServer()
    return server.ServeStdio(s)
}
```

### Claude Code 設定例

`.mcp.json` または Claude Code の設定:

```json
{
  "mcpServers": {
    "djou": {
      "command": "djou",
      "args": ["mcp"]
    }
  }
}
```

## エラーハンドリング

| エラー | 対応 |
|--------|------|
| DB接続失敗 | エラーメッセージを返す（MCP エラーレスポンス） |
| 必須パラメータ不足 | バリデーションエラーを返す |
| 不正なパラメータ値 | バリデーションエラーを返す |
| ファイル書き込み失敗（export） | エラーメッセージを返す |
| 指定IDが存在しない（update/delete） | エラーメッセージを返す |

### エラーレスポンス例

```json
{
  "content": [
    {
      "type": "text",
      "text": "❌ エラー: task_name は必須です"
    }
  ],
  "isError": true
}
```

## テスト方針

### 単体テスト

- 各 Tool ハンドラーのテスト
- パラメータバリデーションのテスト
- DB 操作との統合テスト

### 結合テスト

- MCP クライアントからの呼び出しテスト
- JSON-RPC メッセージの送受信テスト

### 手動テスト

1. `djou mcp` でサーバー起動
2. Claude Code から各 Tool を呼び出し
3. 正常系・異常系の動作確認

## 完了条件

- [x] `djou mcp` コマンドで MCP サーバーが起動する
- [x] 7つの Tools が全て動作する
- [x] Claude Code から呼び出して記録・取得ができる
- [x] エラー時に適切なレスポンスを返す
- [x] 単体テストが通る

## 依存

- Phase 1〜7 の全機能（既存の Repository を再利用）
