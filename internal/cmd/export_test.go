package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kca2250/djou/internal/model"
)

func TestExportCmd_Use(t *testing.T) {
	if exportCmd.Use != "export" {
		t.Errorf("exportCmd.Use = %q, want %q", exportCmd.Use, "export")
	}
}

func TestExportCmd_Flags(t *testing.T) {
	// Check --output flag exists
	outputFlag := exportCmd.Flags().Lookup("output")
	if outputFlag == nil {
		t.Error("--output flag should exist")
	}
	if outputFlag.Shorthand != "o" {
		t.Errorf("--output shorthand = %q, want %q", outputFlag.Shorthand, "o")
	}

	// Check --from flag exists
	fromFlag := exportCmd.Flags().Lookup("from")
	if fromFlag == nil {
		t.Error("--from flag should exist")
	}

	// Check --to flag exists
	toFlag := exportCmd.Flags().Lookup("to")
	if toFlag == nil {
		t.Error("--to flag should exist")
	}
}

func TestGenerateFilename(t *testing.T) {
	// Use a temp directory for tests
	tmpDir := t.TempDir()

	tests := []struct {
		name          string
		existingFiles []string
		wantSuffix    string
	}{
		{
			name:          "no existing file",
			existingFiles: []string{},
			wantSuffix:    ".csv",
		},
		{
			name:          "one existing file",
			existingFiles: []string{"djou_2025-01.csv"},
			wantSuffix:    "_1.csv",
		},
		{
			name:          "two existing files",
			existingFiles: []string{"djou_2025-01.csv", "djou_2025-01_1.csv"},
			wantSuffix:    "_2.csv",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a subdirectory for each test case
			testDir := filepath.Join(tmpDir, tt.name)
			if err := os.MkdirAll(testDir, 0755); err != nil {
				t.Fatal(err)
			}

			// Create existing files
			for _, f := range tt.existingFiles {
				path := filepath.Join(testDir, f)
				if err := os.WriteFile(path, []byte{}, 0644); err != nil {
					t.Fatal(err)
				}
			}

			got := generateFilename(testDir, "2025-01")
			if !strings.HasSuffix(got, tt.wantSuffix) {
				t.Errorf("generateFilename() = %q, want suffix %q", got, tt.wantSuffix)
			}
		})
	}
}

func TestValidateDateFormat(t *testing.T) {
	tests := []struct {
		name    string
		date    string
		wantErr bool
	}{
		{
			name:    "valid date",
			date:    "2025-01-15",
			wantErr: false,
		},
		{
			name:    "valid date end of month",
			date:    "2025-12-31",
			wantErr: false,
		},
		{
			name:    "invalid format - no dash",
			date:    "20250115",
			wantErr: true,
		},
		{
			name:    "invalid format - wrong separator",
			date:    "2025/01/15",
			wantErr: true,
		},
		{
			name:    "invalid month",
			date:    "2025-13-01",
			wantErr: true,
		},
		{
			name:    "invalid day",
			date:    "2025-01-32",
			wantErr: true,
		},
		{
			name:    "empty string",
			date:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseDateFlag(tt.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDateFlag(%q) error = %v, wantErr %v", tt.date, err, tt.wantErr)
			}
		})
	}
}

func TestFormatCSVRow(t *testing.T) {
	memo := "CORSでハマった"
	tags := "新機能,フロント"

	log := model.Log{
		ID:            1,
		CreatedAt:     time.Date(2025, 1, 10, 10, 0, 0, 0, time.Local),
		TaskName:      "ログイン実装",
		EstimateHours: 2.0,
		ActualHours:   3.0,
		Memo:          &memo,
		Tags:          &tags,
	}

	row := formatCSVRow(log)

	if len(row) != 6 {
		t.Errorf("formatCSVRow() returned %d fields, want 6", len(row))
	}

	if row[0] != "2025-01-10" {
		t.Errorf("date = %q, want %q", row[0], "2025-01-10")
	}
	if row[1] != "ログイン実装" {
		t.Errorf("task_name = %q, want %q", row[1], "ログイン実装")
	}
	if row[2] != "2.0" {
		t.Errorf("estimate_hours = %q, want %q", row[2], "2.0")
	}
	if row[3] != "3.0" {
		t.Errorf("actual_hours = %q, want %q", row[3], "3.0")
	}
	if row[4] != "CORSでハマった" {
		t.Errorf("memo = %q, want %q", row[4], "CORSでハマった")
	}
	if row[5] != "新機能,フロント" {
		t.Errorf("tags = %q, want %q", row[5], "新機能,フロント")
	}
}

func TestFormatCSVRow_NilFields(t *testing.T) {
	log := model.Log{
		ID:            1,
		CreatedAt:     time.Date(2025, 1, 10, 10, 0, 0, 0, time.Local),
		TaskName:      "API連携",
		EstimateHours: 1.5,
		ActualHours:   1.5,
		Memo:          nil,
		Tags:          nil,
	}

	row := formatCSVRow(log)

	if row[4] != "" {
		t.Errorf("memo for nil = %q, want empty string", row[4])
	}
	if row[5] != "" {
		t.Errorf("tags for nil = %q, want empty string", row[5])
	}
}

func TestGetCSVHeaders(t *testing.T) {
	headers := getCSVHeaders()

	expected := []string{
		"date",
		"task_name",
		"estimate_hours",
		"actual_hours",
		"memo",
		"tags",
	}

	if len(headers) != len(expected) {
		t.Errorf("getCSVHeaders() returned %d headers, want %d", len(headers), len(expected))
	}

	for i, h := range expected {
		if headers[i] != h {
			t.Errorf("headers[%d] = %q, want %q", i, headers[i], h)
		}
	}
}

func TestWriteCSVWithBOM(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.csv")

	logs := []model.Log{
		{
			CreatedAt:     time.Date(2025, 1, 10, 10, 0, 0, 0, time.Local),
			TaskName:      "テスト",
			EstimateHours: 1.0,
			ActualHours:   1.0,
		},
	}

	err := writeCSV(filePath, logs)
	if err != nil {
		t.Fatalf("writeCSV() error = %v", err)
	}

	// Read file and check BOM
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	// UTF-8 BOM: 0xEF, 0xBB, 0xBF
	if len(content) < 3 || content[0] != 0xEF || content[1] != 0xBB || content[2] != 0xBF {
		t.Error("CSV file should start with UTF-8 BOM")
	}

	// Check content contains header and data
	contentStr := string(content[3:]) // Skip BOM
	if !strings.Contains(contentStr, "date,task_name") {
		t.Error("CSV should contain headers")
	}
	if !strings.Contains(contentStr, "テスト") {
		t.Error("CSV should contain task name")
	}
}
