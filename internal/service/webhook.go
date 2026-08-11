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
	repo, err := s.RepoService.GetRepo(ctx, repoKey)
	if err != nil {
		return err
	}
	if repo == nil {
		return ErrRepoNotFound
	}

	prov, err := s.RepoService.providerMgr.GetByURL(repo.CloneURL, repo.AccessToken)
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

	existing, err := s.WebhookService.FindEventByEventID(event.ID)
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
		Status:    "received",
	}

	if err := s.WebhookService.CreateWebhookEvent(whEvent); err != nil {
		return err
	}

	go s.safeApplyRules(context.Background(), repoKey, whEvent)

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
	s.WebhookService.ApplyRules(ctx, repoKey, event, &s.lastTriggerTime, func(ctx context.Context, taskKey, trigger string) error {
		return s.RunTaskWithTrigger(ctx, taskKey, trigger)
	})
}

func (s *Service) RetryEvent(ctx context.Context, eventID uint) error {
	return s.WebhookService.RetryEvent(ctx, eventID, func(ctx context.Context, repoKey string, event *model.WebhookEvent) {
		s.safeApplyRules(ctx, repoKey, event)
	})
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
