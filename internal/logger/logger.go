// File: internal/logger/logger.go
// Phiên bản: v0.1.2
// Mục đích: Khởi tạo và cung cấp structured logger (zap) cho toàn bộ ứng dụng.
//           Hỗ trợ 2 format: "console" (đẹp khi local) và "json" (production/Render).
// Bảo mật: Không bao giờ log Discord token, MongoDB URI, password, hay secret bất kỳ.
// Ghi chú: Dùng logger.With(fields...) để đính kèm context (userId, guildId, command) vào log.

package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Options cấu hình logger — truyền vào Init() khi khởi động ứng dụng.
// Tất cả field đều có giá trị mặc định an toàn cho production nếu không được đặt.
type Options struct {
	Level         string // "debug" | "info" | "warn" | "error" — mặc định "info"
	Format        string // "console" | "json" — mặc định "json" (production-safe, không ANSI)
	Color         bool   // true = ANSI color, chỉ có tác dụng khi Format="console"
	CallerEnabled bool   // true = in tên file:dòng vào mỗi log entry
}

var globalLogger *zap.Logger

// Init khởi tạo global logger theo Options cho trước.
// Gọi một lần duy nhất khi khởi động ứng dụng, trước khi log bất cứ thứ gì.
func Init(opts Options) error {
	// Giá trị mặc định an toàn — json không in ANSI code vào file/pipe
	if opts.Format == "" {
		opts.Format = "json"
	}

	zapLevel, _ := parseLevel(opts.Level)

	var core zapcore.Core
	if strings.ToLower(opts.Format) == "console" {
		core = buildConsoleCore(zapLevel, opts.Color)
	} else {
		core = buildJSONCore(zapLevel)
	}

	// Stacktrace tự động tại ERROR trở lên; caller tùy cấu hình
	zapOpts := []zap.Option{
		zap.AddStacktrace(zapcore.ErrorLevel),
	}
	if opts.CallerEnabled {
		zapOpts = append(zapOpts, zap.AddCaller())
	}

	globalLogger = zap.New(core, zapOpts...)
	return nil
}

// buildJSONCore xây dựng core cho môi trường production (Render / VPS / Docker).
// Output: {"level":"info","time":"2026-05-25T14:16:07+07:00","logger":"discord.bot","msg":"...","field":"val"}
func buildJSONCore(level zapcore.Level) zapcore.Core {
	encCfg := zap.NewProductionEncoderConfig()
	encCfg.TimeKey = "time"
	encCfg.LevelKey = "level"
	encCfg.NameKey = "logger"
	encCfg.MessageKey = "msg"
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encCfg.EncodeLevel = zapcore.LowercaseLevelEncoder

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encCfg),
		zapcore.AddSync(os.Stdout),
		level,
	)
}

// buildConsoleCore xây dựng core dễ đọc cho local development.
// Output (ví dụ): [14:16:07]  INFO   discord.bot   Khởi động bot   app=tu-tien  version=0.1.1
// Các field primitive được in kiểu key=value; field phức tạp (object/array) in JSON.
func buildConsoleCore(level zapcore.Level, color bool) zapcore.Core {
	encCfg := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     shortTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeName:     paddedNameEncoder,
	}

	if color {
		encCfg.EncodeLevel = paddedColorLevelEncoder
	} else {
		encCfg.EncodeLevel = paddedLevelEncoder
	}

	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(encCfg),
		zapcore.AddSync(os.Stdout),
		level,
	)
}

// shortTimeEncoder định dạng thời gian kiểu [15:04:05] cho console log.
// Dùng giờ địa phương để dễ đối chiếu với đồng hồ máy tính.
func shortTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Local().Format("[15:04:05]"))
}

// paddedLevelEncoder in level hoa căn trái 5 ký tự: "INFO ", "DEBUG", "WARN ", "ERROR".
func paddedLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("%-5s", l.CapitalString()))
}

// paddedColorLevelEncoder in level có màu ANSI, đảm bảo chiều rộng hiển thị 5 ký tự.
// Màu: DEBUG=cyan, INFO=green, WARN=yellow, ERROR=red.
func paddedColorLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	// Mã ANSI không tính vào chiều rộng hiển thị nên không ảnh hưởng căn chỉnh tab
	switch l {
	case zapcore.DebugLevel:
		enc.AppendString("\033[36mDEBUG\033[0m") // cyan
	case zapcore.InfoLevel:
		enc.AppendString("\033[32mINFO \033[0m") // green
	case zapcore.WarnLevel:
		enc.AppendString("\033[33mWARN \033[0m") // yellow
	case zapcore.ErrorLevel:
		enc.AppendString("\033[31mERROR\033[0m") // red
	default:
		enc.AppendString(fmt.Sprintf("%-5s", l.CapitalString()))
	}
}

// paddedNameEncoder căn trái tên logger đến 14 ký tự để các cột trông đều nhau.
// Ví dụ: "discord.bot   ", "profile.service"
func paddedNameEncoder(name string, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("%-14s", name))
}

// L trả về global logger. Panic nếu Init chưa được gọi.
func L() *zap.Logger {
	if globalLogger == nil {
		panic("logger: phải gọi Init() trước khi dùng L()")
	}
	return globalLogger
}

// S trả về global SugaredLogger cho logging kiểu printf.
func S() *zap.SugaredLogger {
	return L().Sugar()
}

// Sync xả các log entry đang buffer. Gọi khi tắt ứng dụng.
func Sync() {
	if globalLogger != nil {
		_ = globalLogger.Sync()
	}
}

func parseLevel(level string) (zapcore.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn", "warning":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, nil
	}
}
