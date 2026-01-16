package ui

import (
	"bytes"
	"testing"
)

func TestColorWriter_WriteColor(t *testing.T) {
	tests := []struct {
		name     string
		colorize bool
		color    string
		text     string
		want     string
	}{
		{
			name:     "colorize enabled",
			colorize: true,
			color:    Red,
			text:     "error",
			want:     Red + "error" + Reset,
		},
		{
			name:     "colorize disabled",
			colorize: false,
			color:    Red,
			text:     "error",
			want:     "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			cw := NewColorWriter(&buf, tt.colorize)
			cw.WriteColor(tt.color, tt.text)

			if got := buf.String(); got != tt.want {
				t.Errorf("WriteColor() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestColorWriter_WriteSuccess(t *testing.T) {
	var buf bytes.Buffer
	cw := NewColorWriter(&buf, true)
	cw.WriteSuccess("success")

	want := Green + "success" + Reset
	if got := buf.String(); got != want {
		t.Errorf("WriteSuccess() = %q, want %q", got, want)
	}
}

func TestColorWriter_WriteError(t *testing.T) {
	var buf bytes.Buffer
	cw := NewColorWriter(&buf, true)
	cw.WriteError("error")

	want := Red + "error" + Reset
	if got := buf.String(); got != want {
		t.Errorf("WriteError() = %q, want %q", got, want)
	}
}

func TestColorWriter_WriteWarning(t *testing.T) {
	var buf bytes.Buffer
	cw := NewColorWriter(&buf, true)
	cw.WriteWarning("warning")

	want := Yellow + "warning" + Reset
	if got := buf.String(); got != want {
		t.Errorf("WriteWarning() = %q, want %q", got, want)
	}
}

func TestColorWriter_WriteInfo(t *testing.T) {
	var buf bytes.Buffer
	cw := NewColorWriter(&buf, true)
	cw.WriteInfo("info")

	want := Cyan + "info" + Reset
	if got := buf.String(); got != want {
		t.Errorf("WriteInfo() = %q, want %q", got, want)
	}
}

func TestColorWriter_WriteBold(t *testing.T) {
	var buf bytes.Buffer
	cw := NewColorWriter(&buf, true)
	cw.WriteBold("bold")

	want := Bold + "bold" + Reset
	if got := buf.String(); got != want {
		t.Errorf("WriteBold() = %q, want %q", got, want)
	}
}

func TestColorWriter_Colorize(t *testing.T) {
	var buf bytes.Buffer

	cw := NewColorWriter(&buf, true)
	if !cw.Colorize() {
		t.Error("Colorize() should return true")
	}

	cw = NewColorWriter(&buf, false)
	if cw.Colorize() {
		t.Error("Colorize() should return false")
	}
}

func TestColorize(t *testing.T) {
	tests := []struct {
		name     string
		colorize bool
		color    string
		text     string
		want     string
	}{
		{
			name:     "colorize enabled",
			colorize: true,
			color:    Blue,
			text:     "hello",
			want:     Blue + "hello" + Reset,
		},
		{
			name:     "colorize disabled",
			colorize: false,
			color:    Blue,
			text:     "hello",
			want:     "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Colorize(tt.colorize, tt.color, tt.text)
			if got != tt.want {
				t.Errorf("Colorize() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIsTTY(t *testing.T) {
	// bytes.Buffer is not a TTY
	var buf bytes.Buffer
	if IsTTY(&buf) {
		t.Error("IsTTY() should return false for bytes.Buffer")
	}
}

func TestNewAutoColorWriter(t *testing.T) {
	var buf bytes.Buffer
	cw := NewAutoColorWriter(&buf)
	// bytes.Buffer is not a TTY, so colorize should be false
	if cw.Colorize() {
		t.Error("NewAutoColorWriter() should detect non-TTY and disable colorize")
	}
}
