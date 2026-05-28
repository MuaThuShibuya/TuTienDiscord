package pve

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

// PvEActionPayload chứa dữ liệu nội bộ của một hành động UI.
type PvEActionPayload struct {
	OwnerID         string
	MenuSessionID   string
	CombatSessionID string
	TargetID        string
	Action          string
	CreatedAt       time.Time
}

// ActionCache lưu trữ payload tương tác UI tạm thời trên RAM.
type ActionCache struct {
	mu    sync.RWMutex
	store map[string]PvEActionPayload
}

func NewActionCache() *ActionCache {
	return &ActionCache{
		store: make(map[string]PvEActionPayload),
	}
}

// Save tạo một token 8 ký tự và lưu payload.
func (c *ActionCache) Save(payload PvEActionPayload) string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	token := hex.EncodeToString(b) // 8 chars

	payload.CreatedAt = time.Now().UTC()
	c.mu.Lock()
	c.store[token] = payload
	c.mu.Unlock()
	return token
}

// Get lấy payload từ token.
func (c *ActionCache) Get(token string) (PvEActionPayload, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	p, ok := c.store[token]
	return p, ok
}
