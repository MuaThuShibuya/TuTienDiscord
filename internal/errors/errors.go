// File: internal/errors/errors.go
// Version: v0.1
// Purpose: Define common application error types and sentinel errors for consistent error handling.
// Security: Do not expose internal error details to Discord users directly; use user-friendly messages.
// Notes: Use errors.Is() and errors.As() for type checking. Wrap with fmt.Errorf("%w", err) to preserve chain.

package errors

import (
	"errors"
	"fmt"
)

// Sentinel errors for common game/app scenarios.
var (
	ErrNotFound          = errors.New("not found")
	ErrAlreadyExists     = errors.New("already exists")
	ErrInvalidInput      = errors.New("invalid input")
	ErrPermissionDenied  = errors.New("permission denied")
	ErrSessionExpired    = errors.New("session expired")
	ErrSessionNotOwner   = errors.New("not the session owner")
	ErrCooldownActive    = errors.New("cooldown is still active")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrDatabaseTimeout   = errors.New("database operation timed out")
	ErrRateLimited       = errors.New("rate limited")
)

// AppError is a structured error carrying a user-facing message and an internal cause.
type AppError struct {
	Code        string // machine-readable code, e.g. "COOLDOWN_ACTIVE"
	UserMessage string // shown to Discord user (in Vietnamese)
	Cause       error  // internal wrapped error, never shown to user
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.UserMessage, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.UserMessage)
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

// New creates an AppError with a code, user-facing message, and optional cause.
func New(code, userMessage string, cause error) *AppError {
	return &AppError{
		Code:        code,
		UserMessage: userMessage,
		Cause:       cause,
	}
}

// IsNotFound returns true if the error chain contains ErrNotFound.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsSessionExpired returns true if the error chain contains ErrSessionExpired.
func IsSessionExpired(err error) bool {
	return errors.Is(err, ErrSessionExpired)
}

// IsCooldownActive returns true if the error chain contains ErrCooldownActive.
func IsCooldownActive(err error) bool {
	return errors.Is(err, ErrCooldownActive)
}

// IsInsufficientFunds returns true if the error chain contains ErrInsufficientFunds.
func IsInsufficientFunds(err error) bool {
	return errors.Is(err, ErrInsufficientFunds)
}

// UserFacing extracts the user-visible message from an AppError,
// or returns the fallback string for unknown errors.
func UserFacing(err error, fallback string) string {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.UserMessage
	}
	return fallback
}
