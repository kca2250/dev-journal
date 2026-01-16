.PHONY: build test lint fmt ci clean coverage

# ビルド
build:
	go build -v ./...

# テスト実行
test:
	go test -v -race ./...

# カバレッジ付きテスト
coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	@echo ""
	@echo "HTMLレポートを生成するには: go tool cover -html=coverage.out -o coverage.html"

# Lint
lint:
	go vet ./...
	@echo "Checking format..."
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "Code is not formatted:"; \
		gofmt -l .; \
		exit 1; \
	fi
	@echo "Lint passed!"

# フォーマット
fmt:
	go fmt ./...

# CI（ローカルで全チェック実行）
ci: lint build test
	@echo ""
	@echo "✅ All CI checks passed!"

# クリーンアップ
clean:
	rm -f coverage.out coverage.html
	rm -f djou

# 開発用ビルド
dev:
	go build -o djou ./cmd/djou
