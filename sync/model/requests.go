package model

type CreateRepoRequest struct {
	Name        string `json:"name"`
	RemoteURL   string `json:"remoteUrl"`
	AccessToken string `json:"accessToken"`
}

type UpdateRepoRequest struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

type TestConnectionResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type PreviewSyncRequest struct {
	SourceRepoKey string `json:"sourceRepoKey"`
	SourceBranch  string `json:"sourceBranch"`
	TargetRepoKey string `json:"targetRepoKey"`
	TargetBranch  string `json:"targetBranch"`
}

type PreviewSyncResult struct {
	CanSync      bool   `json:"canSync"`
	SourceExists bool   `json:"sourceExists"`
	TargetExists bool   `json:"targetExists"`
	CommitCount  int    `json:"commitCount"`
	LatestCommit string `json:"latestCommit"`
	Message      string `json:"message"`
}
