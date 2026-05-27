// File: internal/game/admin/service.go
package admin

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type AuditLog struct {
	ID           primitive.ObjectID     `bson:"_id,omitempty"`
	AdminUserID  string                 `bson:"adminUserId"`
	Action       string                 `bson:"action"`
	TargetUserID string                 `bson:"targetUserId,omitempty"`
	DryRun       bool                   `bson:"dryRun"`
	Success      bool                   `bson:"success"`
	ErrorMessage string                 `bson:"error,omitempty"`
	Details      map[string]interface{} `bson:"details,omitempty"`
	CreatedAt    time.Time              `bson:"createdAt"`
}

type Service interface {
	PreviewMigration(ctx context.Context) (string, error)
	ApplyMigration(ctx context.Context, adminID string) (string, error)
	LogAudit(ctx context.Context, log AuditLog)
}

type adminSvc struct {
	db  *mongo.Database
	log *zap.Logger
}

func NewService(db *mongo.Database, log *zap.Logger) Service {
	return &adminSvc{db: db, log: log.Named("game.admin")}
}

func (s *adminSvc) LogAudit(ctx context.Context, log AuditLog) {
	log.CreatedAt = time.Now().UTC()
	_, err := s.db.Collection("admin_audit_logs").InsertOne(ctx, log)
	if err != nil {
		s.log.Error("Không thể ghi audit log", zap.Error(err))
	}
}

// PreviewMigration thực hiện đếm số lượng document cần chuẩn hóa (Dry-run).
func (s *adminSvc) PreviewMigration(ctx context.Context) (string, error) {
	cultCol := s.db.Collection("cultivations")

	// Đếm realm lỗi tên cũ
	countLK, _ := cultCol.CountDocuments(ctx, bson.M{"realm": "Luyện Khí"})
	countLDK, _ := cultCol.CountDocuments(ctx, bson.M{"realm": "Linh Động Kỳ"})
	countKD, _ := cultCol.CountDocuments(ctx, bson.M{"realm": bson.M{"$in": []string{"Kim Đan", "Kết Đan"}}})

	// Đếm user thiếu tư chất (Aptitude)
	// Giả định collection là "users" hoặc "profiles" có nhúng aptitude
	profileCol := s.db.Collection("profiles")
	countNoAptitude, _ := profileCol.CountDocuments(ctx, bson.M{"aptitude": bson.M{"$exists": false}})

	// Đếm combat session treo
	combatCol := s.db.Collection("combat_sessions")
	countExpiredCombat, _ := combatCol.CountDocuments(ctx, bson.M{"expiresAt": bson.M{"$lt": time.Now().UTC()}})

	report := fmt.Sprintf("Thiên Cơ Soi Chiếu (Dry-run):\n"+
		"- Luyện Khí / Linh Động Kỳ (cũ): %d\n"+
		"- Kim Đan / Kết Đan (cũ): %d\n"+
		"- Hồ sơ thiếu Tư Chất: %d\n"+
		"- Combat session hết hạn cần dọn: %d",
		countLK+countLDK, countKD, countNoAptitude, countExpiredCombat)

	return report, nil
}

// ApplyMigration thực hiện update thật vào DB.
func (s *adminSvc) ApplyMigration(ctx context.Context, adminID string) (string, error) {
	cultCol := s.db.Collection("cultivations")

	var totalUpdated int64

	res, _ := cultCol.UpdateMany(ctx, bson.M{"realm": bson.M{"$in": []string{"Luyện Khí", "Linh Động Kỳ"}}}, bson.M{"$set": bson.M{"realm": "ngung_khi"}})
	totalUpdated += res.ModifiedCount

	res2, _ := cultCol.UpdateMany(ctx, bson.M{"realm": bson.M{"$in": []string{"Kim Đan", "Kết Đan"}}}, bson.M{"$set": bson.M{"realm": "ket_dan"}})
	totalUpdated += res2.ModifiedCount

	// Dọn dẹp combat
	combatCol := s.db.Collection("combat_sessions")
	res3, _ := combatCol.DeleteMany(ctx, bson.M{"expiresAt": bson.M{"$lt": time.Now().UTC()}})

	// TODO: Hàm sinh tư chất ngẫu nhiên cho user thiếu (gọi từ aptitude service nếu có)

	report := fmt.Sprintf("Chuẩn hóa thành công:\n- Đã sửa %d cảnh giới cũ.\n- Đã dọn %d combat session rác.", totalUpdated, res3.DeletedCount)

	// Ghi Audit
	s.LogAudit(ctx, AuditLog{
		AdminUserID: adminID,
		Action:      "APPLY_LEGACY_MIGRATION",
		DryRun:      false,
		Success:     true,
		Details:     map[string]interface{}{"realmsFixed": totalUpdated, "combatsCleared": res3.DeletedCount},
	})

	s.log.Info("Áp dụng Migration thành công",
		zap.String("adminUserId", adminID),
		zap.Int64("realmsFixed", totalUpdated),
		zap.Int64("combatsCleared", res3.DeletedCount),
	)
	return report, nil
}
