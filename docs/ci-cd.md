# CI/CD 設計

## 概要

| 項目 | 状態 | 内容 |
|------|------|------|
| CI | 実装する | build, test, lint の自動実行 |
| CD | 後回し | リリース自動化は将来対応 |

---

## CI 要件

### 必須チェック（マージブロック）

以下が全てpassしないとマージ不可：

| チェック | 内容 | 失敗時 |
|----------|------|--------|
| **build** | `go build` が成功する | マージブロック |
| **test** | `go test` が全てpass | マージブロック |
| **lint** | `go vet` がエラーなし | マージブロック |

### オプションチェック

| チェック | 内容 | 失敗時 |
|----------|------|--------|
| coverage | カバレッジ80%以上 | 警告のみ（将来的にブロック） |

---

## GitHub Actions ワークフロー

### トリガー

```yaml
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
```

### ジョブ構成

```
┌─────────────────────────────────────────┐
│              CI Pipeline                │
├─────────────────────────────────────────┤
│  1. checkout                            │
│  2. setup-go                            │
│  3. go mod download                     │
│  4. go vet (lint)                       │
│  5. go build                            │
│  6. go test -race -coverprofile        │
│  7. coverage report (optional)          │
└─────────────────────────────────────────┘
```

---

## ブランチ保護ルール

### main ブランチ

| 設定 | 値 |
|------|-----|
| Require pull request before merging | Yes |
| Require status checks to pass | Yes |
| Required checks | `build`, `test`, `lint` |
| Require branches to be up to date | Yes |

### 設定手順

1. GitHub リポジトリ → Settings → Branches
2. Add branch protection rule
3. Branch name pattern: `main`
4. 上記設定を有効化

---

## ローカルでのCI確認

```bash
# CIと同じチェックをローカルで実行
make ci

# または個別に実行
go vet ./...
go build ./...
go test -race -cover ./...
```

---

## 将来的なCD（リリース自動化）

### 予定

- タグプッシュでリリース作成
- GoReleaserでマルチプラットフォームビルド
- Homebrew tap への自動公開

### 対応時期

- 初期リリース完了後に検討
