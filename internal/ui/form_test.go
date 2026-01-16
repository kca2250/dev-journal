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

func TestParseNonNegativeInt(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{
			name:    "valid positive int",
			input:   "30",
			want:    30,
			wantErr: false,
		},
		{
			name:    "zero is valid",
			input:   "0",
			want:    0,
			wantErr: false,
		},
		{
			name:    "empty returns zero",
			input:   "",
			want:    0,
			wantErr: false,
		},
		{
			name:    "whitespace returns zero",
			input:   "  ",
			want:    0,
			wantErr: false,
		},
		{
			name:    "valid with whitespace",
			input:   "  15  ",
			want:    15,
			wantErr: false,
		},
		{
			name:    "negative is invalid",
			input:   "-1",
			want:    0,
			wantErr: true,
		},
		{
			name:    "non-numeric is invalid",
			input:   "abc",
			want:    0,
			wantErr: true,
		},
		{
			name:    "float is invalid",
			input:   "1.5",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseNonNegativeInt(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseNonNegativeInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseNonNegativeInt() = %v, want %v", got, tt.want)
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

func TestValidateNonNegativeInt(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid positive",
			input:   "30",
			wantErr: false,
		},
		{
			name:    "zero is valid",
			input:   "0",
			wantErr: false,
		},
		{
			name:    "empty is valid",
			input:   "",
			wantErr: false,
		},
		{
			name:    "negative is invalid",
			input:   "-1",
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
			err := ValidateNonNegativeInt(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNonNegativeInt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFormInput_ToLogInput(t *testing.T) {
	tests := []struct {
		name        string
		input       FormInput
		wantErr     bool
		wantTask    string
		wantEst     float64
		wantActual  float64
		wantAI      int
		wantProblem string
	}{
		{
			name: "valid full input",
			input: FormInput{
				TaskName:      "テストタスク",
				EstimateHours: "2.0",
				ActualHours:   "3.0",
				AIMinutes:     "30",
				Problem:       "問題",
				Solution:      "解決",
				Learning:      "学び",
			},
			wantErr:     false,
			wantTask:    "テストタスク",
			wantEst:     2.0,
			wantActual:  3.0,
			wantAI:      30,
			wantProblem: "問題",
		},
		{
			name: "valid minimal input",
			input: FormInput{
				TaskName:      "タスク",
				EstimateHours: "1",
				ActualHours:   "1.5",
				AIMinutes:     "",
				Problem:       "",
				Solution:      "",
				Learning:      "",
			},
			wantErr:     false,
			wantTask:    "タスク",
			wantEst:     1.0,
			wantActual:  1.5,
			wantAI:      0,
			wantProblem: "",
		},
		{
			name: "trims whitespace",
			input: FormInput{
				TaskName:      "  タスク  ",
				EstimateHours: "  2  ",
				ActualHours:   "  3  ",
				AIMinutes:     "  10  ",
				Problem:       "  問題  ",
				Solution:      "",
				Learning:      "",
			},
			wantErr:     false,
			wantTask:    "タスク",
			wantEst:     2.0,
			wantActual:  3.0,
			wantAI:      10,
			wantProblem: "問題",
		},
		{
			name: "invalid estimate",
			input: FormInput{
				TaskName:      "タスク",
				EstimateHours: "abc",
				ActualHours:   "1",
				AIMinutes:     "",
			},
			wantErr: true,
		},
		{
			name: "invalid actual",
			input: FormInput{
				TaskName:      "タスク",
				EstimateHours: "1",
				ActualHours:   "abc",
				AIMinutes:     "",
			},
			wantErr: true,
		},
		{
			name: "invalid AI minutes",
			input: FormInput{
				TaskName:      "タスク",
				EstimateHours: "1",
				ActualHours:   "1",
				AIMinutes:     "-5",
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
			if got.AIMinutes != tt.wantAI {
				t.Errorf("AIMinutes = %v, want %v", got.AIMinutes, tt.wantAI)
			}
			if got.Problem != tt.wantProblem {
				t.Errorf("Problem = %q, want %q", got.Problem, tt.wantProblem)
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
