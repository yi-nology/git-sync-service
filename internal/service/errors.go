package service

import "errors"

var (
	// ErrRepoNotFound is returned when a repository is not found.
	ErrRepoNotFound = errors.New("repo not found")

	// ErrTaskNotFound is returned when a sync task is not found.
	ErrTaskNotFound = errors.New("task not found")

	// ErrTaskDisabled is returned when attempting to run a disabled task.
	ErrTaskDisabled = errors.New("task is disabled")

	// ErrTaskRunning is returned when a task is already executing (concurrent run skipped).
	ErrTaskRunning = errors.New("task is already running")

	// ErrTooManyConcurrent is returned when the global concurrency limit is reached.
	ErrTooManyConcurrent = errors.New("too many concurrent sync tasks")

	// ErrRuleNotFound is returned when a webhook rule is not found.
	ErrRuleNotFound = errors.New("rule not found")

	// ErrEventNotFound is returned when a webhook event is not found.
	ErrEventNotFound = errors.New("event not found")
)
