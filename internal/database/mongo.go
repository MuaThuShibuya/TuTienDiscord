// File: internal/database/mongo.go
// Phiên bản: v0.1.1
// Mục đích: Quản lý vòng đời kết nối MongoDB Atlas — kết nối, ping, và ngắt kết nối.
// Bảo mật: MongoDB URI đọc từ config (env var). Không bao giờ log hoặc expose URI.
// Ghi chú: Luôn dùng context có timeout cho mọi thao tác DB. Gọi Disconnect khi tắt app.

package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/whiskey/tu-tien-bot/internal/logger"
)

const (
	connectTimeout = 10 * time.Second
	pingTimeout    = 5 * time.Second
	defaultTimeout = 10 * time.Second
)

// Client bọc MongoDB client và expose database mục tiêu.
type Client struct {
	client   *mongo.Client
	database *mongo.Database
	dbName   string
}

// Connect kết nối đến MongoDB Atlas và xác minh bằng ping.
// Trả về Client sẵn dùng, hoặc lỗi nếu kết nối thất bại.
func Connect(ctx context.Context, uri, dbName string) (*Client, error) {
	log := logger.L()

	connectCtx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	clientOpts := options.Client().
		ApplyURI(uri).
		SetServerSelectionTimeout(connectTimeout).
		SetConnectTimeout(connectTimeout)

	mongoClient, err := mongo.Connect(connectCtx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("database: connect thất bại: %w", err)
	}

	pingCtx, pingCancel := context.WithTimeout(ctx, pingTimeout)
	defer pingCancel()

	if err := mongoClient.Ping(pingCtx, nil); err != nil {
		_ = mongoClient.Disconnect(context.Background())
		return nil, fmt.Errorf("database: ping thất bại: %w", err)
	}

	log.Info("Đã kết nối MongoDB Atlas", zap.String("database", dbName))

	return &Client{
		client:   mongoClient,
		database: mongoClient.Database(dbName),
		dbName:   dbName,
	}, nil
}

// DB trả về mongo.Database mục tiêu để truy cập collection.
func (c *Client) DB() *mongo.Database {
	return c.database
}

// Collection trả về collection theo tên từ database mục tiêu.
func (c *Client) Collection(name string) *mongo.Collection {
	return c.database.Collection(name)
}

// NewContext trả về context với timeout chuẩn cho thao tác DB.
func NewContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), defaultTimeout)
}

// IsConnected ping MongoDB và trả về true nếu kết nối còn hoạt động.
func (c *Client) IsConnected() bool {
	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()
	return c.client.Ping(ctx, nil) == nil
}

// Disconnect đóng kết nối MongoDB graceful. Gọi khi tắt ứng dụng.
func (c *Client) Disconnect(ctx context.Context) error {
	logger.L().Info("Đang ngắt kết nối MongoDB")
	return c.client.Disconnect(ctx)
}
