// File: internal/game/cultivation/realm_stats.go
package cultivation

import "github.com/whiskey/tu-tien-bot/internal/game/combat"

// GetRealmBaseStats cung cấp chỉ số gốc phụ thuộc vào cảnh giới và tầng.
func GetRealmBaseStats(realm string, level int) combat.CombatStats {
	id := NormalizeRealmID(realm)
	def, ok := RealmRegistry[id]
	if !ok {
		def = RealmRegistry["pham_nhan"]
	}

	// Clamp level
	if level < 1 {
		level = 1
	}
	if level > def.MaxLevel {
		level = def.MaxLevel
	}

	// Công thức: Base + PerLevel * (Level - 1)
	return combat.CombatStats{
		MaxHP: def.BaseStats.MaxHP + def.PerLevel.MaxHP*int64(level-1),
		ATK:   def.BaseStats.ATK + def.PerLevel.ATK*int64(level-1),
		DEF:   def.BaseStats.DEF + def.PerLevel.DEF*int64(level-1),
		Speed: def.BaseStats.Speed, // Tốc độ giữ nguyên theo đại cảnh giới
	}
}
