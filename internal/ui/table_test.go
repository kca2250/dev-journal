package ui

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderTable(t *testing.T) {
	tests := []struct {
		name     string
		headers  []string
		rows     [][]string
		wantRows int // expected number of lines in output
	}{
		{
			name:     "simple table",
			headers:  []string{"Name", "Age"},
			rows:     [][]string{{"Alice", "30"}, {"Bob", "25"}},
			wantRows: 6, // top border + header + separator + 2 data rows + bottom border
		},
		{
			name:     "empty table",
			headers:  []string{"Name", "Age"},
			rows:     [][]string{},
			wantRows: 4, // top border + header + separator + bottom border
		},
		{
			name:     "japanese characters",
			headers:  []string{"タスク", "時間"},
			rows:     [][]string{{"ログイン実装", "2.0h"}, {"API連携", "1.5h"}},
			wantRows: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			RenderTable(&buf, tt.headers, tt.rows)
			output := buf.String()

			lines := strings.Split(strings.TrimSpace(output), "\n")
			if len(lines) != tt.wantRows {
				t.Errorf("RenderTable() got %d lines, want %d\nOutput:\n%s",
					len(lines), tt.wantRows, output)
			}
		})
	}
}

func TestRenderTable_ContainsData(t *testing.T) {
	var buf bytes.Buffer
	headers := []string{"Task", "Hours"}
	rows := [][]string{{"Login", "2.0"}, {"API", "1.5"}}

	RenderTable(&buf, headers, rows)
	output := buf.String()

	// Check that headers are present
	if !strings.Contains(output, "Task") {
		t.Error("output should contain 'Task'")
	}
	if !strings.Contains(output, "Hours") {
		t.Error("output should contain 'Hours'")
	}

	// Check that data is present
	if !strings.Contains(output, "Login") {
		t.Error("output should contain 'Login'")
	}
	if !strings.Contains(output, "2.0") {
		t.Error("output should contain '2.0'")
	}
}

func TestFormatAIMinutes(t *testing.T) {
	tests := []struct {
		name     string
		minutes  *int
		expected string
	}{
		{
			name:     "nil returns dash",
			minutes:  nil,
			expected: "-",
		},
		{
			name:     "zero returns dash",
			minutes:  intPtr(0),
			expected: "-",
		},
		{
			name:     "positive value returns Xmin",
			minutes:  intPtr(30),
			expected: "30min",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatAIMinutes(tt.minutes)
			if result != tt.expected {
				t.Errorf("FormatAIMinutes() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormatHours(t *testing.T) {
	tests := []struct {
		name     string
		hours    float64
		expected string
	}{
		{
			name:     "whole number",
			hours:    2.0,
			expected: "2.0h",
		},
		{
			name:     "decimal",
			hours:    1.5,
			expected: "1.5h",
		},
		{
			name:     "zero",
			hours:    0,
			expected: "0.0h",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatHours(tt.hours)
			if result != tt.expected {
				t.Errorf("FormatHours() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		maxLen   int
		expected string
	}{
		{
			name:     "short string unchanged",
			s:        "hello",
			maxLen:   10,
			expected: "hello",
		},
		{
			name:     "exact length unchanged",
			s:        "hello",
			maxLen:   5,
			expected: "hello",
		},
		{
			name:     "long string truncated",
			s:        "hello world",
			maxLen:   8,
			expected: "hello...",
		},
		{
			name:     "japanese string truncated",
			s:        "ログイン画面実装テスト",
			maxLen:   10,
			expected: "ログイン画面実...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TruncateString(tt.s, tt.maxLen)
			if result != tt.expected {
				t.Errorf("TruncateString() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func intPtr(i int) *int {
	return &i
}
