package service

import (
	"github.com/yi-nology/git-platform-sdk/pkg/branchfilter"
)

func matchBranch(pattern, branch string) bool {
	return branchfilter.New(pattern).Match(branch)
}
