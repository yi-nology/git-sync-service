package sync

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/yi-nology/git-sync-service/internal/sync/model"
)

func (s *Service) ReceiveWebhook(ctx context.Context, repoKey string, req *http.Request) error {
	repo, err := s.repoDAO.FindByKey(repoKey)
	if err != nil {
		return err
	}
	if repo == nil {
		return fmt.Errorf("repo not found")
	}

	provider, err := s.providerMgr.GetProvider(repo)
	if err != nil {
		return err
	}

	if err := provider.ValidateWebhookSignature(req, repo.WebhookSecret); err != nil {
		return fmt.Errorf("invalid webhook signature: %w", err)
	}

	event, err := provider.ParseWebhookEvent(req, repo.WebhookSecret)
	if err != nil {
		return fmt.Errorf("parse webhook event failed: %w", err)
	}

	existing, _ := s.eventDAO.FindByEventID(event.ID)
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

	if err := s.eventDAO.Create(whEvent); err != nil {
		return err
	}

	go s.applyRules(ctx, repoKey, whEvent)

	return nil
}

func (s *Service) applyRules(ctx context.Context, repoKey string, event *model.WebhookEvent) {
	rules, err := s.ruleDAO.FindByRepoKey(repoKey)
	if err != nil {
		return
	}

	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}

		if !matchEventType(rule.EventType, event.EventType) {
			continue
		}

		if rule.BranchPattern != "" && !matchBranch(rule.BranchPattern, event.Branch) {
			continue
		}

		if rule.Action == "sync" && rule.SyncTaskKeys != "" {
			for _, taskKey := range splitAndTrim(rule.SyncTaskKeys, ",") {
				_ = s.RunTaskWithTrigger(ctx, taskKey, "webhook")
			}
		}
	}
}

func (s *Service) ListRules(repoKey string) ([]*model.WebhookRule, error) {
	return s.ruleDAO.FindByRepoKey(repoKey)
}

func (s *Service) GetRule(id uint) (*model.WebhookRule, error) {
	return s.ruleDAO.FindByID(id)
}

func (s *Service) CreateRule(req *model.CreateRuleRequest) (*model.WebhookRule, error) {
	rule := &model.WebhookRule{
		Name:          req.Name,
		RepoKey:       req.RepoKey,
		EventType:     req.EventType,
		BranchPattern: req.BranchPattern,
		Action:        req.Action,
		SyncTaskKeys:  req.SyncTaskKeys,
		MinInterval:   req.MinInterval,
		Enabled:       req.Enabled,
		Description:   req.Description,
	}

	if err := s.ruleDAO.Create(rule); err != nil {
		return nil, err
	}

	return rule, nil
}

func (s *Service) UpdateRule(req *model.UpdateRuleRequest) (*model.WebhookRule, error) {
	rule, err := s.ruleDAO.FindByID(req.ID)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return nil, fmt.Errorf("rule not found")
	}

	if req.Name != "" {
		rule.Name = req.Name
	}
	if req.EventType != "" {
		rule.EventType = req.EventType
	}
	rule.BranchPattern = req.BranchPattern
	if req.Action != "" {
		rule.Action = req.Action
	}
	rule.SyncTaskKeys = req.SyncTaskKeys
	if req.MinInterval > 0 {
		rule.MinInterval = req.MinInterval
	}
	rule.Enabled = req.Enabled
	rule.Description = req.Description

	if err := s.ruleDAO.Update(rule); err != nil {
		return nil, err
	}

	return rule, nil
}

func (s *Service) DeleteRule(id uint) error {
	return s.ruleDAO.Delete(id)
}

func (s *Service) ListEvents(repoKey string, limit int) ([]*model.WebhookEvent, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.eventDAO.FindByRepoKey(repoKey, limit)
}

func (s *Service) RetryEvent(ctx context.Context, eventID uint) error {
	event, err := s.eventDAO.FindByID(eventID)
	if err != nil {
		return err
	}
	if event == nil {
		return fmt.Errorf("event not found")
	}

	go s.applyRules(ctx, event.RepoKey, event)

	event.ProcessedAt = &time.Time{}
	event.Status = "processed"
	return s.eventDAO.Update(event)
}

func matchEventType(pattern, actual string) bool {
	if pattern == "" || pattern == "*" {
		return true
	}
	return pattern == actual
}
