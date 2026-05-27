// File: internal/discord/handlers/menu_handler.go
// Chức năng: Controller cho lệnh /menu — mở giao diện game tổng hợp.
// Bảo mật: Xác minh người chơi đã đăng ký trước khi tạo phiên menu.
//          Session dùng sessionId sinh ngẫu nhiên (crypto/rand).
// Ghi chú: /menu luôn mở trang Main. Điều hướng xảy ra qua menu.Router.
//          PageLoader inject để handler không trực tiếp gọi DB.

package handlers

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	apperrors "github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/internal/config"
	"github.com/whiskey/tu-tien-bot/internal/discord/menu"
	alchemymenu "github.com/whiskey/tu-tien-bot/internal/discord/menu/alchemy"
	cultivmenu "github.com/whiskey/tu-tien-bot/internal/discord/menu/cultivation"
	equipmenu "github.com/whiskey/tu-tien-bot/internal/discord/menu/equipment"
	invmenu "github.com/whiskey/tu-tien-bot/internal/discord/menu/inventory"
	mainmenu "github.com/whiskey/tu-tien-bot/internal/discord/menu/main"
	profilemenu "github.com/whiskey/tu-tien-bot/internal/discord/menu/profile"
	pvemenu "github.com/whiskey/tu-tien-bot/internal/discord/menu/pve"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
	"github.com/whiskey/tu-tien-bot/internal/game/alchemy"
	"github.com/whiskey/tu-tien-bot/internal/game/cultivation"
	"github.com/whiskey/tu-tien-bot/internal/game/economy"
	"github.com/whiskey/tu-tien-bot/internal/game/equipment"
	"github.com/whiskey/tu-tien-bot/internal/game/inventory"
	"github.com/whiskey/tu-tien-bot/internal/game/item"
	"github.com/whiskey/tu-tien-bot/internal/game/profile"
	"github.com/whiskey/tu-tien-bot/internal/logger"
	"github.com/whiskey/tu-tien-bot/pkg/utils"
)

// MenuHandler xử lý lệnh /menu và cung cấp PageLoader cho menu router.
type MenuHandler struct {
	cfg            *config.Config
	profileSvc     profile.Service
	cultivationSvc cultivation.Service
	economySvc     economy.Service
	inventorySvc   inventory.Service
	equipSvc       equipment.Service
	alchemySvc     alchemy.Service
	sessionSvc     menu.SessionService
	log            *zap.Logger
}

// NewMenuHandler tạo MenuHandler với các service đã inject.
func NewMenuHandler(
	cfg *config.Config,
	profileSvc profile.Service,
	cultivationSvc cultivation.Service,
	economySvc economy.Service,
	inventorySvc inventory.Service,
	equipSvc equipment.Service,
	alchemySvc alchemy.Service,
	sessionSvc menu.SessionService,
) *MenuHandler {
	return &MenuHandler{
		cfg:            cfg,
		profileSvc:     profileSvc,
		cultivationSvc: cultivationSvc,
		economySvc:     economySvc,
		inventorySvc:   inventorySvc,
		equipSvc:       equipSvc,
		alchemySvc:     alchemySvc,
		sessionSvc:     sessionSvc,
		log:            logger.L().Named("handler.menu"),
	}
}

// Handle xử lý lệnh slash /menu.
func (h *MenuHandler) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil || i.GuildID == "" {
		ui.RespondEphemeralError(s, i.Interaction, "Lệnh này chỉ dùng được trong server Discord.")
		return
	}

	userID := i.Member.User.ID
	guildID := i.GuildID
	channelID := i.ChannelID

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	h.log.Debug("/menu được gọi", zap.String("userId", userID), zap.String("guildId", guildID))

	player, err := h.profileSvc.GetPlayer(ctx, userID, guildID)
	if err != nil {
		if apperrors.IsNotFound(err) {
			ui.RespondEphemeralError(s, i.Interaction, ui.MsgNotRegistered)
			return
		}
		h.log.Error("/menu: GetPlayer thất bại", zap.String("userId", userID), zap.Error(err))
		ui.RespondEphemeralError(s, i.Interaction, ui.MsgGenericError)
		return
	}

	cult, err := h.cultivationSvc.GetOrCreate(ctx, userID, guildID)
	if err != nil {
		h.log.Error("/menu: GetOrCreate cultivation thất bại", zap.String("userId", userID), zap.Error(err))
		ui.RespondEphemeralError(s, i.Interaction, ui.MsgGenericError)
		return
	}

	wallet, err := h.economySvc.GetOrCreate(ctx, userID, guildID)
	if err != nil {
		h.log.Error("/menu: GetOrCreate wallet thất bại", zap.String("userId", userID), zap.Error(err))
		ui.RespondEphemeralError(s, i.Interaction, ui.MsgGenericError)
		return
	}

	session, err := h.sessionSvc.OpenMenu(ctx, userID, guildID, channelID, h.cfg.Menu.SessionTTL)
	if err != nil {
		h.log.Error("/menu: OpenMenu thất bại", zap.String("userId", userID), zap.Error(err))
		ui.RespondEphemeralError(s, i.Interaction, ui.MsgGenericError)
		return
	}

	// Thêm nút Admin nếu là Owner
	isAdmin := h.cfg.IsOwner(userID)
	vm := toMainMenuVM(session, player, cult, wallet)

	responseData := mainmenu.BuildMenuResponse(vm, isAdmin)

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: responseData,
	})
	if err != nil {
		h.log.Error("/menu: InteractionRespond thất bại", zap.String("userId", userID), zap.Error(err))
		return
	}

	msg, err := s.InteractionResponse(i.Interaction)
	if err == nil && msg != nil {
		_ = h.sessionSvc.SetMessageID(ctx, session.SessionID, msg.ID)
	}

	h.profileSvc.TouchLastActive(ctx, userID, guildID)
}

// PageLoaders trả về map page → loader function để menu.Router dùng khi điều hướng.
func (h *MenuHandler) PageLoaders() map[menu.Page]menu.PageLoader {
	return map[menu.Page]menu.PageLoader{
		menu.PageMain:        h.loadMainPage,
		menu.PageProfile:     h.loadProfilePage,
		menu.PageCultivation: h.loadCultivationPage,
		menu.PageInventory:   h.loadInventoryPage,
		menu.PageEquipment:   h.loadEquipmentPage,
		menu.PageAlchemy:     h.loadAlchemyPage,
		menu.PagePvE:         pvemenu.PvEMainLoader,
	}
}

func (h *MenuHandler) loadMainPage(ctx context.Context, session *menu.Session) (*discordgo.InteractionResponseData, error) {
	player, err := h.profileSvc.GetPlayer(ctx, session.UserID, session.GuildID)
	if err != nil {
		h.log.Error("loadMainPage failed", zap.String("step", "profile"), zap.Error(err))
		return nil, fmt.Errorf("loadMainPage profile: %w", err)
	}
	cult, err := h.cultivationSvc.GetOrCreate(ctx, session.UserID, session.GuildID)
	if err != nil {
		h.log.Error("loadMainPage failed", zap.String("step", "cultivation"), zap.Error(err))
		return nil, fmt.Errorf("loadMainPage cultivation: %w", err)
	}
	wallet, err := h.economySvc.GetOrCreate(ctx, session.UserID, session.GuildID)
	if err != nil {
		h.log.Error("loadMainPage failed", zap.String("step", "wallet"), zap.Error(err))
		return nil, fmt.Errorf("loadMainPage wallet: %w", err)
	}

	isAdmin := h.cfg.IsOwner(session.UserID)
	vm := toMainMenuVM(session, player, cult, wallet)
	return mainmenu.BuildMenuEdit(vm, isAdmin), nil
}

func (h *MenuHandler) loadProfilePage(ctx context.Context, session *menu.Session) (*discordgo.InteractionResponseData, error) {
	player, err := h.profileSvc.GetPlayer(ctx, session.UserID, session.GuildID)
	if err != nil {
		h.log.Error("loadProfilePage failed", zap.String("step", "profile"), zap.Error(err))
		return nil, fmt.Errorf("loadProfilePage profile: %w", err)
	}
	wallet, err := h.economySvc.GetOrCreate(ctx, session.UserID, session.GuildID)
	if err != nil {
		h.log.Error("loadProfilePage failed", zap.String("step", "wallet"), zap.Error(err))
		return nil, fmt.Errorf("loadProfilePage wallet: %w", err)
	}
	return profilemenu.BuildMenuResponse(toProfileMenuVM(session, player, wallet)), nil
}

func (h *MenuHandler) loadCultivationPage(ctx context.Context, session *menu.Session) (*discordgo.InteractionResponseData, error) {
	player, err := h.profileSvc.GetPlayer(ctx, session.UserID, session.GuildID)
	if err != nil {
		h.log.Error("loadCultivationPage failed", zap.String("step", "profile"), zap.Error(err))
		return nil, fmt.Errorf("loadCultivationPage profile: %w", err)
	}
	cult, err := h.cultivationSvc.GetOrCreate(ctx, session.UserID, session.GuildID)
	if err != nil {
		h.log.Error("loadCultivationPage failed", zap.String("step", "cultivation"), zap.Error(err))
		return nil, fmt.Errorf("loadCultivationPage cultivation: %w", err)
	}
	return cultivmenu.BuildMenuResponse(toCultivationMenuVM(session, player, cult)), nil
}

func (h *MenuHandler) loadInventoryPage(ctx context.Context, session *menu.Session) (*discordgo.InteractionResponseData, error) {
	player, err := h.profileSvc.GetPlayer(ctx, session.UserID, session.GuildID)
	if err != nil {
		h.log.Error("loadInventoryPage failed", zap.String("step", "profile"), zap.Error(err))
		return nil, err
	}
	_, items, err := h.inventorySvc.GetInventory(ctx, session.UserID, session.GuildID)
	if err != nil {
		if apperrors.IsNotFound(err) {
			items = []*item.ItemInstance{} // Fix: Túi đồ trống không phải là lỗi
		} else {
			h.log.Error("loadInventoryPage failed", zap.String("step", "get_inventory"), zap.Error(err))
			return nil, err
		}
	}
	return invmenu.BuildMenuResponse(toInventoryMenuVM(session, player, items)), nil
}

func (h *MenuHandler) loadEquipmentPage(ctx context.Context, session *menu.Session) (*discordgo.InteractionResponseData, error) {
	player, err := h.profileSvc.GetPlayer(ctx, session.UserID, session.GuildID)
	if err != nil {
		return nil, fmt.Errorf("loadEquipmentPage profile: %w", err)
	}
	cult, err := h.cultivationSvc.GetOrCreate(ctx, session.UserID, session.GuildID)
	if err != nil {
		return nil, fmt.Errorf("loadEquipmentPage cultivation: %w", err)
	}
	equipSet, err := h.equipSvc.GetEquipment(ctx, session.UserID, session.GuildID)
	if err != nil {
		return nil, fmt.Errorf("loadEquipmentPage equipment: %w", err)
	}
	_, allItems, err := h.inventorySvc.GetInventory(ctx, session.UserID, session.GuildID)
	if err != nil {
		return nil, fmt.Errorf("loadEquipmentPage inventory: %w", err)
	}
	return equipmenu.BuildMenuResponse(toEquipmentMenuVM(session, player, cult, equipSet, allItems)), nil
}

func (h *MenuHandler) loadAlchemyPage(ctx context.Context, session *menu.Session) (*discordgo.InteractionResponseData, error) {
	profile, err := h.alchemySvc.GetProfile(ctx, session.UserID, session.GuildID)
	if err != nil {
		h.log.Error("loadAlchemyPage failed", zap.String("step", "get_profile"), zap.Error(err))
		return nil, fmt.Errorf("loadAlchemyPage: %w", err)
	}

	expReq := int64(profile.Level * 100)
	expBar := fmt.Sprintf("`%s` %d/%d",
		utils.ProgressBar(int(profile.Exp), int(expReq), 10),
		profile.Exp, expReq)

	// Lấy túi đồ hiện tại để check nguyên liệu
	_, items, err := h.inventorySvc.GetInventory(ctx, session.UserID, session.GuildID)
	playerHas := make(map[string]int64)
	if err == nil {
		for _, it := range items {
			playerHas[it.DefinitionID] += it.Quantity
		}
	}

	var recipes []menu.RecipeVM
	var selectedRecipe *menu.RecipeVM

	for _, r := range alchemy.Recipes {
		var matStr string
		canCraft := true
		for reqDefID, reqQty := range r.RequiredItems {
			hasQty := playerHas[reqDefID]
			name := reqDefID
			if def, ok := item.GetDefinition(reqDefID); ok {
				name = def.Name
			}
			status := "Đủ"
			if hasQty < reqQty {
				status = "Thiếu"
				canCraft = false
			}
			matStr += fmt.Sprintf("• %s: %d/%d - %s\n", name, hasQty, reqQty, status)
		}

		recipeVM := menu.RecipeVM{
			ID:            r.ID,
			Name:          r.Name,
			SuccessRate:   fmt.Sprintf("%.0f%%", r.SuccessRate*100),
			LevelRequired: r.LevelRequired,
			Materials:     matStr,
			CanCraft:      canCraft,
		}

		recipes = append(recipes, recipeVM)

		if session.CurrentCategory == r.ID {
			copied := recipeVM
			selectedRecipe = &copied
		}
	}

	// Mẹo luyện đan
	tip := "Luyện đan có rủi ro thất bại, hãy thu thập đủ nguyên liệu và nâng cấp độ trước khi thử luyện linh đan phẩm chất cao."

	return alchemymenu.BuildMenuResponse(&menu.AlchemyMenuVM{
		SessionID:      session.SessionID,
		Level:          profile.Level,
		ExpBar:         expBar,
		Title:          "Dược Đồng", // Tạm thời hardcode, sau này bạn có thể map logic danh hiệu tại đây
		DailyTip:       tip,
		Recipes:        recipes,
		SelectedRecipe: selectedRecipe,
	}), nil
}

// --- ViewModel mapping ---

func toMainMenuVM(session *menu.Session, player *profile.Player, cult *cultivation.CultivationProfile, wallet *economy.Wallet) *menu.MainMenuVM {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	tip := ui.DailyTips[r.Intn(len(ui.DailyTips))]

	staminaBar := fmt.Sprintf("`%s` %d/%d",
		utils.ProgressBar(cult.Stamina, cult.MaxStamina, 10),
		cult.Stamina, cult.MaxStamina)

	expBar := fmt.Sprintf("`%s` %s/%s",
		utils.ProgressBar(int(cult.CultivationExp), int(cult.CultivationExpRequired), 10),
		utils.FormatNumber(cult.CultivationExp),
		utils.FormatNumber(cult.CultivationExpRequired))

	return &menu.MainMenuVM{
		SessionID:    session.SessionID,
		DaoName:      player.DaoName,
		RealmDisplay: fmt.Sprintf("%s tầng %d", cult.Realm.DisplayName(), cult.RealmLevel),
		CombatPower:  utils.FormatNumber(cult.CombatPower),
		MindState:    fmt.Sprintf("%s (%d/100)", cult.MindStateDisplayName(), cult.MindState),
		PathDisplay:  cult.Path.DisplayName(),
		StaminaBar:   staminaBar,
		ExpBar:       expBar,
		SpiritStones: utils.FormatNumber(wallet.SpiritStones),
		SpiritJades:  utils.FormatNumber(wallet.SpiritJades),
		FateTickets:  fmt.Sprintf("%d vé", wallet.FateTickets),
		DailyTip:     tip,
	}
}

func toProfileMenuVM(session *menu.Session, player *profile.Player, wallet *economy.Wallet) *menu.ProfileMenuVM {
	return &menu.ProfileMenuVM{
		SessionID:    session.SessionID,
		DaoName:      player.DaoName,
		DisplayName:  player.DisplayName,
		JoinedAt:     utils.DiscordTimestamp(player.CreatedAt, "D"),
		LastActive:   utils.DiscordTimestamp(player.LastActiveAt, "R"),
		SpiritStones: utils.FormatNumber(wallet.SpiritStones),
		SpiritJades:  utils.FormatNumber(wallet.SpiritJades),
		FateTickets:  fmt.Sprintf("%d vé", wallet.FateTickets),
	}
}

func toCultivationMenuVM(session *menu.Session, player *profile.Player, cult *cultivation.CultivationProfile) *menu.CultivationMenuVM {
	staminaBar := fmt.Sprintf("`%s` %d/%d",
		utils.ProgressBar(cult.Stamina, cult.MaxStamina, 10),
		cult.Stamina, cult.MaxStamina)

	expBar := fmt.Sprintf("`%s`\n%s / %s tu vi",
		utils.ProgressBar(int(cult.CultivationExp), int(cult.CultivationExpRequired), 12),
		utils.FormatNumber(cult.CultivationExp),
		utils.FormatNumber(cult.CultivationExpRequired))

	return &menu.CultivationMenuVM{
		SessionID:       session.SessionID,
		DaoName:         player.DaoName,
		RealmDisplay:    fmt.Sprintf("%s tầng %d", cult.Realm.DisplayName(), cult.RealmLevel),
		MindState:       fmt.Sprintf("%s (%d/100)", cult.MindStateDisplayName(), cult.MindState),
		PathDisplay:     cult.Path.DisplayName(),
		HasPath:         cult.Path != cultivation.PathNone,
		StaminaBar:      staminaBar,
		ExpBar:          expBar,
		CombatPower:     utils.FormatNumber(cult.CombatPower),
		CanBreakthrough: cult.CanBreakthrough(),
	}
}

func toInventoryMenuVM(session *menu.Session, player *profile.Player, items []*item.ItemInstance) *menu.InventoryMenuVM {
	const itemsPerPage = 20
	page := 1
	if session.CurrentCategory != "" {
		fmt.Sscanf(session.CurrentCategory, "%d", &page)
	}
	if page < 1 {
		page = 1
	}

	totalItems := len(items)
	totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage
	if totalPages == 0 {
		totalPages = 1
	}

	start := (page - 1) * itemsPerPage
	end := start + itemsPerPage
	if start >= totalItems {
		start = 0
	}
	if end > totalItems {
		end = totalItems
	}

	var itemVMs []menu.InventoryItemVM
	for _, it := range items[start:end] {
		vm := menu.InventoryItemVM{
			InstanceID: it.InstanceID,
			Name:       it.DefinitionID,
			Quantity:   it.Quantity,
		}
		if def, ok := item.GetDefinition(it.DefinitionID); ok {
			vm.Name = def.Name
			vm.Rarity = string(def.Rarity)
			vm.IsEquip = def.Type == item.TypeEquipment
			vm.IsUsable = def.Usable
		} else {
			vm.Name = "Vật phẩm lỗi (" + it.DefinitionID + ")"
			vm.Rarity = "D" // Tránh UI bị crash khi thiếu Rarity
		}
		itemVMs = append(itemVMs, vm)
	}

	var usableVMs []menu.InventoryItemVM
	for _, it := range items {
		if len(usableVMs) >= 25 {
			break
		}
		def, ok := item.GetDefinition(it.DefinitionID)
		if !ok || !def.Usable {
			continue
		}
		usableVMs = append(usableVMs, menu.InventoryItemVM{
			InstanceID: it.InstanceID,
			Name:       def.Name,
			Quantity:   it.Quantity,
			Rarity:     string(def.Rarity),
			IsUsable:   true,
		})
	}

	// Fix: Discord API sẽ báo lỗi 400 nếu Select Menu không có option nào.
	// Nếu không có vật phẩm khả dụng, thêm 1 option giả để tránh lỗi hiển thị UI.
	if len(usableVMs) == 0 {
		usableVMs = append(usableVMs, menu.InventoryItemVM{
			InstanceID: "empty",
			Name:       "Không có vật phẩm khả dụng",
			Quantity:   0,
			Rarity:     "D",
		})
	}

	return &menu.InventoryMenuVM{
		SessionID:   session.SessionID,
		DaoName:     player.DaoName,
		SlotUsage:   fmt.Sprintf("%d/50", totalItems),
		Items:       itemVMs,
		UsableItems: usableVMs,
		CurrentPage: page,
		TotalPages:  totalPages,
	}
}

func toEquipmentMenuVM(
	session *menu.Session,
	player *profile.Player,
	cult *cultivation.CultivationProfile,
	equipSet *equipment.EquipmentSet,
	allItems []*item.ItemInstance,
) *menu.EquipmentMenuVM {
	instanceMap := make(map[string]*item.ItemInstance, len(allItems))
	for _, it := range allItems {
		instanceMap[it.InstanceID] = it
	}

	toEquippedVM := func(slot equipment.EquipmentSlot, slotName string) *menu.EquippedItemVM {
		instanceID, ok := equipSet.Slots[string(slot)]
		if !ok || instanceID == "" {
			return nil
		}
		vm := &menu.EquippedItemVM{
			Slot:     string(slot),
			SlotName: slotName,
			Name:     "Vật phẩm lỗi",
			Rarity:   "D",
		}
		if it, found := instanceMap[instanceID]; found {
			if def, ok := item.GetDefinition(it.DefinitionID); ok {
				vm.Name = def.Name
				vm.Rarity = string(def.Rarity)
			} else {
				vm.Name = "Lỗi: " + it.DefinitionID
			}
		}
		return vm
	}

	vm := &menu.EquipmentMenuVM{
		SessionID:   session.SessionID,
		DaoName:     player.DaoName,
		CombatPower: utils.FormatNumber(cult.CombatPower),
		Weapon:      toEquippedVM(equipment.SlotWeapon, "Vũ Khí"),
		Armor:       toEquippedVM(equipment.SlotArmor, "Giáp"),
		Artifact:    toEquippedVM(equipment.SlotArtifact, "Pháp Bảo"),
		Treasure:    toEquippedVM(equipment.SlotTreasure, "Bảo Vật"),
		Boots:       toEquippedVM(equipment.SlotBoots, "Giày"),
	}

	equippedIDs := make(map[string]bool)
	for _, instanceID := range equipSet.Slots {
		if instanceID != "" {
			equippedIDs[instanceID] = true
		}
	}

	for _, it := range allItems {
		if len(vm.Equippable) >= 25 {
			break
		}
		if equippedIDs[it.InstanceID] {
			continue
		}
		def, ok := item.GetDefinition(it.DefinitionID)
		if !ok || def.Type != item.TypeEquipment {
			continue
		}
		slot := equipment.GetSlotForDefinition(it.DefinitionID)
		if slot == "" {
			continue
		}
		vm.Equippable = append(vm.Equippable, menu.EquippableItemVM{
			InstanceID:   it.InstanceID,
			DefinitionID: it.DefinitionID,
			Name:         def.Name,
			Rarity:       string(def.Rarity),
			SlotName:     slotDisplayName(slot),
		})
	}

	// Fix: Tránh lỗi Discord API khi Select Menu không có option.
	if len(vm.Equippable) == 0 {
		vm.Equippable = append(vm.Equippable, menu.EquippableItemVM{
			InstanceID:   "empty",
			DefinitionID: "empty",
			Name:         "Không có trang bị khả dụng",
			Rarity:       "D",
			SlotName:     "Trống",
		})
	}

	return vm
}

func slotDisplayName(slot equipment.EquipmentSlot) string {
	switch slot {
	case equipment.SlotWeapon:
		return "Vũ Khí"
	case equipment.SlotArmor:
		return "Giáp"
	case equipment.SlotArtifact:
		return "Pháp Bảo"
	case equipment.SlotTreasure:
		return "Bảo Vật"
	case equipment.SlotBoots:
		return "Giày"
	default:
		return string(slot)
	}
}
