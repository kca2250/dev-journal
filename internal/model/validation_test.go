package model

import (
	"errors"
	"testing"
)

func TestValidateInput(t *testing.T) {
	tests := []struct {
		name    string
		input   LogInput
		wantErr error
	}{
		{
			name: "valid input with all fields",
			input: LogInput{
				TaskName:      "テストタスク",
				EstimateHours: 2.0,
				ActualHours:   3.0,
				AIMinutes:     30,
				Problem:       "問題",
				Solution:      "解決",
				Learning:      "学び",
			},
			wantErr: nil,
		},
		{
			name: "valid input with required fields only",
			input: LogInput{
				TaskName:      "タスク",
				EstimateHours: 1.0,
				ActualHours:   1.0,
			},
			wantErr: nil,
		},
		{
			name: "empty task name",
			input: LogInput{
				TaskName:      "",
				EstimateHours: 1.0,
				ActualHours:   1.0,
			},
			wantErr: ErrEmptyTaskName,
		},
		{
			name: "whitespace only task name",
			input: LogInput{
				TaskName:      "   ",
				EstimateHours: 1.0,
				ActualHours:   1.0,
			},
			wantErr: ErrEmptyTaskName,
		},
		{
			name: "zero estimate hours",
			input: LogInput{
				TaskName:      "タスク",
				EstimateHours: 0,
				ActualHours:   1.0,
			},
			wantErr: ErrInvalidEstimate,
		},
		{
			name: "negative estimate hours",
			input: LogInput{
				TaskName:      "タスク",
				EstimateHours: -1.0,
				ActualHours:   1.0,
			},
			wantErr: ErrInvalidEstimate,
		},
		{
			name: "zero actual hours",
			input: LogInput{
				TaskName:      "タスク",
				EstimateHours: 1.0,
				ActualHours:   0,
			},
			wantErr: ErrInvalidActual,
		},
		{
			name: "negative actual hours",
			input: LogInput{
				TaskName:      "タスク",
				EstimateHours: 1.0,
				ActualHours:   -1.0,
			},
			wantErr: ErrInvalidActual,
		},
		{
			name: "negative AI minutes",
			input: LogInput{
				TaskName:      "タスク",
				EstimateHours: 1.0,
				ActualHours:   1.0,
				AIMinutes:     -10,
			},
			wantErr: ErrInvalidAIMinutes,
		},
		{
			name: "zero AI minutes is valid",
			input: LogInput{
				TaskName:      "タスク",
				EstimateHours: 1.0,
				ActualHours:   1.0,
				AIMinutes:     0,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateInput(&tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidateInput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
