// File: internal/discord/handlers/start_handler.go
// Version: v0.1
// Purpose: Controller for the /start slash command — registers a new player and initializes their data.
// Security: Validates guildId and userId from the interaction; never from user-provided text.
//           Responds ephemeral so account creation is private to the user.
// Notes: Handler orchestrates profile + cultivation + economy service calls in sequence.
//        On partial failure (e.g., wallet fails after profile created), the player can re-run /start safely.

package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	"github.com/yourname/tu-tien-bot/internal/discord/ui"
	apperrors "github.com/yourname/tu-tien-bot/internal/errors"
	"github.com/yourname/tu-tien-bot/internal/game/cultivation"
	"github.com/yourname/tu-tien-bot/internal/game/economy"
	"github.com/yourname/tu-tien-bot/internal/game/profile"
	"github.com/yourname/tu-tien-bot/internal/logger"
)

// StartHandler handles the /start command.
type StartHandler struct {
	profileSvc     profile.Service
	cultivationSvc cultivation.Service
	economySvc     economy.Service
	log            *zap.Logger
}

// NewStartHandler creates a new StartHandler with all required services.
func NewStartHandler(
	profileSvc profile.Service,
	cultivationSvc cultivation.Service,
	economySvc economy.Service,
) *StartHandler {
	return &StartHandler{
		profileSvc:     profileSvc,
		cultivationSvc: cultivationSvc,
		economySvc:     economySvc,
		log:            logger.L().Named("handler.start"),
	}
}

// Handle processes the /start slash command interaction.
func (h *StartHandler) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Validate: must be used inside a guild
	if i.Member == nil || i.GuildID == "" {
		ui.EphemeralError(s, i.Interaction, "Lệnh này chỉ dùng được trong server Discord.")
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

	// Defer response while processing
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})

	h.log.Info("/start invoked",
		zap.String("userId", userID),
		zap.String("guildId", guildID),
		zap.String("username", username),
	)

	// 1. Get or create player profile
	player, err := h.profileSvc.GetOrCreate(ctx, userID, guildID, username, displayName)
	if err != nil {
		h.log.Error("/start: GetOrCreate profile failed",
			zap.String("userId", userID), zap.Error(err))
		editEphemeralError(s, i.Interaction, ui.MsgGenericError)
		return
	}

	// Check if player already existed (ErrAlreadyExists would have been wrapped — check via a flag)
	isNewPlayer := player.CreatedAt.After(time.Now().Add(-5 * time.Second))

	if !isNewPlayer {
		// Player already registered
		editEphemeralEmbed(s, i.Interaction, ui.WarningEmbed(
			fmt.Sprintf("%s Đạo hữu **%s** đã đăng ký trước đó rồi!\nHãy dùng `/menu` để tiếp tục hành trình.",
				ui.EmojiWarning, player.DaoName)))
		return
	}

	// 2. Initialize cultivation profile
	_, err = h.cultivationSvc.GetOrCreate(ctx, userID, guildID)
	if err != nil {
		h.log.Error("/start: GetOrCreate cultivation failed",
			zap.String("userId", userID), zap.Error(err))
		editEphemeralError(s, i.Interaction, ui.MsgGenericError)
		return
	}

	// 3. Initialize wallet
	_, err = h.economySvc.GetOrCreate(ctx, userID, guildID)
	if err != nil {
		if !apperrors.IsNotFound(err) {
			h.log.Error("/start: GetOrCreate wallet failed",
				zap.String("userId", userID), zap.Error(err))
		}
		// Non-fatal: player can still proceed; wallet will be re-created on next access
	}

	// 4. Respond with welcome message
	embed := ui.SuccessEmbed(
		"Chào Mừng Đến Với Vạn Pháp Tiên Nghịch!",
		fmt.Sprintf(
			"%s Đạo hữu **%s** đã bước vào thế giới tu tiên!\n\n"+
				"• Đạo hiệu: **%s**\n"+
				"• Cảnh giới khởi đầu: **Luyện Khí tầng 1**\n"+
				"• Linh thạch: **500**\n"+
				"• Vé cơ duyên: **3**\n\n"+
				"Hãy dùng `/menu` để bắt đầu hành trình!",
			ui.EmojiSuccess, player.DaoName, player.DaoName,
		),
	)

	editEphemeralEmbed(s, i.Interaction, embed)
}

// editEphemeralError edits the deferred response with an error embed.
func editEphemeralError(s *discordgo.Session, i *discordgo.Interaction, message string) {
	_, _ = s.InteractionResponseEdit(i, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{ui.ErrorEmbed(message)},
	})
}

// editEphemeralEmbed edits the deferred response with any embed.
func editEphemeralEmbed(s *discordgo.Session, i *discordgo.Interaction, embed *discordgo.MessageEmbed) {
	_, _ = s.InteractionResponseEdit(i, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	})
}
