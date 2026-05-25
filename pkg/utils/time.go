// File: pkg/utils/time.go
// Version: v0.1
// Purpose: Time formatting helpers for Vietnamese display and Discord timestamps.
// Notes: FormatDuration produces human-readable Vietnamese durations for cooldown display.

package utils

import (
	"fmt"
	"time"
)

// FormatDuration converts a duration into a human-readable Vietnamese string.
// Example: 1h30m15s → "1 giờ 30 phút 15 giây"
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

// DiscordTimestamp returns a Discord-formatted timestamp string.
// style: "R" = relative ("5 minutes ago"), "F" = full date/time, "d" = short date
func DiscordTimestamp(t time.Time, style string) string {
	return fmt.Sprintf("<t:%d:%s>", t.Unix(), style)
}

// StartOfDay returns midnight UTC for the given time.
func StartOfDay(t time.Time) time.Time {
	y, m, d := t.UTC().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// IsSameDay returns true if two times fall on the same UTC calendar day.
func IsSameDay(a, b time.Time) bool {
	ay, am, ad := a.UTC().Date()
	by, bm, bd := b.UTC().Date()
	return ay == by && am == bm && ad == bd
}
