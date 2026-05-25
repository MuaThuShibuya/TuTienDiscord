// File: pkg/utils/numbers.go
// Phiên bản: v0.1.1
// Mục đích: Hàm tiện ích định dạng số để hiển thị trong Discord embed.
// Ghi chú: FormatNumber định dạng số nguyên lớn với dấu chấm phân cách theo chuẩn Việt Nam.

package utils

import (
	"fmt"
	"strings"
)

// FormatNumber định dạng số nguyên với dấu chấm phân cách hàng nghìn.
// Ví dụ: 1234567 → "1.234.567"
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

// Clamp giới hạn giá trị trong khoảng [min, max].
func Clamp(val, min, max int64) int64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// ProgressBar render thanh tiến độ ASCII/Unicode cho Discord embed.
// Ví dụ: ProgressBar(3, 10, 10) → "███░░░░░░░"
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
