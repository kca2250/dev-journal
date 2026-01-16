package mcp

import (
	"testing"

	"github.com/kca2250/djou/internal/model"
)

func TestGetMatchFields(t *testing.T) {
	problem := "CORSエラーが発生"
	solution := "プロキシ設定を追加"
	learning := "API連携は早めに確認"

	tests := []struct {
		name      string
		log       model.Log
		keywords  []string
		wantCount int
	}{
		{
			name: "match in task name",
			log: model.Log{
				TaskName: "ログイン実装",
			},
			keywords:  []string{"ログイン"},
			wantCount: 1,
		},
		{
			name: "match in problem",
			log: model.Log{
				TaskName: "API実装",
				Problem:  &problem,
			},
			keywords:  []string{"CORS"},
			wantCount: 1,
		},
		{
			name: "match in solution",
			log: model.Log{
				TaskName: "API実装",
				Solution: &solution,
			},
			keywords:  []string{"プロキシ"},
			wantCount: 1,
		},
		{
			name: "match in learning",
			log: model.Log{
				TaskName: "API実装",
				Learning: &learning,
			},
			keywords:  []string{"API"},
			wantCount: 2, // matches task_name and learning
		},
		{
			name: "case insensitive",
			log: model.Log{
				TaskName: "Login Implementation",
			},
			keywords:  []string{"login"},
			wantCount: 1,
		},
		{
			name: "no match",
			log: model.Log{
				TaskName: "テスト",
			},
			keywords:  []string{"存在しない"},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := getMatchFields(tt.log, tt.keywords)
			if len(matches) != tt.wantCount {
				t.Errorf("getMatchFields() returned %d matches, want %d", len(matches), tt.wantCount)
			}
		})
	}
}

func TestContainsString(t *testing.T) {
	tests := []struct {
		name  string
		slice []string
		item  string
		want  bool
	}{
		{
			name:  "item exists",
			slice: []string{"a", "b", "c"},
			item:  "b",
			want:  true,
		},
		{
			name:  "item does not exist",
			slice: []string{"a", "b", "c"},
			item:  "d",
			want:  false,
		},
		{
			name:  "empty slice",
			slice: []string{},
			item:  "a",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := containsString(tt.slice, tt.item); got != tt.want {
				t.Errorf("containsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEscapeCSV(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "no escape needed",
			input: "hello",
			want:  "hello",
		},
		{
			name:  "contains comma",
			input: "hello,world",
			want:  "\"hello,world\"",
		},
		{
			name:  "contains quote",
			input: "hello\"world",
			want:  "\"hello\"\"world\"",
		},
		{
			name:  "contains newline",
			input: "hello\nworld",
			want:  "\"hello\nworld\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := escapeCSV(tt.input); got != tt.want {
				t.Errorf("escapeCSV() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGenerateFilename(t *testing.T) {
	tmpDir := t.TempDir()

	// First call should return base name
	got := generateFilename(tmpDir, "2026-01")
	want := tmpDir + "/djou_2026-01.csv"
	if got != want {
		t.Errorf("generateFilename() = %q, want %q", got, want)
	}
}
