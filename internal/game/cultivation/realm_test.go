// File: internal/game/cultivation/realm_test.go
package cultivation_test

import (
	"testing"

	"github.com/whiskey/tu-tien-bot/internal/game/cultivation"
)

func TestNormalizeRealmID(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Luyện Khí", "ngung_khi"},
		{"Linh Động Kỳ", "ngung_khi"},
		{"Kim Đan", "ket_dan"},
		{"Kết Đan", "ket_dan"},
		{"nguyen_anh", "nguyen_anh"},
		{"", "pham_nhan"},
		{"Cảnh Giới Ảo", "pham_nhan"},
	}

	for _, tt := range tests {
		res := cultivation.NormalizeRealmID(tt.input)
		if res != tt.expected {
			t.Errorf("Normalize(%q) = %q; mong đợi %q", tt.input, res, tt.expected)
		}
	}
}

func TestNextRealm(t *testing.T) {
	// Lên tầng
	r, lvl, adv, _ := cultivation.NextRealm("ngung_khi", 1)
	if r != "ngung_khi" || lvl != 2 || adv {
		t.Errorf("Lên tầng 2 lỗi")
	}

	// Đột phá cảnh giới
	r, lvl, adv, _ = cultivation.NextRealm("ngung_khi", 10)
	if r != "truc_co" || lvl != 1 || !adv {
		t.Errorf("Đột phá Trúc Cơ lỗi")
	}

	// Đột phá Dương Thực -> Khuy Niết
	r, lvl, adv, _ = cultivation.NextRealm("duong_thuc", 10)
	if r != "khuy_niet" || lvl != 1 || !adv {
		t.Errorf("Đột phá Khuy Niết lỗi")
	}
}

func TestGetRealmBaseStats(t *testing.T) {
	st := cultivation.GetRealmBaseStats("ngung_khi", 1)
	if st.MaxHP <= 0 || st.ATK <= 0 || st.Speed <= 0 {
		t.Errorf("Chỉ số không được âm hoặc 0: %+v", st)
	}
}
