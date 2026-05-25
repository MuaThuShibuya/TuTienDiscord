// File: internal/database/mongo.go
// Version: v0.1
// Purpose: Manage MongoDB Atlas connection lifecycle — connect, ping, and disconnect.
// Security: MongoDB URI is read from config (env var). Never logged or exposed.
// Notes: Always use context with timeout for every DB operation. Call Disconnect on shutdown.

package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/yourname/tu-tien-bot/internal/logger"
)

const (
	connectTimeout = 10 * time.Second
	pingTimeout    = 5 * time.Second
	defaultTimeout = 10 * time.Second
)

// Client wraps the MongoDB client and exposes the target database.
type Client struct {
	client   *mongo.Client
	database *mongo.Database
	dbName   string
}

// Connect establishes a connection to MongoDB Atlas and verifies it with a ping.
// Returns a Client ready for use, or an error if connection fails.
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
		return nil, fmt.Errorf("database: connect failed: %w", err)
	}

	pingCtx, pingCancel := context.WithTimeout(ctx, pingTimeout)
	defer pingCancel()

	if err := mongoClient.Ping(pingCtx, nil); err != nil {
		_ = mongoClient.Disconnect(context.Background())
		return nil, fmt.Errorf("database: ping failed: %w", err)
	}

	log.Info("Connected to MongoDB Atlas", zap.String("database", dbName))

	return &Client{
		client:   mongoClient,
		database: mongoClient.Database(dbName),
		dbName:   dbName,
	}, nil
}

// DB returns the target mongo.Database for collection access.
func (c *Client) DB() *mongo.Database {
	return c.database
}

// Collection returns a named collection from the target database.
func (c *Client) Collection(name string) *mongo.Collection {
	return c.database.Collection(name)
}

// NewContext returns a context with the standard DB operation timeout.
func NewContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), defaultTimeout)
}

// IsConnected pings MongoDB and returns true if the connection is healthy.
func (c *Client) IsConnected() bool {
	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()
	return c.client.Ping(ctx, nil) == nil
}

// Disconnect cleanly closes the MongoDB connection. Call on application shutdown.
func (c *Client) Disconnect(ctx context.Context) error {
	logger.L().Info("Disconnecting from MongoDB")
	return c.client.Disconnect(ctx)
}
