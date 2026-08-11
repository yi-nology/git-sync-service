package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/yi-nology/git-sync-service/sync/model"
	"gorm.io/gorm"
)

func (s *Service) ReceiveWebhook(ctx context.Context, repoKey string, req *http.Request) error {
	repo, err := s.repoService.GetRepo(ctx, repoKey)
	if err != nil {
		return err
	}
	if repo == nil {
		return ErrRepoNotFound
	}

	prov, err := s.repoService.providerMgr.GetByURL(repo.CloneURL, repo.AccessToken)
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

	existing, err := s.webhookService.FindEventByEventID(event.ID)
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

	if err := s.webhookService.CreateWebhookEvent(whEvent); err != nil {
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
	s.webhookService.ApplyRules(ctx, repoKey, event, &s.lastTriggerTime, func(ctx context.Context, taskKey, trigger string) error {
		return s.RunTaskWithTrigger(ctx, taskKey, trigger)
	})
}

func (s *Service) ListRules(ctx context.Context, repoKey string) ([]*model.WebhookRule, error) {
	return s.webhookService.ListRules(ctx, repoKey)
}

func (s *Service) GetRule(ctx context.Context, id uint) (*model.WebhookRule, error) {
	return s.webhookService.GetRule(ctx, id)
}

func (s *Service) CreateRule(ctx context.Context, req *model.CreateRuleRequest) (*model.WebhookRule, error) {
	return s.webhookService.CreateRule(ctx, req)
}

func (s *Service) UpdateRule(ctx context.Context, req *model.UpdateRuleRequest) (*model.WebhookRule, error) {
	return s.webhookService.UpdateRule(ctx, req)
}

func (s *Service) DeleteRule(ctx context.Context, id uint) error {
	return s.webhookService.DeleteRule(ctx, id)
}

func (s *Service) ListEvents(ctx context.Context, repoKey string, offset, limit int) ([]*model.WebhookEvent, int64, error) {
	return s.webhookService.ListEvents(ctx, repoKey, offset, limit)
}

func (s *Service) RetryEvent(ctx context.Context, eventID uint) error {
	return s.webhookService.RetryEvent(ctx, eventID, func(ctx context.Context, repoKey string, event *model.WebhookEvent) {
		s.safeApplyRules(ctx, repoKey, event)
	})
}

func matchEventType(pattern, actual string) bool {
	if pattern == "" || pattern == "*" {
		return true
	}
	return pattern == actual
}
