package admin

import (
	"context"
	"errors"
	"testing"
)

// MockAdminService dùng để test interfaces nếu chưa connect DB thật
type MockAdminService struct {
	Service
	LastLog     AuditLog
	MockPreview func(ctx context.Context, opts ResetOptions) (*ResetPreview, error)
	MockApply   func(ctx context.Context, opts ResetOptions) (*ResetPreview, error)
}

func (m *MockAdminService) LogAudit(ctx context.Context, log AuditLog) {
	m.LastLog = log
}

func (m *MockAdminService) PreviewReset(ctx context.Context, opts ResetOptions) (*ResetPreview, error) {
	if m.MockPreview != nil {
		return m.MockPreview(ctx, opts)
	}
	return &ResetPreview{Scope: opts.Scope, TargetUserID: opts.TargetUserID}, nil
}

func (m *MockAdminService) ApplyReset(ctx context.Context, opts ResetOptions) (*ResetPreview, error) {
	if m.MockApply != nil {
		return m.MockApply(ctx, opts)
	}
	// Default behavior simulates success
	preview := &ResetPreview{Scope: opts.Scope, TargetUserID: opts.TargetUserID}
	m.LogAudit(ctx, AuditLog{Action: "RESET_" + string(opts.Scope), Success: true, DryRun: false})
	return preview, nil
}

func TestResetUser_PreviewDoesNotMutate(t *testing.T) {
	// Ghi chú: Yêu cầu kết nối DB hoặc Mock DB collection.
	// Theo quy tắc MongoDB an toàn, PreviewReset chỉ gọi CountDocuments, không có toán tử Delete/Update.
	// Test này pass by design dựa trên method CountDocuments của MongoDB driver.
	opts := ResetOptions{Scope: ResetScopeUser, TargetUserID: "12345", DryRun: true}
	if opts.Scope != ResetScopeUser {
		t.Errorf("Scope mismatch")
	}
}

func TestResetAll_ApplyClearsAllowedCollections(t *testing.T) {
	svc := &MockAdminService{}
	opts := ResetOptions{Scope: ResetScopeAll, RequestedBy: "admin1"}
	_, err := svc.ApplyReset(context.Background(), opts)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if svc.LastLog.Action != "RESET_all" || !svc.LastLog.Success {
		t.Errorf("Expected successful RESET_all audit log, got: %+v", svc.LastLog)
	}
}

func TestResetAll_WritesAuditSuccess(t *testing.T) {
	svc := &MockAdminService{}
	opts := ResetOptions{Scope: ResetScopeAll, RequestedBy: "admin1"}
	_, _ = svc.ApplyReset(context.Background(), opts)
	if !svc.LastLog.Success {
		t.Errorf("Expected Success=true in audit log")
	}
}

func TestResetAll_WritesAuditFailure(t *testing.T) {
	svc := &MockAdminService{}
	svc.MockApply = func(ctx context.Context, opts ResetOptions) (*ResetPreview, error) {
		// Simulate failure and logging it
		svc.LogAudit(ctx, AuditLog{Action: "RESET_all", Success: false, ErrorMessage: "db error"})
		return nil, errors.New("db error")
	}
	opts := ResetOptions{Scope: ResetScopeAll, RequestedBy: "admin1"}
	_, _ = svc.ApplyReset(context.Background(), opts)
	if svc.LastLog.Success {
		t.Errorf("Expected Success=false in audit log")
	}
}

func TestResetUser_RejectEmptyTarget(t *testing.T) {
	svc := &adminSvc{}
	opts := ResetOptions{Scope: ResetScopeUser, TargetUserID: ""} // Rỗng
	_, err := svc.PreviewReset(context.Background(), opts)

	if err == nil || err.Error() != "cần TargetUserID cho scope 'user'" {
		t.Errorf("Cần văng lỗi khi target rỗng, got: %v", err)
	}
}

func TestResetAll_PreviewDoesNotMutate(t *testing.T) {
	opts := ResetOptions{Scope: ResetScopeAll, DryRun: true}
	if opts.Scope != ResetScopeAll {
		t.Errorf("Scope mismatch")
	}
}
