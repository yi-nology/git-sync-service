package model

import (
	"time"

	"gorm.io/gorm"
)

// Platform 存储 Git 平台配置
type Platform struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	Key            string         `json:"key" gorm:"uniqueIndex;size:255;not null"`
	Name           string         `json:"name" gorm:"size:100;not null"`
	Type           string         `json:"type" gorm:"size:50;not null"`           // github, gitlab, gitea, gitee, gitcode, atomgit, tencent_code, custom
	InstanceURL    string         `json:"instanceUrl" gorm:"size:255"`            // 实例地址，如 github.com, gitlab.company.com
	APIURL         string         `json:"apiUrl" gorm:"size:255;not null"`        // API 地址，如 https://api.github.com
	AccessToken    string         `json:"-" gorm:"type:text"`                     // 访问令牌（加密存储）
	SkipTLSVerify  bool           `json:"skipTlsVerify" gorm:"default:false"`     // 跳过 TLS 证书验证
	CACertPath     string         `json:"caCertPath" gorm:"size:500"`             // 自定义 CA 证书路径
	ProxyURL       string         `json:"proxyUrl" gorm:"size:255"`               // HTTP 代理地址
	IsDefault      bool           `json:"isDefault" gorm:"default:false"`         // 是否为默认平台
	Status         string         `json:"status" gorm:"size:20;default:active"`   // 状态: active, error
	LastTestAt     *time.Time     `json:"lastTestAt"`                             // 最后测试时间
	LastTestResult string         `json:"lastTestResult" gorm:"size:500"`         // 最后测试结果
	RepoCount      int            `json:"repoCount" gorm:"default:0"`             // 关联的仓库数量
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Platform) TableName() string {
	return "platforms"
}

// PlatformType 平台类型常量
const (
	PlatformTypeGitHub      = "github"
	PlatformTypeGitLab      = "gitlab"
	PlatformTypeGitea       = "gitea"
	PlatformTypeGitee       = "gitee"
	PlatformTypeGitCode     = "gitcode"
	PlatformTypeAtomGit     = "atomgit"
	PlatformTypeTencentCode = "tencent_code"
	PlatformTypeCustom      = "custom"
)

// PlatformStatus 平台状态常量
const (
	PlatformStatusActive = "active"
	PlatformStatusError  = "error"
)

// PlatformAPIPaths 各平台的 API 路径
var PlatformAPIPaths = map[string]string{
	PlatformTypeGitHub:      "/api/v3",
	PlatformTypeGitLab:      "/api/v4",
	PlatformTypeGitea:       "/api/v1",
	PlatformTypeGitee:       "/api/v5",
	PlatformTypeGitCode:     "/api/v5",
	PlatformTypeAtomGit:     "/api/v1",
	PlatformTypeTencentCode: "/api/v3",
}

// PlatformDefaultInstances 各平台的默认实例地址
var PlatformDefaultInstances = map[string]string{
	PlatformTypeGitHub:      "github.com",
	PlatformTypeGitLab:      "gitlab.com",
	PlatformTypeGitea:       "gitea.com",
	PlatformTypeGitee:       "gitee.com",
	PlatformTypeGitCode:     "gitcode.com",
	PlatformTypeAtomGit:     "atomgit.com",
	PlatformTypeTencentCode: "git.code.tencent.com",
}

// GetAPIURL 根据实例地址生成 API URL
func GetAPIURL(platformType, instanceURL string) string {
	if instanceURL == "" {
		instanceURL = PlatformDefaultInstances[platformType]
	}
	apiPath := PlatformAPIPaths[platformType]
	if apiPath == "" {
		return ""
	}
	return "https://" + instanceURL + apiPath
}
