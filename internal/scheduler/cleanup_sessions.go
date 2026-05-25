// File: internal/scheduler/cleanup_sessions.go
// Phiên bản: v0.1.1
// Mục đích: Dọn dẹp định kỳ phiên menu đã hết hạn chưa được MongoDB TTL xóa tự động.
// Bảo mật: Chỉ xóa phiên đã qua expiresAt. Không có dữ liệu người chơi nào bị mất vĩnh viễn.
// Ghi chú: MongoDB TTL index xử lý phần lớn việc dọn dẹp tự động. Scheduler này là lớp dự phòng
//          chạy mỗi giờ. Có thể bỏ nếu TTL index hoạt động ổn định.

package scheduler

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/whiskey/tu-tien-bot/internal/logger"
)

const (
	sessionCleanupInterval = 1 * time.Hour
	sessionCleanupTimeout  = 30 * time.Second
)

// SessionCleaner xóa các document menu_sessions đã hết hạn theo lịch.
type SessionCleaner struct {
	col  *mongo.Collection
	log  *zap.Logger
	stop chan struct{}
}

// NewSessionCleaner tạo session cleaner cho MongoDB collection cho trước.
func NewSessionCleaner(db *mongo.Database) *SessionCleaner {
	return &SessionCleaner{
		col:  db.Collection("menu_sessions"),
		log:  logger.L().Named("scheduler.session_cleaner"),
		stop: make(chan struct{}),
	}
}

// Start bắt đầu vòng lặp dọn dẹp trong goroutine. Không blocking.
func (sc *SessionCleaner) Start() {
	sc.log.Info("Session cleaner đã bật", zap.Duration("interval", sessionCleanupInterval))
	go sc.run()
}

// Stop phát tín hiệu dừng vòng lặp dọn dẹp.
func (sc *SessionCleaner) Stop() {
	close(sc.stop)
}

func (sc *SessionCleaner) run() {
	ticker := time.NewTicker(sessionCleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			sc.cleanExpired()
		case <-sc.stop:
			sc.log.Info("Session cleaner đã dừng")
			return
		}
	}
}

func (sc *SessionCleaner) cleanExpired() {
	ctx, cancel := context.WithTimeout(context.Background(), sessionCleanupTimeout)
	defer cancel()

	filter := bson.M{"expiresAt": bson.M{"$lte": time.Now().UTC()}}
	result, err := sc.col.DeleteMany(ctx, filter)
	if err != nil {
		sc.log.Error("SessionCleaner: không xóa được phiên hết hạn", zap.Error(err))
		return
	}

	if result.DeletedCount > 0 {
		sc.log.Info("SessionCleaner: đã xóa phiên hết hạn",
			zap.Int64("count", result.DeletedCount))
	}
}
