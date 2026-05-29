// File: internal/discord/ui/messages.go
// Version: v0.1
// Purpose: Standard Vietnamese message strings for common bot responses.
// Notes: Keep all user-facing text here for easy localization.

package ui

// System messages
const (
	MsgNotYourMenu    = "Đây không phải giao diện của đạo hữu. Hãy dùng `/menu` để mở giao diện riêng của mình."
	MsgSessionExpired = "Giao diện đã hết hạn. Hãy dùng `/menu` để mở lại."
	MsgGenericError   = "Đã xảy ra lỗi. Hãy thử lại sau."
	MsgCooldownActive = "Đạo hữu đang trong thời gian hồi phục."
	MsgComingSoon     = "Tính năng này đang được phát triển. Hãy chờ đón trong bản cập nhật tiếp theo!"
	MsgNotRegistered  = "Đạo hữu chưa đăng ký. Hãy dùng `/start` để bắt đầu hành trình tu tiên!"
	MsgAlreadyStarted = "Đạo hữu đã bắt đầu hành trình tu tiên rồi. Hãy dùng `/menu` để tiếp tục!"
)

// Tips shown in the main menu (rotate through these)
var DailyTips = []string{
	"Tĩnh tu mỗi ngày để tích lũy tu vi. Đừng bỏ lỡ buổi tu luyện nào!",
	"Boss server xuất hiện theo giờ. Tham gia cùng đạo hữu để nhận thưởng!",
	"Đan dược có thể tăng tốc tu luyện. Hãy khám phá cửa hàng đan dược!",
	"Tông môn cung cấp nhiều buff đặc biệt. Hãy gia nhập tông môn sớm!",
	"Gacha cơ duyên dùng vé cơ duyên, không dùng tiền thật. Tích lũy vé để quay!",
}
