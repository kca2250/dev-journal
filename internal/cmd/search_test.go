package cmd

import (
	"bytes"
	"testing"
	"time"

	"github.com/kca2250/djou/internal/model"
	"github.com/kca2250/djou/internal/ui"
)

func TestSearchCmd_Use(t *testing.T) {
	if searchCmd.Use != "search [keywords...]" {
		t.Errorf("searchCmd.Use = %q, want %q", searchCmd.Use, "search [keywords...]")
	}
}

func TestSearchCmd_Flags(t *testing.T) {
	// Check --limit flag exists
	limitFlag := searchCmd.Flags().Lookup("limit")
	if limitFlag == nil {
		t.Error("--limit flag should exist")
	}
	if limitFlag.Shorthand != "l" {
		t.Errorf("--limit shorthand = %q, want %q", limitFlag.Shorthand, "l")
	}
	if limitFlag.DefValue != "0" {
		t.Errorf("--limit default = %q, want %q", limitFlag.DefValue, "0")
	}
}

func TestSearchCmd_RequiresArgs(t *testing.T) {
	// The command requires at least 1 argument
	if searchCmd.Args == nil {
		t.Error("searchCmd.Args should not be nil")
	}
}

func TestGetMatchFields(t *testing.T) {
	localizer := ui.NewLocalizer("ja")

	problem := "CORSエラーが発生"
	solution := "プロキシを設定"
	learning := "APIの仕様を確認すべき"

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
			wantCount: 2, // matches task name and learning
		},
		{
			name: "multiple keywords match multiple fields",
			log: model.Log{
				TaskName: "ログイン実装",
				Problem:  &problem,
			},
			keywords:  []string{"ログイン", "CORS"},
			wantCount: 2,
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
			matches := getMatchFields(tt.log, tt.keywords, localizer)
			if len(matches) != tt.wantCount {
				t.Errorf("getMatchFields() returned %d matches, want %d", len(matches), tt.wantCount)
			}
		})
	}
}

func TestContains(t *testing.T) {
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
			if got := contains(tt.slice, tt.item); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRenderSearchTable(t *testing.T) {
	localizer := ui.NewLocalizer("ja")

	problem := "CORSエラー"
	logs := []model.Log{
		{
			ID:        1,
			CreatedAt: time.Date(2025, 1, 10, 10, 0, 0, 0, time.Local),
			TaskName:  "ログイン実装",
			Problem:   &problem,
		},
	}

	var buf bytes.Buffer

	// Test the table rendering logic
	headers := []string{
		localizer.Get(ui.MsgHeaderDate),
		localizer.Get(ui.MsgHeaderTask),
		localizer.Get(ui.MsgHeaderMatchField),
	}

	keywords := []string{"CORS"}
	rows := make([][]string, len(logs))
	for i, log := range logs {
		matchFields := getMatchFields(log, keywords, localizer)
		rows[i] = []string{
			log.CreatedAt.Format("2006-01-02"),
			ui.TruncateString(log.TaskName, 20),
			"",
		}
		if len(matchFields) > 0 {
			rows[i][2] = matchFields[0]
		}
	}

	ui.RenderTable(&buf, headers, rows)
	output := buf.String()

	// Verify output contains expected data
	if !bytes.Contains([]byte(output), []byte("2025-01-10")) {
		t.Error("output should contain date '2025-01-10'")
	}
	if !bytes.Contains([]byte(output), []byte("ログイン実装")) {
		t.Error("output should contain task name")
	}
	if !bytes.Contains([]byte(output), []byte("マッチ箇所")) {
		t.Error("output should contain match field header")
	}
}
