// File: internal/discord/handlers/start_handler.go
// Phiên bản: v0.1.1
// Mục đích: Controller cho lệnh /start — đăng ký người chơi mới và khởi tạo dữ liệu.
// Bảo mật: guildId và userId lấy từ interaction, không từ input người dùng.
//           Response ephemeral nên chỉ người dùng đó thấy thông tin tài khoản.
// Ghi chú: Handler điều phối gọi tuần tự profile → cultivation → economy service.
//          Nếu wallet tạo thất bại sau khi profile đã tạo, người dùng có thể /start lại an toàn.

package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
	"github.com/whiskey/tu-tien-bot/internal/game/cultivation"
	"github.com/whiskey/tu-tien-bot/internal/game/economy"
	"github.com/whiskey/tu-tien-bot/internal/game/inventory"
	"github.com/whiskey/tu-tien-bot/internal/game/profile"
	"github.com/whiskey/tu-tien-bot/internal/logger"
)

// StartHandler xử lý lệnh /start.
type StartHandler struct {
	profileSvc     profile.Service
	cultivationSvc cultivation.Service
	economySvc     economy.Service
	inventorySvc   inventory.Service
	log            *zap.Logger
}

// NewStartHandler tạo StartHandler với các service đã inject.
func NewStartHandler(
	profileSvc profile.Service,
	cultivationSvc cultivation.Service,
	economySvc economy.Service,
	inventorySvc inventory.Service,
) *StartHandler {
	return &StartHandler{
		profileSvc:     profileSvc,
		cultivationSvc: cultivationSvc,
		economySvc:     economySvc,
		inventorySvc:   inventorySvc,
		log:            logger.L().Named("handler.start"),
	}
}

// Handle xử lý lệnh slash /start.
func (h *StartHandler) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Chỉ dùng được trong server
	if i.Member == nil || i.GuildID == "" {
		ui.RespondEphemeralError(s, i.Interaction, "Lệnh này chỉ dùng được trong server Discord.")
		return
	}

	userID := i.Member.User.ID
	guildID := i.GuildID
	username := i.Member.User.Username
	displayName := i.Member.Nick
	if displayName == "" {
		displayName = i.Member.User.GlobalName
	}
	if displayName == "" {
		displayName = username
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Defer response trong khi xử lý (chỉ người dùng thấy)
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})

	h.log.Info("/start được gọi",
		zap.String("userId", userID),
		zap.String("guildId", guildID),
		zap.String("username", username),
	)

	// 1. GetOrCreate hồ sơ người chơi
	player, err := h.profileSvc.GetOrCreate(ctx, userID, guildID, username, displayName)
	if err != nil {
		h.log.Error("/start: GetOrCreate profile thất bại", zap.String("userId", userID), zap.Error(err))
		ui.EditEphemeralError(s, i.Interaction, ui.MsgGenericError)
		return
	}

	// Cấp quà tân thủ (Túi đồ, Đan dược, Kiếm gỗ)
	_ = h.inventorySvc.GrantStarterItems(ctx, userID, guildID)

	// Kiểm tra người chơi đã đăng ký trước đó chưa (CreatedAt cách đây trên 5 giây)
	isNewPlayer := player.CreatedAt.After(time.Now().Add(-5 * time.Second))

	if !isNewPlayer {
		ui.EditEphemeralEmbed(s, i.Interaction, ui.WarningEmbed(
			fmt.Sprintf("Đạo hữu **%s** đã đăng ký trước đó rồi!\nHãy dùng `/menu` để tiếp tục hành trình.", player.DaoName),
		))
		return
	}

	// 2. Khởi tạo hồ sơ tu luyện
	cult, err := h.cultivationSvc.GetOrCreate(ctx, userID, guildID)
	if err != nil {
		h.log.Error("/start: GetOrCreate cultivation thất bại", zap.String("userId", userID), zap.Error(err))
		ui.EditEphemeralError(s, i.Interaction, ui.MsgGenericError)
		return
	}

	// 3. Khởi tạo ví — không fatal nếu thất bại, sẽ tự tạo lại khi truy cập lần sau
	_, err = h.economySvc.GetOrCreate(ctx, userID, guildID)
	if err != nil {
		h.log.Warn("/start: GetOrCreate wallet thất bại (không fatal)", zap.String("userId", userID), zap.Error(err))
	}

	// 4. Gửi thông báo chào mừng
	embed := ui.SuccessEmbed(
		"Chào Mừng Đến Với Vạn Pháp Tiên Nghịch!",
		fmt.Sprintf(
			"Đạo hữu **%s** đã bước vào thế giới tu tiên!\n\n"+
				"• Đạo hiệu: **%s**\n"+
				"• Cảnh giới khởi đầu: **%s tầng %d**\n"+
				"• Linh thạch: **500**\n"+
				"• Vé cơ duyên: **3**\n\n"+
				"Hãy dùng `/menu` để bắt đầu hành trình!",
			player.DaoName, player.DaoName, cult.Realm.DisplayName(), cult.RealmLevel,
		),
	)

	ui.EditEphemeralEmbed(s, i.Interaction, embed)
}
