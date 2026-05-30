// File: internal/discord/menu/custom_id.go
// Phiên bản: v0.1.1
// Mục đích: Xây dựng và phân tích custom_id cho button/select menu Discord.
//           Tất cả custom_id trong hệ thống menu đều được tạo và đọc qua file này.
// Bảo mật: sessionId phải được nhúng vào custom_id và xác thực trước khi xử lý.
// Ghi chú: Format: "<domain>:<action>:<sessionId>[:<extra>]"
//          Ví dụ: "nav:back:abc123:main" | "cultivation:meditate:abc123"

package menu

import (
	"fmt"
	"strings"
)

const idSep = ":"

// --- Domain constants (nhóm chức năng) ---
const (
	DomainNav         = "nav"         // Điều hướng (làm mới, quay lại, đóng)
	DomainMenuSelect  = "menu"        // Select menu chọn category
	DomainProfile     = "profile"     // Hồ sơ người chơi
	DomainCultivation = "cultivation" // Tu luyện
	DomainInventory   = "inventory"   // Túi đồ
	DomainEquipment   = "equipment"   // Trang bị
	DomainAlchemy     = "alchemy"     // Lò luyện đan
	DomainPvE         = "pve"         // PvE Combat (Du Ngoạn / Bí Cảnh)
	DomainAdmin       = "admin"       // Thiên Cơ Các (Owner Admin)
	DomainShop        = "shop"        // Cửa hàng NPC & Đấu Giá Người Chơi
	// TODO v0.3+: thêm domain khi xây dựng tính năng mới
)

// --- Action constants (hành động cụ thể) ---
const (
	// Điều hướng
	ActionRefresh = "refresh" // Làm mới trang hiện tại
	ActionBack    = "back"    // Quay lại trang trước
	ActionClose   = "close"   // Đóng menu

	// Menu select
	ActionNavSelect  = "nav"  // Select menu hàng 1
	ActionNav2Select = "nav2" // Select menu hàng 2

	// Profile
	ActionRename = "rename" // Đổi đạo hiệu

	// Cultivation
	ActionMeditate     = "meditate"     // Tĩnh tu
	ActionClosedDoor   = "closeddoor"   // Bế quan
	ActionBodyTraining = "bodytraining" // Luyện thể
	ActionBreakthrough = "breakthrough" // Đột phá
	ActionChoosePath   = "choosepath"   // Chọn đạo lộ (Select Menu)

	// Inventory
	ActionInventoryUse              = "use"          // Dùng vật phẩm
	ActionInventoryDismantle        = "dismantle"    // Phân giải
	ActionInventoryDismantleConfirm = "dismantle_ok" // Xác nhận phân giải
	ActionInventoryPage             = "page"         // Phân trang

	// Equipment
	ActionEquipmentEquip   = "equip"   // Mặc trang bị
	ActionEquipmentUnequip = "unequip" // Tháo trang bị
	ActionEquipmentEnhance = "enhance" // Cường hóa

	// Alchemy
	ActionAlchemyCraft  = "craft"  // Luyện đan
	ActionAlchemyView   = "view"   // Xem chi tiết đan dược
	ActionAlchemyCancel = "cancel" // Hủy thao tác xem, quay lại

	// PvE
	ActionPvEMain    = "main"    // Menu tổng PvE
	ActionPvEDuNgoan = "dungoan" // Mở menu Du Ngoạn
	ActionPvEBiCanh  = "bicanh"  // Mở menu Bí Cảnh
	ActionPvEStart   = "start"   // Bắt đầu ải
	ActionPvEAttack  = "attack"  // Đánh thường
	ActionPvESkill   = "skill"   // Dùng kỹ năng
	ActionPvEAuto    = "auto"    // Tự động đánh
	ActionPvEClaim   = "claim"   // Nhận thưởng thắng trận
	ActionPvEEscape  = "escape"  // Bỏ chạy

	// Shop & Auction
	ActionShopMain                   = "main"         // Menu chính phường thị
	ActionShopGoNPC                  = "go_npc"       // Sang shop NPC
	ActionShopGoPlayer               = "go_player"    // Sang sàn đấu giá
	ActionShopNPCBuy                 = "npc_buy"      // Mua từ NPC
	ActionShopNPCSell                = "npc_sell"     // Bán cho NPC
	ActionShopNPCRefresh             = "npc_ref"      // Làm mới NPC Shop
	ActionShopNPCCategory            = "npc_cat"      // Chọn danh mục NPC
	ActionShopNPCPage                = "npc_page"     // Phân trang danh sách mua/bán NPC
	ActionShopNPCModeBuy             = "npc_mbuy"     // Chuyển sang mode Mua
	ActionShopNPCModeSell            = "npc_msell"    // Chuyển sang mode Bán
	ActionShopPlayerBuy              = "p_buy"        // Mua từ Đấu giá
	ActionShopPlayerList             = "p_list"       // Mở chọn đồ đăng bán (Không mở form)
	ActionShopPlayerListSelect_modal = "p_sel_modal"  // BẮT BUỘC kết thúc bằng _modal để mở Form
	ActionShopPlayerListApply        = "p_list_apply" // Xác nhận đăng bán
	ActionShopPlayerManage           = "p_manage"     // Quản lý hàng đang bán
	ActionShopPlayerCancel           = "p_cancel"     // Hủy đăng bán

	// Admin
	ActionAdminMain              = "main"
	ActionAdminMigrateDryRun     = "migrate_dry"
	ActionAdminMigrateModal      = "migrate_modal"
	ActionAdminMigrateApply      = "migrate_apply"
	ActionAdminPlayerLookupModal = "lookup_modal"
	ActionAdminPlayerLookupApply = "lookup_apply"
	ActionAdminResetUserModal    = "reset_user_modal"
	ActionAdminResetUserPreview  = "reset_user_preview"
	ActionAdminResetUserApply    = "reset_user_apply"
	ActionAdminResetAllPreview   = "reset_all_preview"
	ActionAdminResetAllApply     = "reset_all_apply"
	ActionAdminConfirmResetModal = "confirm_reset_modal"
	ActionAdminCombatCleanModal  = "clean_modal"
	ActionAdminCombatCleanApply  = "clean_apply"
	ActionAdminAuditLogs         = "audit_logs"
)

// ParsedID kết quả phân tích custom_id.
type ParsedID struct {
	Domain    string // Nhóm chức năng
	Action    string // Hành động
	SessionID string // ID phiên menu (để xác thực chủ sở hữu)
	Extra     string // Thông tin bổ sung tuỳ chọn (ví dụ: tên trang đích)
}

// Build tạo custom_id theo format chuẩn.
// extra là tuỳ chọn — chỉ thêm vào nếu có giá trị.
func Build(domain, action, sessionID string, extra ...string) string {
	parts := []string{domain, action, sessionID}
	if len(extra) > 0 && extra[0] != "" {
		parts = append(parts, extra[0])
	}
	return strings.Join(parts, idSep)
}

// Parse phân tích custom_id. Trả về lỗi nếu format không đúng.
func Parse(customID string) (*ParsedID, error) {
	parts := strings.SplitN(customID, idSep, 4)
	if len(parts) < 3 {
		return nil, fmt.Errorf("custom_id không hợp lệ: %q (cần ít nhất 3 phần)", customID)
	}
	result := &ParsedID{
		Domain:    parts[0],
		Action:    parts[1],
		SessionID: parts[2],
	}
	if len(parts) == 4 {
		result.Extra = parts[3]
	}
	return result, nil
}
