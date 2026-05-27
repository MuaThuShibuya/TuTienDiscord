// File: internal/game/combat/actor.go
// Chức năng: Định nghĩa thực thể tham gia chiến đấu (Người chơi, Quái vật, Linh thú) kèm theo chỉ số (Stats) và máu hiện tại.

package combat

type ActorType string

const (
	ActorTypePlayer  ActorType = "player"
	ActorTypeMonster ActorType = "monster"
	ActorTypePet     ActorType = "pet" // Dành cho v0.6
)

// CombatActor đại diện cho một thực thể tham gia trận chiến (Người chơi, Quái, Linh Thú).
type CombatActor struct {
	ID            string         `bson:"id" json:"id"` // UserID hoặc MonsterInstanceID
	Type          ActorType      `bson:"type" json:"type"`
	Name          string         `bson:"name" json:"name"`
	Level         int            `bson:"level" json:"level"`
	Stats         CombatStats    `bson:"stats" json:"stats"`
	CurrentHP     int64          `bson:"currentHp" json:"currentHp"`
	CurrentRage   int64          `bson:"currentRage" json:"currentRage"`
	CurrentEnergy int64          `bson:"currentEnergy" json:"currentEnergy"`
	StatusEffects []StatusEffect `bson:"statusEffects" json:"statusEffects"`
}
