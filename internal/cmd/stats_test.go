package cmd

import (
	"bytes"
	"testing"

	"github.com/kca2250/djou/internal/db"
	"github.com/kca2250/djou/internal/ui"
)

func TestStatsCmd_Use(t *testing.T) {
	if statsCmd.Use != "stats" {
		t.Errorf("statsCmd.Use = %q, want %q", statsCmd.Use, "stats")
	}
}

func TestStatsCmd_Flags(t *testing.T) {
	// Check --month flag exists
	monthFlag := statsCmd.Flags().Lookup("month")
	if monthFlag == nil {
		t.Error("--month flag should exist")
	}
	if monthFlag.Shorthand != "m" {
		t.Errorf("--month shorthand = %q, want %q", monthFlag.Shorthand, "m")
	}
}

func TestFormatAccuracy(t *testing.T) {
	tests := []struct {
		name     string
		estimate float64
		actual   float64
		want     string
	}{
		{
			name:     "normal accuracy",
			estimate: 25.0,
			actual:   32.5,
			want:     "77%",
		},
		{
			name:     "over 100%",
			estimate: 30.0,
			actual:   28.0,
			want:     "107%",
		},
		{
			name:     "zero actual",
			estimate: 10.0,
			actual:   0,
			want:     "0%",
		},
		{
			name:     "both zero",
			estimate: 0,
			actual:   0,
			want:     "0%",
		},
		{
			name:     "exact match",
			estimate: 10.0,
			actual:   10.0,
			want:     "100%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatAccuracy(tt.estimate, tt.actual)
			if got != tt.want {
				t.Errorf("formatAccuracy(%v, %v) = %q, want %q", tt.estimate, tt.actual, got, tt.want)
			}
		})
	}
}

func TestFormatAIRate(t *testing.T) {
	tests := []struct {
		name      string
		aiMinutes int
		actual    float64
		want      string
	}{
		{
			name:      "normal rate",
			aiMinutes: 60,
			actual:    5.0,
			want:      "20%",
		},
		{
			name:      "zero actual",
			aiMinutes: 30,
			actual:    0,
			want:      "0%",
		},
		{
			name:      "zero ai minutes",
			aiMinutes: 0,
			actual:    10.0,
			want:      "0%",
		},
		{
			name:      "both zero",
			aiMinutes: 0,
			actual:    0,
			want:      "0%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatAIRate(tt.aiMinutes, tt.actual)
			if got != tt.want {
				t.Errorf("formatAIRate(%v, %v) = %q, want %q", tt.aiMinutes, tt.actual, got, tt.want)
			}
		})
	}
}

func TestRenderMonthlyStatsTable(t *testing.T) {
	stats := []db.MonthlyStats{
		{Month: "2025-01", Count: 15, TotalEstimate: 25.0, TotalActual: 32.5},
		{Month: "2024-12", Count: 20, TotalEstimate: 30.0, TotalActual: 28.0},
	}

	localizer := ui.NewLocalizer("ja")
	var buf bytes.Buffer

	renderMonthlyStatsTable(&buf, stats, localizer)
	output := buf.String()

	// Verify output contains expected data
	if !bytes.Contains([]byte(output), []byte("2025-01")) {
		t.Error("output should contain '2025-01'")
	}
	if !bytes.Contains([]byte(output), []byte("2024-12")) {
		t.Error("output should contain '2024-12'")
	}
	if !bytes.Contains([]byte(output), []byte("15")) {
		t.Error("output should contain count '15'")
	}
	if !bytes.Contains([]byte(output), []byte("25.0h")) {
		t.Error("output should contain '25.0h'")
	}
	if !bytes.Contains([]byte(output), []byte("77%")) {
		t.Error("output should contain accuracy '77%'")
	}
	if !bytes.Contains([]byte(output), []byte("107%")) {
		t.Error("output should contain accuracy '107%'")
	}
}

func TestRenderMonthlyStatsTable_Empty(t *testing.T) {
	stats := []db.MonthlyStats{}

	localizer := ui.NewLocalizer("ja")
	var buf bytes.Buffer

	renderMonthlyStatsTable(&buf, stats, localizer)
	output := buf.String()

	// Should show no data message
	if !bytes.Contains([]byte(output), []byte(localizer.Get(ui.MsgStatsNoData))) {
		t.Error("output should contain no data message")
	}
}

func TestValidateMonthFormat(t *testing.T) {
	tests := []struct {
		name    string
		month   string
		wantErr bool
	}{
		{
			name:    "valid format",
			month:   "2025-01",
			wantErr: false,
		},
		{
			name:    "valid format december",
			month:   "2024-12",
			wantErr: false,
		},
		{
			name:    "invalid format - no dash",
			month:   "202501",
			wantErr: true,
		},
		{
			name:    "invalid format - wrong separator",
			month:   "2025/01",
			wantErr: true,
		},
		{
			name:    "invalid month - 13",
			month:   "2025-13",
			wantErr: true,
		},
		{
			name:    "invalid month - 0",
			month:   "2025-00",
			wantErr: true,
		},
		{
			name:    "empty string",
			month:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMonthFormat(tt.month)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateMonthFormat(%q) error = %v, wantErr %v", tt.month, err, tt.wantErr)
			}
		})
	}
}
