package ui

import (
	"os"
	"testing"
)

func TestNewLocalizer(t *testing.T) {
	tests := []struct {
		name     string
		lang     string
		wantLang string
	}{
		{
			name:     "Japanese",
			lang:     "ja",
			wantLang: "ja",
		},
		{
			name:     "English",
			lang:     "en",
			wantLang: "en",
		},
		{
			name:     "Japanese with encoding",
			lang:     "ja_JP.UTF-8",
			wantLang: "ja",
		},
		{
			name:     "English with encoding",
			lang:     "en_US.UTF-8",
			wantLang: "en",
		},
		{
			name:     "Unsupported language defaults to English",
			lang:     "fr",
			wantLang: "en",
		},
		{
			name:     "Empty string defaults to English",
			lang:     "",
			wantLang: "en",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLocalizer(tt.lang)
			if l.Lang() != tt.wantLang {
				t.Errorf("NewLocalizer(%q).Lang() = %q, want %q", tt.lang, l.Lang(), tt.wantLang)
			}
		})
	}
}

func TestLocalizer_Get(t *testing.T) {
	tests := []struct {
		name string
		lang string
		key  string
		want string
	}{
		{
			name: "Japanese message",
			lang: "ja",
			key:  MsgRecordSuccess,
			want: "ログを記録しました",
		},
		{
			name: "English message",
			lang: "en",
			key:  MsgRecordSuccess,
			want: "Log recorded successfully",
		},
		{
			name: "Unknown key returns key",
			lang: "ja",
			key:  "unknown_key",
			want: "unknown_key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLocalizer(tt.lang)
			got := l.Get(tt.key)
			if got != tt.want {
				t.Errorf("Get(%q) = %q, want %q", tt.key, got, tt.want)
			}
		})
	}
}

func TestLocalizer_Getf(t *testing.T) {
	tests := []struct {
		name string
		lang string
		key  string
		args []any
		want string
	}{
		{
			name: "Japanese formatted message",
			lang: "ja",
			key:  MsgStatsTotalLogs,
			args: []any{10},
			want: "総ログ数: 10件",
		},
		{
			name: "English formatted message",
			lang: "en",
			key:  MsgStatsTotalLogs,
			args: []any{10},
			want: "Total logs: 10",
		},
		{
			name: "Multiple arguments",
			lang: "ja",
			key:  MsgStatsTotalHours,
			args: []any{10.5, 12.0},
			want: "合計時間: 見積 10.5h / 実績 12.0h",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLocalizer(tt.lang)
			got := l.Getf(tt.key, tt.args...)
			if got != tt.want {
				t.Errorf("Getf(%q, %v) = %q, want %q", tt.key, tt.args, got, tt.want)
			}
		})
	}
}

func TestLocalizer_IsJapanese(t *testing.T) {
	tests := []struct {
		name string
		lang string
		want bool
	}{
		{
			name: "Japanese",
			lang: "ja",
			want: true,
		},
		{
			name: "English",
			lang: "en",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLocalizer(tt.lang)
			if got := l.IsJapanese(); got != tt.want {
				t.Errorf("IsJapanese() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		name     string
		langEnv  string
		wantLang string
	}{
		{
			name:     "Japanese locale",
			langEnv:  "ja_JP.UTF-8",
			wantLang: "ja",
		},
		{
			name:     "English locale",
			langEnv:  "en_US.UTF-8",
			wantLang: "en",
		},
		{
			name:     "Empty defaults to English",
			langEnv:  "",
			wantLang: "en",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save and restore LANG
			origLang := os.Getenv("LANG")
			defer os.Setenv("LANG", origLang)

			os.Setenv("LANG", tt.langEnv)
			got := DetectLanguage()
			if got != tt.wantLang {
				t.Errorf("DetectLanguage() = %q, want %q", got, tt.wantLang)
			}
		})
	}
}

func TestGlobalLocalizer(t *testing.T) {
	// Test SetGlobalLocalizer and GetGlobalLocalizer
	l := NewLocalizer("ja")
	SetGlobalLocalizer(l)

	got := GetGlobalLocalizer()
	if got.Lang() != "ja" {
		t.Errorf("GetGlobalLocalizer().Lang() = %q, want %q", got.Lang(), "ja")
	}

	// Test T function
	msg := T(MsgRecordSuccess)
	if msg != "ログを記録しました" {
		t.Errorf("T(%q) = %q, want %q", MsgRecordSuccess, msg, "ログを記録しました")
	}

	// Test Tf function
	msg = Tf(MsgStatsTotalLogs, 5)
	if msg != "総ログ数: 5件" {
		t.Errorf("Tf(%q, 5) = %q, want %q", MsgStatsTotalLogs, msg, "総ログ数: 5件")
	}
}

func TestNormalizeLang(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"ja", "ja"},
		{"ja_JP", "ja"},
		{"ja_JP.UTF-8", "ja"},
		{"JA_JP.UTF-8", "ja"},
		{"en", "en"},
		{"en_US", "en"},
		{"en_US.UTF-8", "en"},
		{"EN_US.UTF-8", "en"},
		{"fr_FR.UTF-8", "en"}, // Unsupported defaults to en
		{"", "en"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := normalizeLang(tt.input)
			if got != tt.want {
				t.Errorf("normalizeLang(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
