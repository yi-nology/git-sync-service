package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/yi-nology/git-sync-service/sync/model"
	"gorm.io/gorm"
)

func (s *Service) ReceiveWebhook(ctx context.Context, repoKey string, req *http.Request) error {
	repo, err := s.repos.GetRepo(ctx, repoKey)
	if err != nil {
		return err
	}
	if repo == nil {
		return ErrRepoNotFound
	}

	prov, err := s.repos.GetProvider(repo.CloneURL, repo.AccessToken)
	if err != nil {
		return err
	}

	if err := prov.ValidateWebhookSignature(req, repo.WebhookSecret); err != nil {
		return fmt.Errorf("invalid webhook signature: %w", err)
	}

	event, err := prov.ParseWebhookEvent(req, repo.WebhookSecret)
	if err != nil {
		return fmt.Errorf("parse webhook event failed: %w", err)
	}

	existing, err := s.webhooks.FindEventByEventID(event.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existing != nil {
		return nil
	}

	actorName := ""
	if event.Actor != nil {
		actorName = event.Actor.Name
	}

	whEvent := &model.WebhookEvent{
		EventID:   event.ID,
		RepoKey:   repoKey,
		EventType: event.Type,
		Source:    string(event.Source),
		ActorName: actorName,
		Branch:    event.Branch,
		CommitSHA: event.CommitSHA,
		Payload:   event.RawPayload,
		Status:    model.StatusReceived,
	}

	if err := s.webhooks.CreateWebhookEvent(whEvent); err != nil {
		// 并发去重:event_id 有唯一索引,插入冲突说明另一个请求已处理该事件,直接幂等返回
		if dup, _ := s.webhooks.FindEventByEventID(event.ID); dup != nil {
			return nil
		}
		return err
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.safeApplyRules(s.bgCtx, repoKey, whEvent)
	}()

	return nil
}

func (s *Service) safeApplyRules(ctx context.Context, repoKey string, event *model.WebhookEvent) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic in applyRules", "repoKey", repoKey, "error", r)
		}
	}()
	s.applyRules(ctx, repoKey, event)
}

func (s *Service) applyRules(ctx context.Context, repoKey string, event *model.WebhookEvent) {
	eventID := event.ID
	s.webhooks.ApplyRules(ctx, repoKey, event, &s.lastTriggerTime, func(ctx context.Context, taskKey, trigger string, webhookEventID *uint) error {
		return s.RunTaskWithTrigger(ctx, taskKey, trigger, webhookEventID)
	}, &eventID)
}

func (s *Service) RetryEvent(ctx context.Context, eventID uint) error {
	event, err := s.webhooks.MarkEventProcessing(ctx, eventID)
	if err != nil {
		return err
	}
	// 纳入 WaitGroup + 用 bgCtx,优雅关停时能被等待/取消,不再泄露 goroutine
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.safeApplyRules(s.bgCtx, event.RepoKey, event)
		if err := s.webhooks.MarkEventProcessed(event); err != nil {
			slog.Error("mark event processed failed", "eventID", eventID, "error", err)
		}
	}()
	return nil
}

// ListRules returns webhook rules for a repository.
func (s *Service) ListRules(ctx context.Context, repoKey string) ([]*model.WebhookRule, error) {
	return s.webhooks.ListRules(ctx, repoKey)
}

// GetRule returns a webhook rule by ID.
func (s *Service) GetRule(ctx context.Context, id uint) (*model.WebhookRule, error) {
	return s.webhooks.GetRule(ctx, id)
}

// CreateRule creates a new webhook rule.
func (s *Service) CreateRule(ctx context.Context, req *model.CreateRuleRequest) (*model.WebhookRule, error) {
	return s.webhooks.CreateRule(ctx, req)
}

// UpdateRule updates an existing webhook rule.
func (s *Service) UpdateRule(ctx context.Context, req *model.UpdateRuleRequest) (*model.WebhookRule, error) {
	return s.webhooks.UpdateRule(ctx, req)
}

// DeleteRule deletes a webhook rule by ID.
func (s *Service) DeleteRule(ctx context.Context, id uint) error {
	return s.webhooks.DeleteRule(ctx, id)
}

// ListEvents returns webhook events for a repository.
func (s *Service) ListEvents(ctx context.Context, repoKey string, offset, limit int) ([]*model.WebhookEvent, int64, error) {
	return s.webhooks.ListEvents(ctx, repoKey, offset, limit)
}

func matchEventType(pattern, actual string) bool {
	if pattern == "" || pattern == "*" {
		return true
	}

	// Support comma-separated patterns
	for _, p := range splitAndTrim(pattern) {
		if p == actual {
			return true
		}
	}
	return false
}

func splitAndTrim(s string) []string {
	var result []string
	for _, part := range splitString(s, ",") {
		part = trimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

func splitString(s, sep string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, sep)
}

func trimSpace(s string) string {
	return strings.TrimSpace(s)
}
