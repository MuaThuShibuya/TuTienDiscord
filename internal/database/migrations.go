// File: internal/database/migrations.go
// Phiên bản: v0.2
// Mục đích: Tự động chạy các script cập nhật cấu trúc dữ liệu (schema) mỗi khi khởi động.
//           Điều này giúp "triệt tiêu lỗi" từ gốc thay vì chỉ bắt lỗi tạm bợ.

package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/whiskey/tu-tien-bot/internal/logger"
)

// AutoMigrate tự động kiểm tra và nâng cấp cấu trúc dữ liệu cho phù hợp với phiên bản code hiện tại.
func AutoMigrate(ctx context.Context, db *mongo.Database) error {
	log := logger.L().Named("database.migration")

	// 1. Chạy migration sửa lỗi schema v0.2 cho cultivation_profiles (string -> int/long)
	if err := migrateCultivationSchemaV02(ctx, db, log); err != nil {
		return err
	}

	return nil
}

func migrateCultivationSchemaV02(ctx context.Context, db *mongo.Database, log *zap.Logger) error {
	col := db.Collection("cultivation_profiles")

	// Tìm các document đang bị sai kiểu dữ liệu (vẫn còn là string từ bản cũ)
	filter := bson.M{
		"$or": []bson.M{
			{"mindState": bson.M{"$type": "string"}},
			{"realmLevel": bson.M{"$type": "string"}},
			{"cultivationExp": bson.M{"$type": "string"}},
			{"cultivationExpRequired": bson.M{"$type": "string"}},
			{"combatPower": bson.M{"$type": "string"}},
			{"stamina": bson.M{"$type": "string"}},
			{"maxStamina": bson.M{"$type": "string"}},
		},
	}

	// Cập nhật bằng Aggregation Pipeline để ép kiểu về số
	// Dùng $convert, nếu lỗi (ví dụ chuỗi "Bình tĩnh" không biến thành số được) thì đưa về default.
	update := []bson.M{
		{
			"$set": bson.M{
				"mindState":              bson.M{"$convert": bson.M{"input": "$mindState", "to": "int", "onError": 50, "onNull": 50}},
				"realmLevel":             bson.M{"$convert": bson.M{"input": "$realmLevel", "to": "int", "onError": 1, "onNull": 1}},
				"cultivationExp":         bson.M{"$convert": bson.M{"input": "$cultivationExp", "to": "long", "onError": 0, "onNull": 0}},
				"cultivationExpRequired": bson.M{"$convert": bson.M{"input": "$cultivationExpRequired", "to": "long", "onError": 200, "onNull": 200}},
				"combatPower":            bson.M{"$convert": bson.M{"input": "$combatPower", "to": "long", "onError": 100, "onNull": 100}},
				"stamina":                bson.M{"$convert": bson.M{"input": "$stamina", "to": "int", "onError": 100, "onNull": 100}},
				"maxStamina":             bson.M{"$convert": bson.M{"input": "$maxStamina", "to": "int", "onError": 100, "onNull": 100}},
			},
		},
	}

	res, err := col.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.ModifiedCount > 0 {
		log.Info("Đã tự động fix schema cho cultivation_profiles (string -> int/long)",
			zap.Int64("modifiedCount", res.ModifiedCount))
	}

	return nil
}
