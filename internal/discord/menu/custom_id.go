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
