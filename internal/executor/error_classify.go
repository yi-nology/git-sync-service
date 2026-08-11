package executor

import (
	"strings"

	"github.com/yi-nology/git-sync-service/sync/model"
)

// ClassifyError categorizes an error into one of the predefined error types
// based on the error message content.
func ClassifyError(err error) string {
	if err == nil {
		return ""
	}

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
