// File: internal/game/characterstats/service_test.go
package characterstats_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/whiskey/tu-tien-bot/internal/game/aptitude"
	"github.com/whiskey/tu-tien-bot/internal/game/characterstats"
	"github.com/whiskey/tu-tien-bot/internal/game/combat"
	"github.com/whiskey/tu-tien-bot/internal/game/cultivation"
	"github.com/whiskey/tu-tien-bot/internal/game/equipment"
	"github.com/whiskey/tu-tien-bot/internal/logger"
)

// --- Fakes ---
type fakeAptitude struct {
	aptitude.Service
	def *aptitude.AptitudeDefinition
}

func (f *fakeAptitude) GetProfile(ctx context.Context, userID string) (*aptitude.AptitudeProfile, *aptitude.AptitudeDefinition, error) {
	return nil, f.def, nil
}

type fakeCultivation struct {
	cultivation.Service
	prof *cultivation.CultivationProfile
	err  error
}

func (f *fakeCultivation) GetProfile(ctx context.Context, userID, guildID string) (*cultivation.CultivationProfile, error) {
	return f.prof, f.err
}

type fakeEquipment struct {
	equipment.Service
	stats combat.CombatStats
}

func (f *fakeEquipment) GetEffectiveStats(ctx context.Context, userID, guildID string) (equipment.CombatStats, error) {
	return equipment.CombatStats{MaxHP: f.stats.MaxHP, ATK: f.stats.ATK, DEF: f.stats.DEF, CritRate: f.stats.CritRate}, nil
}

func TestMain(m *testing.M) {
	if err := logger.Init(logger.Options{Level: "error", Format: "json"}); err != nil {
		panic("logger init thất bại: " + err.Error())
	}
	m.Run()
}

// --- Tests ---

func TestGetEffectiveStats_MissingAptitudeReturnsClearError(t *testing.T) {
	fApt := &fakeAptitude{def: nil}
	fCult := &fakeCultivation{prof: &cultivation.CultivationProfile{Realm: "ngung_khi", RealmLevel: 1}}
	fEq := &fakeEquipment{}
	svc := characterstats.NewPipelineService(fApt, fCult, fEq)
	_, err := svc.GetEffectiveStats(context.Background(), "u1")
	if err == nil {
		t.Fatal("Phải trả về lỗi")
	}
	if !strings.Contains(err.Error(), "missing aptitude") || !strings.Contains(err.Error(), "u1") {
		t.Errorf("Lỗi không chứa thông tin cần thiết: %v", err)
	}
}

func TestGetEffectiveStats_MissingCultivationReturnsClearError(t *testing.T) {
	apt := aptitude.Registry["apt_pham_tu"]
	fApt := &fakeAptitude{def: &apt}
	fCult := &fakeCultivation{prof: nil, err: fmt.Errorf("not found")}
	fEq := &fakeEquipment{}
	svc := characterstats.NewPipelineService(fApt, fCult, fEq)
	_, err := svc.GetEffectiveStats(context.Background(), "u1")
	if err == nil {
		t.Fatal("Phải trả về lỗi")
	}
	if !strings.Contains(err.Error(), "missing cultivation") || !strings.Contains(err.Error(), "u1") {
		t.Errorf("Lỗi không chứa thông tin cần thiết: %v", err)
	}
}

func TestEffectiveStats_NoEquipmentStillPositive(t *testing.T) {
	apt := aptitude.Registry["apt_pham_tu"]
	fApt := &fakeAptitude{def: &apt}
	fCult := &fakeCultivation{prof: &cultivation.CultivationProfile{Realm: "ngung_khi", RealmLevel: 1}}
	fEq := &fakeEquipment{} // Zero stats

	svc := characterstats.NewPipelineService(fApt, fCult, fEq)
	stats, err := svc.GetEffectiveStats(context.Background(), "u1")

	if err != nil {
		t.Fatalf("Lỗi: %v", err)
	}
	if stats.MaxHP <= 0 || stats.ATK <= 0 || stats.Speed <= 0 {
		t.Errorf("Stats cởi đồ phải lớn hơn 0: HP=%d, ATK=%d, Speed=%d", stats.MaxHP, stats.ATK, stats.Speed)
	}
}

func TestEffectiveStats_HighAptitudeStrongerThanMortal(t *testing.T) {
	aptMortal := aptitude.Registry["apt_pham_tu"]
	aptHeavenly := aptitude.Registry["apt_nghich_thien_dao_thai"]

	fCult := &fakeCultivation{prof: &cultivation.CultivationProfile{Realm: "ngung_khi", RealmLevel: 1}}
	fEq := &fakeEquipment{}

	svcMortal := characterstats.NewPipelineService(&fakeAptitude{def: &aptMortal}, fCult, fEq)
	svcHeavenly := characterstats.NewPipelineService(&fakeAptitude{def: &aptHeavenly}, fCult, fEq)

	statM, _ := svcMortal.GetEffectiveStats(context.Background(), "u1")
	statH, _ := svcHeavenly.GetEffectiveStats(context.Background(), "u2")

	if statH.CombatPower <= statM.CombatPower {
		t.Errorf("Nghịch Thiên (%d) phải mạnh hơn Phàm Tư (%d)", statH.CombatPower, statM.CombatPower)
	}
}

func TestEffectiveStats_HigherRealmStronger(t *testing.T) {
	apt := aptitude.Registry["apt_pham_tu"]
	fApt := &fakeAptitude{def: &apt}
	fEq := &fakeEquipment{}

	svc1 := characterstats.NewPipelineService(fApt, &fakeCultivation{prof: &cultivation.CultivationProfile{Realm: "ngung_khi", RealmLevel: 1}}, fEq)
	svc2 := characterstats.NewPipelineService(fApt, &fakeCultivation{prof: &cultivation.CultivationProfile{Realm: "truc_co", RealmLevel: 1}}, fEq)

	st1, _ := svc1.GetEffectiveStats(context.Background(), "u1")
	st2, _ := svc2.GetEffectiveStats(context.Background(), "u1")

	if st2.MaxHP <= st1.MaxHP || st2.ATK <= st1.ATK {
		t.Errorf("Trúc cơ phải mạnh hơn Ngưng Khí")
	}
}

func TestEffectiveStats_EquipmentAddsNotReplaces(t *testing.T) {
	apt := aptitude.Registry["apt_pham_tu"]
	fApt := &fakeAptitude{def: &apt}
	fCult := &fakeCultivation{prof: &cultivation.CultivationProfile{Realm: "ngung_khi", RealmLevel: 1}}

	svcBase := characterstats.NewPipelineService(fApt, fCult, &fakeEquipment{})
	svcEquip := characterstats.NewPipelineService(fApt, fCult, &fakeEquipment{stats: combat.CombatStats{ATK: 100}})

	stBase, _ := svcBase.GetEffectiveStats(context.Background(), "u1")
	stEquip, _ := svcEquip.GetEffectiveStats(context.Background(), "u1")

	if stEquip.ATK != stBase.ATK+100 {
		t.Errorf("Equipment phải cộng thêm 100 ATK (Có: %d, Chờ: %d)", stEquip.ATK, stBase.ATK+100)
	}
}

func TestEffectiveStats_CritRateNormalized(t *testing.T) {
	apt := aptitude.Registry["apt_pham_tu"] // base crit 0
	fApt := &fakeAptitude{def: &apt}
	fCult := &fakeCultivation{prof: &cultivation.CultivationProfile{Realm: "ngung_khi", RealmLevel: 1}}
	svc := characterstats.NewPipelineService(fApt, fCult, &fakeEquipment{stats: combat.CombatStats{CritRate: 1.5}}) // Buff lố 150%

	st, _ := svc.GetEffectiveStats(context.Background(), "u1")
	if st.CritRate > 1.0 {
		t.Errorf("CritRate không được vượt quá 1.0, nhận %f", st.CritRate)
	}
}
