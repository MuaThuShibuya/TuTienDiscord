// File: pkg/utils/numbers.go
// Version: v0.1
// Purpose: Number formatting helpers for display in Discord embeds.
// Notes: FormatNumber formats large integers with dot separator for Vietnamese readability.

package utils

import (
	"fmt"
	"strings"
)

// FormatNumber formats an integer with dot thousands separator.
// Example: 1234567 → "1.234.567"
func FormatNumber(n int64) string {
	s := fmt.Sprintf("%d", n)
	if n < 0 {
		s = s[1:]
	}

	var parts []string
	for len(s) > 3 {
		parts = append([]string{s[len(s)-3:]}, parts...)
		s = s[:len(s)-3]
	}
	if len(s) > 0 {
		parts = append([]string{s}, parts...)
	}

	result := strings.Join(parts, ".")
	if n < 0 {
		return "-" + result
	}
	return result
}

// Clamp constrains a value within [min, max].
func Clamp(val, min, max int64) int64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// ProgressBar renders a simple ASCII/Unicode progress bar for Discord embeds.
// Example: ProgressBar(3, 10, 10) → "███░░░░░░░"
func ProgressBar(current, max, width int) string {
	if max <= 0 || width <= 0 {
		return strings.Repeat("░", width)
	}
	filled := int(float64(current) / float64(max) * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}
