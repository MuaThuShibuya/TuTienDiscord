// File: internal/game/aptitude/registry_test.go
package aptitude_test

import (
	"math/rand"
	"testing"

	"github.com/whiskey/tu-tien-bot/internal/game/aptitude"
)

func TestAptitudeRegistry_NotEmpty(t *testing.T) {
	if len(aptitude.Registry) == 0 {
		t.Fatal("Registry tư chất không được rỗng")
	}
}

func TestAptitudeRegistry_MinimumDiversity(t *testing.T) {
	if len(aptitude.Registry) < 12 {
		t.Errorf("Cần ít nhất 12 tư chất, hiện có %d", len(aptitude.Registry))
	}

	rarities := make(map[aptitude.AptitudeRarity]bool)
	for _, def := range aptitude.Registry {
		rarities[def.Rarity] = true
	}

	expectedRarities := []aptitude.AptitudeRarity{
		aptitude.AptitudeMortal, aptitude.AptitudeCommon, aptitude.AptitudeUncommon,
		aptitude.AptitudeRare, aptitude.AptitudeEpic, aptitude.AptitudeLegendary,
		aptitude.AptitudeMythic, aptitude.AptitudeHeavenly,
	}

	for _, r := range expectedRarities {
		if !rarities[r] {
			t.Errorf("Thiếu độ hiếm: %s", r)
		}
	}
}

func TestAptitudeRegistry_WeightsPositive(t *testing.T) {
	for _, def := range aptitude.Registry {
		if def.Weight <= 0 {
			t.Errorf("Tư chất %s có weight <= 0: %d", def.ID, def.Weight)
		}
	}
}

func TestAptitudeRegistry_HighRarityLowerWeight(t *testing.T) {
	heavenly := aptitude.Registry["apt_nghich_thien_dao_thai"].Weight
	mythic := aptitude.Registry["apt_hon_don_dao_the"].Weight
	common := aptitude.Registry["apt_tap_linh_can"].Weight

	if heavenly >= common || mythic >= common {
		t.Errorf("Độ hiếm cao phải có tỷ lệ (weight) thấp hơn. Heavenly: %d, Mythic: %d, Common: %d", heavenly, mythic, common)
	}
}

func TestGetRandomAptitude_ReturnsValidDefinition(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		def := aptitude.GetRandomAptitude(rng.Intn)
		if def.ID == "" {
			t.Fatal("Roll ra tư chất rỗng")
		}
		if _, ok := aptitude.Registry[def.ID]; !ok {
			t.Fatalf("Roll ra ID không tồn tại trong registry: %s", def.ID)
		}
	}
}

func TestGetRandomAptitude_DistributionSanity(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	counts := make(map[aptitude.AptitudeRarity]int)
	for i := 0; i < 5000; i++ {
		def := aptitude.GetRandomAptitude(rng.Intn)
		counts[def.Rarity]++
	}
	if counts[aptitude.AptitudeMortal] <= counts[aptitude.AptitudeHeavenly] {
		t.Errorf("Phân phối sai: Phàm tư (%d) phải xuất hiện nhiều hơn Thiên mệnh (%d)", counts[aptitude.AptitudeMortal], counts[aptitude.AptitudeHeavenly])
	}
}
