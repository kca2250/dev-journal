package cmd

import (
	"bytes"
	"testing"
	"time"

	"github.com/kca2250/djou/internal/model"
	"github.com/kca2250/djou/internal/ui"
)

func TestListCmd_Use(t *testing.T) {
	if listCmd.Use != "list" {
		t.Errorf("listCmd.Use = %q, want %q", listCmd.Use, "list")
	}
}

func TestListCmd_Flags(t *testing.T) {
	// Check --week flag exists
	weekFlag := listCmd.Flags().Lookup("week")
	if weekFlag == nil {
		t.Error("--week flag should exist")
	}
	if weekFlag.Shorthand != "w" {
		t.Errorf("--week shorthand = %q, want %q", weekFlag.Shorthand, "w")
	}

	// Check --month flag exists
	monthFlag := listCmd.Flags().Lookup("month")
	if monthFlag == nil {
		t.Error("--month flag should exist")
	}
	if monthFlag.Shorthand != "m" {
		t.Errorf("--month shorthand = %q, want %q", monthFlag.Shorthand, "m")
	}

	// Check --limit flag exists
	limitFlag := listCmd.Flags().Lookup("limit")
	if limitFlag == nil {
		t.Error("--limit flag should exist")
	}
	if limitFlag.Shorthand != "l" {
		t.Errorf("--limit shorthand = %q, want %q", limitFlag.Shorthand, "l")
	}
	if limitFlag.DefValue != "10" {
		t.Errorf("--limit default = %q, want %q", limitFlag.DefValue, "10")
	}
}

func TestErrExclusiveOptions(t *testing.T) {
	if ErrExclusiveOptions.Error() != "--week and --month cannot be used together" {
		t.Errorf("ErrExclusiveOptions = %q, want %q",
			ErrExclusiveOptions.Error(), "--week and --month cannot be used together")
	}
}

func TestRenderLogTable(t *testing.T) {
	// Create test logs
	aiMinutes := 30
	logs := []model.Log{
		{
			ID:            1,
			CreatedAt:     time.Date(2025, 1, 10, 10, 0, 0, 0, time.Local),
			TaskName:      "ログイン実装",
			EstimateHours: 2.0,
			ActualHours:   3.0,
			AIMinutes:     &aiMinutes,
		},
		{
			ID:            2,
			CreatedAt:     time.Date(2025, 1, 9, 10, 0, 0, 0, time.Local),
			TaskName:      "API連携",
			EstimateHours: 1.5,
			ActualHours:   1.5,
			AIMinutes:     nil,
		},
	}

	localizer := ui.NewLocalizer("ja")
	var buf bytes.Buffer

	// Note: renderLogTable expects *os.File, but for testing we need to
	// verify the table rendering logic works with the logs
	// We'll test the underlying table render function instead
	headers := []string{
		localizer.Get(ui.MsgHeaderDate),
		localizer.Get(ui.MsgHeaderTask),
		localizer.Get(ui.MsgHeaderEstimate),
		localizer.Get(ui.MsgHeaderActual),
		localizer.Get(ui.MsgHeaderAI),
	}

	rows := make([][]string, len(logs))
	for i, log := range logs {
		rows[i] = []string{
			log.CreatedAt.Format("2006-01-02"),
			ui.TruncateString(log.TaskName, 20),
			ui.FormatHours(log.EstimateHours),
			ui.FormatHours(log.ActualHours),
			ui.FormatAIMinutes(log.AIMinutes),
		}
	}

	ui.RenderTable(&buf, headers, rows)
	output := buf.String()

	// Verify output contains expected data
	if !bytes.Contains([]byte(output), []byte("2025-01-10")) {
		t.Error("output should contain date '2025-01-10'")
	}
	if !bytes.Contains([]byte(output), []byte("ログイン実装")) {
		t.Error("output should contain task 'ログイン実装'")
	}
	if !bytes.Contains([]byte(output), []byte("2.0h")) {
		t.Error("output should contain '2.0h'")
	}
	if !bytes.Contains([]byte(output), []byte("30min")) {
		t.Error("output should contain '30min'")
	}
	if !bytes.Contains([]byte(output), []byte("-")) {
		t.Error("output should contain '-' for nil AI minutes")
	}
}

func TestRenderLogTable_LongTaskName(t *testing.T) {
	logs := []model.Log{
		{
			ID:            1,
			CreatedAt:     time.Now(),
			TaskName:      "これはとても長いタスク名で20文字を超えています",
			EstimateHours: 1.0,
			ActualHours:   1.0,
		},
	}

	localizer := ui.NewLocalizer("ja")
	var buf bytes.Buffer

	rows := make([][]string, len(logs))
	for i, log := range logs {
		rows[i] = []string{
			log.CreatedAt.Format("2006-01-02"),
			ui.TruncateString(log.TaskName, 20),
			ui.FormatHours(log.EstimateHours),
			ui.FormatHours(log.ActualHours),
			ui.FormatAIMinutes(log.AIMinutes),
		}
	}

	headers := []string{
		localizer.Get(ui.MsgHeaderDate),
		localizer.Get(ui.MsgHeaderTask),
		localizer.Get(ui.MsgHeaderEstimate),
		localizer.Get(ui.MsgHeaderActual),
		localizer.Get(ui.MsgHeaderAI),
	}

	ui.RenderTable(&buf, headers, rows)
	output := buf.String()

	// Task name should be truncated to 20 characters with "..."
	if bytes.Contains([]byte(output), []byte("これはとても長いタスク名で20文字を超えています")) {
		t.Error("long task name should be truncated")
	}
	if !bytes.Contains([]byte(output), []byte("...")) {
		t.Error("truncated task name should contain '...'")
	}
}
