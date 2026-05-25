// File: internal/scheduler/cleanup_sessions.go
// Version: v0.1
// Purpose: Periodic cleanup of expired menu sessions not yet removed by MongoDB TTL.
// Security: Only deletes sessions past their expiresAt. No user data is permanently lost.
// Notes: MongoDB TTL index handles most cleanup automatically. This is a belt-and-suspenders
//        fallback that runs every hour. Can be removed if TTL index proves sufficient.

package scheduler

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/yourname/tu-tien-bot/internal/logger"
)

const (
	sessionCleanupInterval = 1 * time.Hour
	sessionCleanupTimeout  = 30 * time.Second
)

// SessionCleaner removes expired menu_sessions documents on a schedule.
type SessionCleaner struct {
	col  *mongo.Collection
	log  *zap.Logger
	stop chan struct{}
}

// NewSessionCleaner creates a session cleaner for the given MongoDB collection.
func NewSessionCleaner(db *mongo.Database) *SessionCleaner {
	return &SessionCleaner{
		col:  db.Collection("menu_sessions"),
		log:  logger.L().Named("scheduler.session_cleaner"),
		stop: make(chan struct{}),
	}
}

// Start begins the cleanup loop in a background goroutine. Non-blocking.
func (sc *SessionCleaner) Start() {
	sc.log.Info("Session cleaner started", zap.Duration("interval", sessionCleanupInterval))
	go sc.run()
}

// Stop signals the cleanup loop to exit.
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
			sc.log.Info("Session cleaner stopped")
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
		sc.log.Error("SessionCleaner: failed to delete expired sessions", zap.Error(err))
		return
	}

	if result.DeletedCount > 0 {
		sc.log.Info("SessionCleaner: removed expired sessions",
			zap.Int64("count", result.DeletedCount))
	}
}
