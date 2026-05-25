// File: pkg/utils/id.go
// Version: v0.1
// Purpose: Generate unique IDs for menu sessions and other entities.
// Notes: Uses crypto/rand for session IDs to prevent guessing attacks.

package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// NewSessionID generates a cryptographically random 16-byte session ID (32 hex chars).
func NewSessionID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// Fallback: should not happen in practice
		panic(fmt.Sprintf("utils.NewSessionID: crypto/rand failed: %v", err))
	}
	return hex.EncodeToString(b)
}

// NewShortID generates a shorter 8-byte random ID (16 hex chars) for non-security-critical use.
func NewShortID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("utils.NewShortID: crypto/rand failed: %v", err))
	}
	return hex.EncodeToString(b)
}
