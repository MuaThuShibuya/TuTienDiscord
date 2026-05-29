// File: internal/game/pvecombat/adapters.go
// Chức năng: Đóng vai trò là cầu nối (Adapters) để chuyển đổi dữ liệu từ các service bên ngoài (Equipment, PvE Progress, Inventory)
//            thành các Interface mà PvECombat Service yêu cầu mà không tạo ra Import Cycle.

package pvecombat

import (
	"context"
	"errors"
	"math/rand"

	"github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/internal/game/characterstats"
	"github.com/whiskey/tu-tien-bot/internal/game/combat"
	"github.com/whiskey/tu-tien-bot/internal/game/inventory"
	"github.com/whiskey/tu-tien-bot/internal/game/item"
	"github.com/whiskey/tu-tien-bot/internal/game/pve"
	"github.com/whiskey/tu-tien-bot/internal/logger"
	"go.uber.org/zap"
)

// StatsAdapter gọi thẳng tới CharacterStats Pipeline
type StatsAdapter struct {
	statSvc characterstats.Provider
}

func NewStatsAdapter(statSvc characterstats.Provider) StatsProvider {
	return &StatsAdapter{statSvc: statSvc}
}

func (a *StatsAdapter) GetEffectiveStats(ctx context.Context, userID string) (combat.CombatStats, error) {
	return a.statSvc.GetEffectiveStats(ctx, userID)
}

// PvEAdapter kết nối PvE Progress Service & Registry với Combat Engine.
type PvEAdapter struct {
	progSvc pve.ProgressService
}

func NewPvEAdapter(progSvc pve.ProgressService) PvEProvider {
	return &PvEAdapter{progSvc: progSvc}
}

func (a *PvEAdapter) GetArea(areaID string) (pve.PvEAreaDefinition, error) {
	def, ok := pve.AreaRegistry[areaID]
	if !ok {
		return pve.PvEAreaDefinition{}, errors.New("không tìm thấy khu vực")
	}
	return def, nil
}
func (a *PvEAdapter) GetNextStage(ctx context.Context, userID, areaID string) (int, error) {
	return a.progSvc.GetNextStage(ctx, userID, areaID)
}
func (a *PvEAdapter) CanEnterArea(ctx context.Context, userID string, area pve.PvEAreaDefinition, cp int64, realm string) error {
	return a.progSvc.CanEnterArea(ctx, userID, area, cp, realm)
}
func (a *PvEAdapter) GenerateEncounter(area pve.PvEAreaDefinition, stage int, rng *rand.Rand) (*pve.EncounterDefinition, error) {
	return pve.GenerateEncounter(area, stage, rng)
}
func (a *PvEAdapter) MarkStageCleared(ctx context.Context, userID, areaID string, stage int) error {
	return a.progSvc.MarkStageCleared(ctx, userID, areaID, stage)
}

// GrantAdapter kết nối Inventory & Economy để trao thưởng.
type GrantAdapter struct{ invSvc inventory.Service }

func NewGrantAdapter(invSvc inventory.Service) RewardGrantService {
	return &GrantAdapter{invSvc: invSvc}
}
func (a *GrantAdapter) GrantExp(ctx context.Context, userID string, amount int64) error {
	if amount <= 0 {
		return errors.New("invalid reward amount")
	}
	return nil
} // TODO: Nối với Cultivation sau
func (a *GrantAdapter) GrantStones(ctx context.Context, userID string, amount int64) error {
	if amount <= 0 {
		return errors.New("invalid reward amount")
	}
	return nil
} // TODO: Nối với Economy sau
func (a *GrantAdapter) GrantItem(ctx context.Context, userID, defID string, quantity int64) error {
	if quantity <= 0 {
		return errors.New("invalid reward amount")
	}
	logger.L().Info("GrantAdapter: Trao vật phẩm",
		zap.String("userId", userID),
		zap.String("defId", defID),
		zap.Int64("qty", quantity))
	return a.invSvc.AddItem(ctx, userID, "", defID, quantity) // Truyền guildID rỗng để dùng kho global
}

func (a *GrantAdapter) PreflightInventoryCapacity(ctx context.Context, userID string, items []RewardItemPlan) error {
	inv, currentItems, err := a.invSvc.GetInventory(ctx, userID, "")
	if err != nil {
		return err
	}
	requiredNewSlots := 0
	stackableCount := make(map[string]int64)
	for _, it := range currentItems {
		stackableCount[it.DefinitionID] += it.Quantity
	}
	for _, req := range items {
		def, ok := item.GetDefinition(req.ItemID)
		if !ok {
			return errors.New("vật phẩm không tồn tại")
		}
		if def.Stackable {
			if _, exists := stackableCount[req.ItemID]; !exists {
				requiredNewSlots++
				stackableCount[req.ItemID] = int64(req.Quantity)
			}
		} else {
			requiredNewSlots += req.Quantity
		}
	}
	if inv.SlotUsage+requiredNewSlots > inv.SlotLimit {
		return apperrors.ErrInventoryFull
	}
	return nil
}
