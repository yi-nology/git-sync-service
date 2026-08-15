package executor

import (
	"errors"
	"strings"

	sdkprov "github.com/yi-nology/git-platform-sdk/provider"
	"github.com/yi-nology/git-sync-service/sync/model"
)

// ClassifyError categorizes an error into one of the predefined error types.
// 优先识别 SDK 结构化错误(provider HTTP 层),再回落到 git 操作关键词匹配。
func ClassifyError(err error) string {
	if err == nil {
		return ""
	}

	// SDK provider 结构化错误:鉴权/限流/未找到等按语义精确分类
	var perr *sdkprov.ProviderError
	if errors.As(err, &perr) {
		switch {
		case sdkprov.IsAuthentication(err):
			return model.ErrorAuth
		case sdkprov.IsRateLimited(err):
			return model.ErrorNetwork
		case sdkprov.IsNotFound(err):
			return model.ErrorConfig
		case perr.IsServerError():
			return model.ErrorNetwork
		case perr.IsClientError():
			return model.ErrorConfig
		}
	}

	// gitbackend / git 命令错误:关键词兜底
	msg := strings.ToLower(err.Error())

	// Auth errors: token issues, 401, 403
	authKeywords := []string{"401", "403", "authentication", "unauthorized", "permission denied", "token", "credential"}
	for _, kw := range authKeywords {
		if strings.Contains(msg, kw) {
			return model.ErrorAuth
		}
	}

	// Network errors: timeout, connection refused, DNS
	networkKeywords := []string{"timeout", "connection refused", "dial tcp", "network", "i/o error", "eof", "no route to host"}
	for _, kw := range networkKeywords {
		if strings.Contains(msg, kw) {
			return model.ErrorNetwork
		}
	}

	// Config errors: repo/branch not found
	configKeywords := []string{"not found", "does not exist", "repository not found", "branch not found"}
	for _, kw := range configKeywords {
		if strings.Contains(msg, kw) {
			return model.ErrorConfig
		}
	}

	// Git operation errors: non-fast-forward, conflict, rejected
	gitKeywords := []string{"non-fast-forward", "conflict", "rejected", "failed to push", "fetch first"}
	for _, kw := range gitKeywords {
		if strings.Contains(msg, kw) {
			return model.ErrorGit
		}
	}

	return model.ErrorUnknown
}
