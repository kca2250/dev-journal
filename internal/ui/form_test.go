package ui

import (
	"testing"
)

func TestParsePositiveFloat(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    float64
		wantErr bool
	}{
		{
			name:    "valid positive float",
			input:   "2.5",
			want:    2.5,
			wantErr: false,
		},
		{
			name:    "valid integer as float",
			input:   "3",
			want:    3.0,
			wantErr: false,
		},
		{
			name:    "valid with whitespace",
			input:   "  1.5  ",
			want:    1.5,
			wantErr: false,
		},
		{
			name:    "zero is invalid",
			input:   "0",
			want:    0,
			wantErr: true,
		},
		{
			name:    "negative is invalid",
			input:   "-1",
			want:    0,
			wantErr: true,
		},
		{
			name:    "empty is invalid",
			input:   "",
			want:    0,
			wantErr: true,
		},
		{
			name:    "non-numeric is invalid",
			input:   "abc",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePositiveFloat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePositiveFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParsePositiveFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}


func TestValidateTaskName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid task name",
			input:   "ログイン実装",
			wantErr: false,
		},
		{
			name:    "empty is invalid",
			input:   "",
			wantErr: true,
		},
		{
			name:    "whitespace only is invalid",
			input:   "   ",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTaskName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTaskName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePositiveFloat(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid",
			input:   "2.5",
			wantErr: false,
		},
		{
			name:    "zero is invalid",
			input:   "0",
			wantErr: true,
		},
		{
			name:    "empty is invalid",
			input:   "",
			wantErr: true,
		},
		{
			name:    "non-numeric is invalid",
			input:   "abc",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePositiveFloat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePositiveFloat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFormInput_ToLogInput(t *testing.T) {
	tests := []struct {
		name       string
		input      FormInput
		wantErr    bool
		wantTask   string
		wantEst    float64
		wantActual float64
		wantMemo   string
	}{
		{
			name: "valid full input",
			input: FormInput{
				TaskName:      "テストタスク",
				EstimateHours: "2.0",
				ActualHours:   "3.0",
				Memo:          "メモ内容",
			},
			wantErr:    false,
			wantTask:   "テストタスク",
			wantEst:    2.0,
			wantActual: 3.0,
			wantMemo:   "メモ内容",
		},
		{
			name: "valid minimal input",
			input: FormInput{
				TaskName:      "タスク",
				EstimateHours: "1",
				ActualHours:   "1.5",
				Memo:          "",
			},
			wantErr:    false,
			wantTask:   "タスク",
			wantEst:    1.0,
			wantActual: 1.5,
			wantMemo:   "",
		},
		{
			name: "trims whitespace",
			input: FormInput{
				TaskName:      "  タスク  ",
				EstimateHours: "  2  ",
				ActualHours:   "  3  ",
				Memo:          "  メモ  ",
			},
			wantErr:    false,
			wantTask:   "タスク",
			wantEst:    2.0,
			wantActual: 3.0,
			wantMemo:   "メモ",
		},
		{
			name: "invalid estimate",
			input: FormInput{
				TaskName:      "タスク",
				EstimateHours: "abc",
				ActualHours:   "1",
				Memo:          "",
			},
			wantErr: true,
		},
		{
			name: "invalid actual",
			input: FormInput{
				TaskName:      "タスク",
				EstimateHours: "1",
				ActualHours:   "abc",
				Memo:          "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.ToLogInput()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToLogInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if got.TaskName != tt.wantTask {
				t.Errorf("TaskName = %q, want %q", got.TaskName, tt.wantTask)
			}
			if got.EstimateHours != tt.wantEst {
				t.Errorf("EstimateHours = %v, want %v", got.EstimateHours, tt.wantEst)
			}
			if got.ActualHours != tt.wantActual {
				t.Errorf("ActualHours = %v, want %v", got.ActualHours, tt.wantActual)
			}
			if got.Memo != tt.wantMemo {
				t.Errorf("Memo = %q, want %q", got.Memo, tt.wantMemo)
			}
		})
	}
}

func TestNewRecordForm(t *testing.T) {
	l := NewLocalizer("ja")
	form := NewRecordForm(l)

	if form == nil {
		t.Fatal("NewRecordForm() returned nil")
	}

	if form.input == nil {
		t.Error("form.input should not be nil")
	}

	if form.localizer != l {
		t.Error("form.localizer should be set")
	}
}
