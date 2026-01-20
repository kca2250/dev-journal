# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-01-21

### Added

- **Core Commands**
  - `djou record` - インタラクティブフォームで開発記録を作成
  - `djou record -q` - クイック記録モード（ワンライナー）
  - `djou list` - 記録一覧表示（インタラクティブモード対応）
  - `djou search` - キーワード検索（`--detail` オプション対応）
  - `djou stats` - 集計・統計機能
  - `djou export` - CSV/JSON出力機能
  - `djou config` - 設定管理
  - `djou version` - バージョン情報表示

- **MCP Server**
  - `djou mcp` - Model Context Protocol サーバー機能
  - 記録の作成・取得・更新・削除ツール

- **Features**
  - タグ機能（記録にタグを付与）
  - フィルタリング（日付、タグ、キーワード）
  - 多言語対応（日本語/英語）

- **Infrastructure**
  - GoReleaser によるマルチプラットフォームビルド
  - GitHub Actions CI/CD
  - SQLite データベース

[0.1.0]: https://github.com/kca2250/dev-journal/releases/tag/v0.1.0
