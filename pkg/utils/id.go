// File: pkg/utils/id.go
// Phiên bản: v0.1.1
// Mục đích: Sinh ID ngẫu nhiên cho phiên menu và các entity khác.
// Bảo mật: Dùng crypto/rand để sinh session ID — ngăn tấn công đoán session.

package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// NewSessionID sinh session ID ngẫu nhiên bảo mật 16 byte (32 ký tự hex).
func NewSessionID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// Không nên xảy ra trong thực tế
		panic(fmt.Sprintf("utils.NewSessionID: crypto/rand thất bại: %v", err))
	}
	return hex.EncodeToString(b)
}

// NewShortID sinh ID ngắn hơn 8 byte ngẫu nhiên (16 ký tự hex) cho mục đích không bảo mật cao.
func NewShortID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("utils.NewShortID: crypto/rand thất bại: %v", err))
	}
	return hex.EncodeToString(b)
}
