// File: internal/apperrors/errors.go
// Phiên bản: v0.1.1
// Mục đích: Định nghĩa các loại lỗi chuẩn cho toàn bộ ứng dụng.
//           AppError mang thông điệp hiển thị cho user và nguyên nhân nội bộ tách biệt.
// Bảo mật: Không bao giờ trả về chi tiết lỗi nội bộ (stack trace, SQL, URI) cho Discord user.
// Ghi chú: Dùng errors.Is() để kiểm tra loại lỗi. Dùng UserFacing() để lấy thông điệp hiển thị.

package apperrors

import (
	"errors"
	"fmt"
)

// --- Sentinel errors (lỗi mẫu để kiểm tra bằng errors.Is) ---

var (
	ErrNotFound                   = errors.New("not found")          // Không tìm thấy dữ liệu
	ErrAlreadyExists              = errors.New("already exists")     // Dữ liệu đã tồn tại
	ErrInvalidInput               = errors.New("invalid input")      // Input không hợp lệ
	ErrInvalidDaoName             = errors.New("invalid dao name")   // Đạo hiệu không hợp lệ (quá ngắn, quá dài, rỗng)
	ErrPermissionDenied           = errors.New("permission denied")  // Không có quyền
	ErrSessionExpired             = errors.New("session expired")    // Phiên menu đã hết hạn
	ErrSessionNotOwner            = errors.New("not session owner")  // Không phải chủ phiên menu
	ErrCooldownActive             = errors.New("cooldown active")    // Đang trong thời gian hồi chiêu
	ErrInsufficientFunds          = errors.New("insufficient funds") // Không đủ tài nguyên
	ErrDatabaseTimeout            = errors.New("database timeout")   // Thao tác DB bị timeout
	ErrInsufficientCultivationExp = errors.New("insufficient cultivation exp")
	ErrInsufficientMindState      = errors.New("insufficient mind state")
	ErrMaxRealmReached            = errors.New("max realm reached")
	ErrBreakthroughFailed         = errors.New("breakthrough failed")
	ErrInvalidAction              = errors.New("invalid action")
	ErrPathAlreadyChosen          = errors.New("path already chosen") // Đã chọn đạo lộ

	ErrInventoryFull            = errors.New("inventory full")
	ErrItemNotFound             = errors.New("item not found")
	ErrItemNotUsable            = errors.New("item not usable")
	ErrInsufficientItemQuantity = errors.New("insufficient item quantity")
)

// --- AppError: lỗi có cấu trúc mang thông điệp tiếng Việt cho user ---

// AppError chứa mã lỗi máy đọc được, thông điệp tiếng Việt cho user, và nguyên nhân nội bộ.
type AppError struct {
	Code        string // Mã lỗi định danh, ví dụ: "COOLDOWN_ACTIVE"
	UserMessage string // Thông điệp hiển thị cho Discord user (tiếng Việt)
	Cause       error  // Nguyên nhân thật sự — không hiển thị cho user
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.UserMessage, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.UserMessage)
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

// New tạo một AppError với mã lỗi, thông điệp user, và nguyên nhân (có thể nil).
func New(code, userMessage string, cause error) *AppError {
	return &AppError{Code: code, UserMessage: userMessage, Cause: cause}
}

// --- Kiểm tra loại lỗi ---

func IsNotFound(err error) bool                 { return errors.Is(err, ErrNotFound) }
func IsAlreadyExists(err error) bool            { return errors.Is(err, ErrAlreadyExists) }
func IsInvalidDaoName(err error) bool           { return errors.Is(err, ErrInvalidDaoName) }
func IsSessionExpired(err error) bool           { return errors.Is(err, ErrSessionExpired) }
func IsSessionNotOwner(err error) bool          { return errors.Is(err, ErrSessionNotOwner) }
func IsCooldownActive(err error) bool           { return errors.Is(err, ErrCooldownActive) }
func IsInvalidInput(err error) bool             { return errors.Is(err, ErrInvalidInput) }
func IsInsufficientFunds(err error) bool        { return errors.Is(err, ErrInsufficientFunds) }
func IsPathAlreadyChosen(err error) bool        { return errors.Is(err, ErrPathAlreadyChosen) }
func IsInventoryFull(err error) bool            { return errors.Is(err, ErrInventoryFull) }
func IsItemNotFound(err error) bool             { return errors.Is(err, ErrItemNotFound) }
func IsItemNotUsable(err error) bool            { return errors.Is(err, ErrItemNotUsable) }
func IsInsufficientItemQuantity(err error) bool { return errors.Is(err, ErrInsufficientItemQuantity) }

// UserFacing trích xuất thông điệp tiếng Việt từ AppError.
// Nếu không phải AppError, trả về fallback.
func UserFacing(err error, fallback string) string {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.UserMessage
	}
	return fallback
}

// CooldownError mang thông tin thời gian còn lại để hiển thị cho user.
type CooldownError struct {
	Action    string
	Remaining string // Đã format, VD: "4 phút 30 giây"
}

func (e *CooldownError) Error() string {
	return fmt.Sprintf("cooldown active cho action %s, còn lại %s", e.Action, e.Remaining)
}
