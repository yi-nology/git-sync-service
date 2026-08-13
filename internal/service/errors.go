package service

import "errors"

var (
	// ErrRepoNotFound is returned when a repository is not found.
	ErrRepoNotFound = errors.New("repo not found")

	// ErrTaskNotFound is returned when a sync task is not found.
	ErrTaskNotFound = errors.New("task not found")

	// ErrTaskDisabled is returned when attempting to run a disabled task.
	ErrTaskDisabled = errors.New("task is disabled")

	// ErrRuleNotFound is returned when a webhook rule is not found.
	ErrRuleNotFound = errors.New("rule not found")

	// ErrEventNotFound is returned when a webhook event is not found.
	ErrEventNotFound = errors.New("event not found")
)
