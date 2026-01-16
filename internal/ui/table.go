package ui

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

// RenderTable renders a table with headers and rows to the given writer
func RenderTable(w io.Writer, headers []string, rows [][]string) {
	if len(headers) == 0 {
		return
	}

	// Calculate column widths
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = displayWidth(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) {
				if w := displayWidth(cell); w > widths[i] {
					widths[i] = w
				}
			}
		}
	}

	// Render top border
	renderBorder(w, widths, "┌", "┬", "┐")

	// Render headers
	renderRow(w, headers, widths)

	// Render separator
	renderBorder(w, widths, "├", "┼", "┤")

	// Render data rows
	for _, row := range rows {
		renderRow(w, row, widths)
	}

	// Render bottom border
	renderBorder(w, widths, "└", "┴", "┘")
}

func renderBorder(w io.Writer, widths []int, left, mid, right string) {
	fmt.Fprint(w, left)
	for i, width := range widths {
		fmt.Fprint(w, strings.Repeat("─", width+2))
		if i < len(widths)-1 {
			fmt.Fprint(w, mid)
		}
	}
	fmt.Fprintln(w, right)
}

func renderRow(w io.Writer, cells []string, widths []int) {
	fmt.Fprint(w, "│")
	for i, width := range widths {
		cell := ""
		if i < len(cells) {
			cell = cells[i]
		}
		padding := width - displayWidth(cell)
		fmt.Fprintf(w, " %s%s │", cell, strings.Repeat(" ", padding))
	}
	fmt.Fprintln(w)
}

// displayWidth returns the display width of a string
// accounting for wide characters (CJK)
func displayWidth(s string) int {
	width := 0
	for _, r := range s {
		if isWideRune(r) {
			width += 2
		} else {
			width += 1
		}
	}
	return width
}

// isWideRune returns true if the rune is a wide character
func isWideRune(r rune) bool {
	// CJK characters and full-width characters
	return (r >= 0x1100 && r <= 0x115F) || // Hangul Jamo
		(r >= 0x2E80 && r <= 0x9FFF) || // CJK
		(r >= 0xAC00 && r <= 0xD7AF) || // Hangul Syllables
		(r >= 0xF900 && r <= 0xFAFF) || // CJK Compatibility Ideographs
		(r >= 0xFE10 && r <= 0xFE1F) || // Vertical forms
		(r >= 0xFE30 && r <= 0xFE6F) || // CJK Compatibility Forms
		(r >= 0xFF00 && r <= 0xFF60) || // Full-width forms
		(r >= 0xFFE0 && r <= 0xFFE6) // Full-width symbols
}

// FormatHours formats hours for display as "X.Xh"
func FormatHours(hours float64) string {
	return fmt.Sprintf("%.1fh", hours)
}

// TruncateString truncates a string to maxLen runes, adding "..." if truncated
func TruncateString(s string, maxLen int) string {
	if maxLen <= 3 {
		return s
	}

	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}

	return string(runes[:maxLen-3]) + "..."
}

// TruncateStringDisplay truncates a string to maxLen display width, adding "..." if truncated
func TruncateStringDisplay(s string, maxLen int) string {
	if maxLen <= 3 {
		return s
	}

	if displayWidth(s) <= maxLen {
		return s
	}

	runes := []rune(s)
	width := 0
	cutoff := 0

	for i, r := range runes {
		runeWidth := 1
		if isWideRune(r) {
			runeWidth = 2
		}
		if width+runeWidth > maxLen-3 {
			cutoff = i
			break
		}
		width += runeWidth
		cutoff = i + 1
	}

	return string(runes[:cutoff]) + "..."
}

// StringWidth returns the display width of a string (exported for testing)
func StringWidth(s string) int {
	return displayWidth(s)
}

// RuneCount returns the number of runes in a string
func RuneCount(s string) int {
	return utf8.RuneCountInString(s)
}
