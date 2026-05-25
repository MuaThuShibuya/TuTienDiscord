// File: pkg/utils/time.go
// Phiên bản: v0.1.1
// Mục đích: Hàm tiện ích định dạng thời gian cho hiển thị tiếng Việt và Discord timestamp.
// Ghi chú: FormatDuration tạo chuỗi thời gian tiếng Việt dễ đọc để hiển thị cooldown.

package utils

import (
	"fmt"
	"time"
)

// FormatDuration chuyển đổi duration thành chuỗi tiếng Việt dễ đọc.
// Ví dụ: 1h30m15s → "1 giờ 30 phút"
func FormatDuration(d time.Duration) string {
	if d <= 0 {
		return "0 giây"
	}

	d = d.Round(time.Second)
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60

	switch {
	case h > 0 && m > 0:
		return fmt.Sprintf("%d giờ %d phút", h, m)
	case h > 0:
		return fmt.Sprintf("%d giờ", h)
	case m > 0 && s > 0:
		return fmt.Sprintf("%d phút %d giây", m, s)
	case m > 0:
		return fmt.Sprintf("%d phút", m)
	default:
		return fmt.Sprintf("%d giây", s)
	}
}

// DiscordTimestamp trả về chuỗi timestamp định dạng Discord.
// style: "R" = tương đối ("5 phút trước"), "F" = đầy đủ ngày giờ, "D" = ngày ngắn
func DiscordTimestamp(t time.Time, style string) string {
	return fmt.Sprintf("<t:%d:%s>", t.Unix(), style)
}

// StartOfDay trả về 00:00:00 UTC của ngày cho trước.
func StartOfDay(t time.Time) time.Time {
	y, m, d := t.UTC().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// IsSameDay trả về true nếu hai thời điểm thuộc cùng ngày UTC.
func IsSameDay(a, b time.Time) bool {
	ay, am, ad := a.UTC().Date()
	by, bm, bd := b.UTC().Date()
	return ay == by && am == bm && ad == bd
}
