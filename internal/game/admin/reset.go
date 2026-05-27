package admin

import "fmt"

type ResetScope string

const (
	ResetScopeUser ResetScope = "user"
	ResetScopeAll  ResetScope = "all"
)

type ResetOptions struct {
	Scope        ResetScope
	TargetUserID string
	DryRun       bool
	RequestedBy  string
}

type ResetPreview struct {
	Scope        ResetScope
	TargetUserID string
	Collections  []CollectionResetPreview
	TotalMatched int64
	Warnings     []string
}

func (p *ResetPreview) Summary() string {
	if p.TotalMatched == 0 {
		return "Không tìm thấy dữ liệu nào để reset."
	}
	summary := ""
	if p.Scope == ResetScopeUser {
		summary = fmt.Sprintf("Sẽ xóa **%d** ấn ký nhân quả của đạo hữu <@%s> từ các collection sau:\n", p.TotalMatched, p.TargetUserID)
	} else {
		summary = fmt.Sprintf("Sẽ xóa **%d** ấn ký nhân quả từ toàn bộ thiên địa:\n", p.TotalMatched)
	}
	for _, c := range p.Collections {
		if c.Action == "Xóa" {
			summary += fmt.Sprintf("- **%s**: %d bản ghi\n", c.Collection, c.Matched)
		}
	}
	return summary
}

type CollectionResetPreview struct {
	Collection string
	Matched    int64
	Action     string // "Xóa" hoặc "Bỏ qua"
}

type CollectionResetSpec struct {
	Name          string
	UserField     string
	SupportsAll   bool
	PreserveOnAll bool
}
