// File: internal/game/aptitude/models.go
package aptitude

type AptitudeRarity string

const (
	AptitudeMortal    AptitudeRarity = "mortal"    // Phàm tư
	AptitudeCommon    AptitudeRarity = "common"    // Bình thường
	AptitudeUncommon  AptitudeRarity = "uncommon"  // Khá
	AptitudeRare      AptitudeRarity = "rare"      // Linh căn tốt
	AptitudeEpic      AptitudeRarity = "epic"      // Thiên tài
	AptitudeLegendary AptitudeRarity = "legendary" // Yêu nghiệt
	AptitudeMythic    AptitudeRarity = "mythic"    // Nghịch thiên
	AptitudeHeavenly  AptitudeRarity = "heavenly"  // Thiên mệnh
)

type BaseStatBonus struct {
	MaxHP      int64
	ATK        int64
	DEF        int64
	Speed      int64
	CritRate   float64
	CritDamage float64
}

type GrowthStatMultiplier struct {
	MaxHPMultiplier float64
	ATKMultiplier   float64
	DEFMultiplier   float64
	SpeedMultiplier float64
}

type CultivationModifier struct {
	ExpGainMultiplier            float64
	BreakthroughChanceBonus      float64
	StaminaCostMultiplier        float64
	MeditationCooldownMultiplier float64
}

type AptitudeAffinity struct {
	Type       string // element, weapon_type, dao_path...
	RefID      string
	Multiplier float64
	FlatBonus  float64
}

type AptitudeDefinition struct {
	ID                   string
	Name                 string
	Description          string
	Rarity               AptitudeRarity
	Weight               int64
	BaseStats            BaseStatBonus
	GrowthStats          GrowthStatMultiplier
	CultivationModifiers CultivationModifier
	Affinities           []AptitudeAffinity
	Tags                 []string
}
