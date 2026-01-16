package ui

import (
	"io"
	"os"
)

// Color codes for terminal output
const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	Bold    = "\033[1m"
)

// ColorWriter wraps an io.Writer and adds color support
type ColorWriter struct {
	w        io.Writer
	colorize bool
}

// NewColorWriter creates a new ColorWriter
// If colorize is true, color codes will be written
// If false, color codes are stripped
func NewColorWriter(w io.Writer, colorize bool) *ColorWriter {
	return &ColorWriter{w: w, colorize: colorize}
}

// NewAutoColorWriter creates a ColorWriter that auto-detects TTY
func NewAutoColorWriter(w io.Writer) *ColorWriter {
	colorize := IsTTY(w)
	return NewColorWriter(w, colorize)
}

// IsTTY checks if the writer is a terminal
func IsTTY(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		fi, err := f.Stat()
		if err != nil {
			return false
		}
		return (fi.Mode() & os.ModeCharDevice) != 0
	}
	return false
}

// Write implements io.Writer
func (c *ColorWriter) Write(p []byte) (n int, err error) {
	return c.w.Write(p)
}

// WriteColor writes text with the specified color
func (c *ColorWriter) WriteColor(color, text string) (n int, err error) {
	if c.colorize {
		return io.WriteString(c.w, color+text+Reset)
	}
	return io.WriteString(c.w, text)
}

// WriteSuccess writes text in green (success color)
func (c *ColorWriter) WriteSuccess(text string) (n int, err error) {
	return c.WriteColor(Green, text)
}

// WriteError writes text in red (error color)
func (c *ColorWriter) WriteError(text string) (n int, err error) {
	return c.WriteColor(Red, text)
}

// WriteWarning writes text in yellow (warning color)
func (c *ColorWriter) WriteWarning(text string) (n int, err error) {
	return c.WriteColor(Yellow, text)
}

// WriteInfo writes text in cyan (info color)
func (c *ColorWriter) WriteInfo(text string) (n int, err error) {
	return c.WriteColor(Cyan, text)
}

// WriteBold writes text in bold
func (c *ColorWriter) WriteBold(text string) (n int, err error) {
	return c.WriteColor(Bold, text)
}

// Colorize returns whether color output is enabled
func (c *ColorWriter) Colorize() bool {
	return c.colorize
}

// Colorize adds color codes to text if colorize is true
func Colorize(colorize bool, color, text string) string {
	if colorize {
		return color + text + Reset
	}
	return text
}
