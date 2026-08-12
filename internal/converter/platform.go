package converter

import (
	"time"

	platformmodel "github.com/yi-nology/git-sync-service/biz/model/platform"
	"github.com/yi-nology/git-sync-service/sync/model"
)

// ToPlatformInfo 将平台模型转为 IDL 生成的对外结构 platformmodel.PlatformInfo。
// 日期按 RFC3339 格式化，与 model.Platform（time.Time 默认 JSON）保持一致。
func ToPlatformInfo(p *model.Platform) *platformmodel.PlatformInfo {
	if p == nil {
		return nil
	}
	lastTestAt := ""
	if p.LastTestAt != nil {
		lastTestAt = p.LastTestAt.Format(time.RFC3339)
	}
	return &platformmodel.PlatformInfo{
		ID:             int64(p.ID),
		Key:            p.Key,
		Name:           p.Name,
		Type:           p.Type,
		InstanceUrl:    p.InstanceURL,
		ApiUrl:         p.APIURL,
		SkipTlsVerify:  p.SkipTLSVerify,
		CaCertPath:     p.CACertPath,
		ProxyUrl:       p.ProxyURL,
		IsDefault:      p.IsDefault,
		Status:         p.Status,
		LastTestResult: p.LastTestResult,
		RepoCount:      int32(p.RepoCount),
		CreatedAt:      p.CreatedAt.Format(time.RFC3339),
		LastTestAt:     lastTestAt,
		UpdatedAt:      p.UpdatedAt.Format(time.RFC3339),
	}
}

// ToPlatformList 批量转换平台。
func ToPlatformList(platforms []*model.Platform) []*platformmodel.PlatformInfo {
	result := make([]*platformmodel.PlatformInfo, 0, len(platforms))
	for _, p := range platforms {
		result = append(result, ToPlatformInfo(p))
	}
	return result
}
