// File: internal/game/admin/service.go
package admin

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	GetPlayerInfo(ctx context.Context, targetUserID string) (string, error)
	GetRecentAudits(ctx context.Context) (string, error)
	CleanupCombat(ctx context.Context, adminID string) (string, error)
	PreviewReset(ctx context.Context, opts ResetOptions) (*ResetPreview, error)
	ApplyReset(ctx context.Context, opts ResetOptions) (*ResetPreview, error)
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
	countExpiredCombat, _ := combatCol.CountDocuments(ctx, bson.M{
		"expiresAt": bson.M{"$lt": time.Now().UTC()},
		"status":    bson.M{"$ne": "won"}, // Không đụng tới cơ duyên chưa nhận
	})

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
	res3, _ := combatCol.DeleteMany(ctx, bson.M{
		"expiresAt": bson.M{"$lt": time.Now().UTC()},
		"status":    bson.M{"$ne": "won"},
	})

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

func (s *adminSvc) GetPlayerInfo(ctx context.Context, targetUserID string) (string, error) {
	var profile bson.M // Using bson.M for flexibility as we don't need strong typing here
	err := s.db.Collection("players").FindOne(ctx, bson.M{"userId": targetUserID}).Decode(&profile)
	if err != nil {
		err = s.db.Collection("profiles").FindOne(ctx, bson.M{"userId": targetUserID}).Decode(&profile)
	}
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", fmt.Errorf("Không tìm thấy tàn hồn của đạo hữu này trong thiên địa.")
		}
		return "", err
	}

	var cult bson.M
	errCult := s.db.Collection("cultivation_profiles").FindOne(ctx, bson.M{"userId": targetUserID}).Decode(&cult)
	if errCult != nil {
		_ = s.db.Collection("cultivations").FindOne(ctx, bson.M{"userId": targetUserID}).Decode(&cult)
	}
	invCount, _ := s.db.Collection("inventories").CountDocuments(ctx, bson.M{"userId": targetUserID})
	var combat bson.M
	_ = s.db.Collection("combat_sessions").FindOne(ctx, bson.M{"userId": targetUserID, "expiresAt": bson.M{"$gt": time.Now().UTC()}}).Decode(&combat)

	info := fmt.Sprintf("**Đạo hiệu:** %v\n", profile["daoName"])
	if cult != nil {
		info += fmt.Sprintf("**Cảnh giới:** %v (Cấp %v)\n", cult["realm"], cult["level"])
		if path, ok := cult["path"]; ok && path != "" {
			info += fmt.Sprintf("**Đạo lộ:** %v\n", path)
		}
	} else {
		info += "**Cảnh giới:** Phàm nhân\n"
	}
	info += fmt.Sprintf("**Túi đồ:** %d vật phẩm\n", invCount)
	if combat != nil {
		info += fmt.Sprintf("**Trạng thái:** Đang độ kiếp tại bí cảnh '%v'\n", combat["areaId"])
	}
	return info, nil
}

func (s *adminSvc) GetRecentAudits(ctx context.Context) (string, error) {
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetLimit(5)
	cursor, err := s.db.Collection("admin_audit_logs").Find(ctx, bson.M{}, opts)
	if err != nil {
		return "", err
	}
	defer cursor.Close(ctx)

	var logs []AuditLog
	if err := cursor.All(ctx, &logs); err != nil {
		return "", err
	}

	if len(logs) == 0 {
		return "Sổ Thiên Cơ trống rỗng. Chưa có đạo pháp nào can thiệp thiên cơ.", nil
	}

	res := ""
	for _, l := range logs {
		status := "Thành công"
		if !l.Success {
			status = "Thất bại"
		}
		target := l.TargetUserID
		if target == "" {
			target = "Toàn cục"
		}
		res += fmt.Sprintf("`[%s]` **%s** | %s | Target: %s\n", l.CreatedAt.Format("01/02 15:04"), l.Action, status, target)
	}
	return res, nil
}

func (s *adminSvc) CleanupCombat(ctx context.Context, adminID string) (string, error) {
	combatCol := s.db.Collection("combat_sessions")
	res, err := combatCol.DeleteMany(ctx, bson.M{"expiresAt": bson.M{"$lt": time.Now().UTC()}, "status": bson.M{"$ne": "won"}})
	if err != nil {
		return "", err
	}
	s.LogAudit(ctx, AuditLog{AdminUserID: adminID, Action: "CLEANUP_COMBAT", Success: true, Details: map[string]interface{}{"deleted": res.DeletedCount}})
	return fmt.Sprintf("Đã dọn dẹp %d tàn dư combat vô chủ.", res.DeletedCount), nil
}

var resettableCollections = []CollectionResetSpec{
	{Name: "players", UserField: "userId", SupportsAll: true, PreserveOnAll: false},
	{Name: "profiles", UserField: "userId", SupportsAll: true, PreserveOnAll: false},
	{Name: "users", UserField: "userId", SupportsAll: true, PreserveOnAll: false},
	{Name: "cultivation_profiles", UserField: "userId", SupportsAll: true, PreserveOnAll: false},
	{Name: "cultivations", UserField: "userId", SupportsAll: true, PreserveOnAll: false},
	{Name: "wallets", UserField: "userId", SupportsAll: true, PreserveOnAll: false},
	{Name: "economies", UserField: "userId", SupportsAll: true, PreserveOnAll: false},
	{Name: "inventories", UserField: "userId", SupportsAll: true, PreserveOnAll: false},
	{Name: "item_instances", UserField: "userId", SupportsAll: true, PreserveOnAll: false},
	{Name: "equipment_sets", UserField: "userId", SupportsAll: true, PreserveOnAll: false},
	{Name: "equipments", UserField: "userId", SupportsAll: true, PreserveOnAll: false},
	{Name: "aptitudes", UserField: "userId", SupportsAll: true, PreserveOnAll: false},
	{Name: "pve_progress", UserField: "userId", SupportsAll: true, PreserveOnAll: false},
	{Name: "combat_sessions", UserField: "userId", SupportsAll: true, PreserveOnAll: false},
	{Name: "cooldowns", UserField: "key", SupportsAll: true, PreserveOnAll: false}, // Cooldown key is "userId:action"
	{Name: "menu_sessions", UserField: "userId", SupportsAll: true, PreserveOnAll: false},
	{Name: "admin_audit_logs", UserField: "adminUserId", SupportsAll: true, PreserveOnAll: true},
}

func (s *adminSvc) PreviewReset(ctx context.Context, opts ResetOptions) (*ResetPreview, error) {
	preview := &ResetPreview{
		Scope:        opts.Scope,
		TargetUserID: opts.TargetUserID,
	}

	for _, spec := range resettableCollections {
		var filter bson.M
		if opts.Scope == ResetScopeUser {
			if opts.TargetUserID == "" {
				return nil, fmt.Errorf("cần TargetUserID cho scope 'user'")
			}
			if spec.UserField == "key" { // Special case for cooldowns
				filter = bson.M{spec.UserField: primitive.Regex{Pattern: fmt.Sprintf("^%s:", opts.TargetUserID)}}
			} else {
				filter = bson.M{spec.UserField: opts.TargetUserID}
			}
		} else if opts.Scope == ResetScopeAll {
			if !spec.SupportsAll {
				continue
			}
			filter = bson.M{}
		} else {
			return nil, fmt.Errorf("scope không hợp lệ: %s", opts.Scope)
		}

		count, err := s.db.Collection(spec.Name).CountDocuments(ctx, filter)
		if err != nil {
			s.log.Warn("PreviewReset: không thể đếm collection", zap.String("collection", spec.Name), zap.Error(err))
			continue
		}

		if count > 0 {
			action := "Xóa"
			if opts.Scope == ResetScopeAll && spec.PreserveOnAll {
				action = "Bỏ qua"
			}
			preview.Collections = append(preview.Collections, CollectionResetPreview{
				Collection: spec.Name,
				Matched:    count,
				Action:     action,
			})
			if action == "Xóa" {
				preview.TotalMatched += count
			}
		}
	}

	return preview, nil
}

func (s *adminSvc) ApplyReset(ctx context.Context, opts ResetOptions) (*ResetPreview, error) {
	// Chạy lại preview để lấy số liệu chính xác trước khi xóa
	result, err := s.PreviewReset(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ApplyReset: preview thất bại: %w", err)
	}

	actionName := fmt.Sprintf("RESET_%s", strings.ToUpper(string(opts.Scope)))
	auditLog := AuditLog{
		AdminUserID:  opts.RequestedBy,
		Action:       actionName,
		TargetUserID: opts.TargetUserID,
		DryRun:       false,
		Details:      make(map[string]interface{}),
	}

	var errors []string
	for _, res := range result.Collections {
		if res.Action != "Xóa" {
			continue
		}

		var spec *CollectionResetSpec
		for i := range resettableCollections {
			if resettableCollections[i].Name == res.Collection {
				spec = &resettableCollections[i]
				break
			}
		}
		if spec == nil {
			continue
		}

		var filter bson.M
		if opts.Scope == ResetScopeUser {
			if spec.UserField == "key" {
				filter = bson.M{spec.UserField: primitive.Regex{Pattern: fmt.Sprintf("^%s:", opts.TargetUserID)}}
			} else {
				filter = bson.M{spec.UserField: opts.TargetUserID}
			}
		} else {
			filter = bson.M{}
		}

		delRes, err := s.db.Collection(spec.Name).DeleteMany(ctx, filter)
		if err != nil {
			s.log.Error("ApplyReset: xóa collection thất bại", zap.String("collection", spec.Name), zap.Error(err))
			errors = append(errors, spec.Name)
			continue
		}
		auditLog.Details[spec.Name] = delRes.DeletedCount
		s.log.Info("ApplyReset: dọn dẹp thành công", zap.String("collection", spec.Name), zap.Int64("deletedCount", delRes.DeletedCount))
	}

	if len(errors) > 0 {
		auditLog.Success = false
		auditLog.ErrorMessage = fmt.Sprintf("Lỗi xóa các collection: %s", strings.Join(errors, ", "))
		s.LogAudit(ctx, auditLog)
		return result, fmt.Errorf(auditLog.ErrorMessage)
	}

	auditLog.Success = true
	s.LogAudit(ctx, auditLog)
	return result, nil
}
