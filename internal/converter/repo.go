package converter

import (
	repomodel "github.com/yi-nology/git-sync-service/biz/model/repo"
	"github.com/yi-nology/git-sync-service/sync/model"
)

func ToRepoInfo(r *model.Repo) *repomodel.RepoInfo {
	if r == nil {
		return nil
	}
	return &repomodel.RepoInfo{
		ID:            int64(r.ID),
		Key:           r.Key,
		Name:          r.Name,
		Platform:      r.Platform,
		PlatformOwner: r.PlatformOwner,
		PlatformRepo:  r.PlatformRepo,
		CloneUrl:      r.CloneURL,
		SshUrl:        r.SSHURL,
		DefaultBranch: r.DefaultBranch,
		Status:        r.Status,
		CreatedAt:     r.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToRepoInfoList(repos []*model.Repo) []*repomodel.RepoInfo {
	result := make([]*repomodel.RepoInfo, 0, len(repos))
	for _, r := range repos {
		result = append(result, ToRepoInfo(r))
	}
	return result
}
