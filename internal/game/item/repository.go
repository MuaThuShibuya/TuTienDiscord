// File: internal/game/item/repository.go
package item

import "context"

type Repository interface {
	CreateInstance(ctx context.Context, inst *ItemInstance) error
	GetInstancesByUser(ctx context.Context, userID, guildID string) ([]*ItemInstance, error)
	GetInstanceByID(ctx context.Context, instanceID, userID, guildID string) (*ItemInstance, error)
	AdjustQuantity(ctx context.Context, instanceID, userID, guildID string, amount int64) error
	DeleteInstance(ctx context.Context, instanceID, userID, guildID string) error
	UpdateMetadata(ctx context.Context, instanceID, userID, guildID string, metadata map[string]interface{}) error
}
