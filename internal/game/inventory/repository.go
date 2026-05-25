// File: internal/game/inventory/repository.go
package inventory

import "context"

type Repository interface {
	GetOrCreate(ctx context.Context, userID, guildID string) (*Inventory, error)
	MarkStarterGranted(ctx context.Context, userID, guildID string) error
}
