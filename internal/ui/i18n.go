package ui

import (
	"fmt"
	"os"
	"strings"
)

// SupportedLangs lists all supported languages
var SupportedLangs = []string{"ja", "en"}

// DefaultLang is the default language
const DefaultLang = "en"

// Localizer handles message localization
type Localizer struct {
	lang string
}

// NewLocalizer creates a new Localizer with the specified language
func NewLocalizer(lang string) *Localizer {
	// Normalize language code
	lang = normalizeLang(lang)
	if !isSupported(lang) {
		lang = DefaultLang
	}
	return &Localizer{lang: lang}
}

// NewAutoLocalizer creates a Localizer that auto-detects language from LANG env
func NewAutoLocalizer() *Localizer {
	lang := DetectLanguage()
	return NewLocalizer(lang)
}

// DetectLanguage detects language from LANG environment variable
func DetectLanguage() string {
	lang := os.Getenv("LANG")
	if lang == "" {
		return DefaultLang
	}
	return normalizeLang(lang)
}

// normalizeLang normalizes a language string to a supported language code
func normalizeLang(lang string) string {
	// Remove encoding suffix (e.g., "ja_JP.UTF-8" -> "ja_JP")
	if idx := strings.Index(lang, "."); idx != -1 {
		lang = lang[:idx]
	}

	// Convert to lowercase for comparison
	lang = strings.ToLower(lang)

	// Check for Japanese
	if strings.HasPrefix(lang, "ja") {
		return "ja"
	}

	// Default to English
	return "en"
}

// isSupported checks if a language is supported
func isSupported(lang string) bool {
	for _, l := range SupportedLangs {
		if l == lang {
			return true
		}
	}
	return false
}

// Get returns the localized message for the given key
func (l *Localizer) Get(key string) string {
	if msgs, ok := messages[l.lang]; ok {
		if msg, ok := msgs[key]; ok {
			return msg
		}
	}
	// Fallback to English
	if msgs, ok := messages[DefaultLang]; ok {
		if msg, ok := msgs[key]; ok {
			return msg
		}
	}
	return key
}

// Getf returns the localized message with format arguments
func (l *Localizer) Getf(key string, args ...any) string {
	return fmt.Sprintf(l.Get(key), args...)
}

// Lang returns the current language
func (l *Localizer) Lang() string {
	return l.lang
}

// IsJapanese returns true if the current language is Japanese
func (l *Localizer) IsJapanese() bool {
	return l.lang == "ja"
}

// Global localizer instance for convenience
var globalLocalizer *Localizer

// InitGlobalLocalizer initializes the global localizer
func InitGlobalLocalizer() {
	globalLocalizer = NewAutoLocalizer()
}

// SetGlobalLocalizer sets the global localizer
func SetGlobalLocalizer(l *Localizer) {
	globalLocalizer = l
}

// GetGlobalLocalizer returns the global localizer
func GetGlobalLocalizer() *Localizer {
	if globalLocalizer == nil {
		globalLocalizer = NewAutoLocalizer()
	}
	return globalLocalizer
}

// T is a shortcut for getting localized messages using the global localizer
func T(key string) string {
	return GetGlobalLocalizer().Get(key)
}

// Tf is a shortcut for getting formatted localized messages using the global localizer
func Tf(key string, args ...any) string {
	return GetGlobalLocalizer().Getf(key, args...)
}
