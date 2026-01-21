# Contributing to djou

Thank you for your interest in contributing to djou!

## Code of Conduct

Please be respectful and constructive in all interactions.

## How to Contribute

### Reporting Bugs

1. Check existing issues to avoid duplicates
2. Use the bug report template
3. Include version info (`djou version`)
4. Provide steps to reproduce

### Suggesting Features

1. Check existing issues and discussions
2. Use the feature request template
3. Explain the problem you're trying to solve

### Submitting Code

1. Fork the repository
2. Create a feature branch from `develop`
   ```bash
   git checkout -b feature/your-feature develop
   ```
3. Make your changes
4. Run tests
   ```bash
   go test -v ./...
   ```
5. Run linter
   ```bash
   go vet ./...
   gofmt -w .
   ```
6. Commit with conventional commit format
   ```
   feat: add new feature
   fix: fix bug
   docs: update documentation
   refactor: refactor code
   test: add tests
   chore: maintenance
   ```
7. Push to your fork
8. Create a Pull Request to `develop`

### Branch Naming

- `feature/*` - New features
- `fix/*` - Bug fixes
- `docs/*` - Documentation
- `refactor/*` - Code refactoring

## Development Setup

```bash
# Clone
git clone https://github.com/kca2250/dev-journal.git
cd dev-journal

# Install dependencies
go mod download

# Build
go build -o djou ./cmd/djou

# Test
go test -v ./...

# Run
./djou
```

## Questions?

Use [GitHub Discussions](https://github.com/kca2250/dev-journal/discussions) for questions.
