package service

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/yi-nology/git-sync-service/internal/dao"
	"github.com/yi-nology/git-sync-service/sync/model"
	"github.com/yi-nology/git-platform-sdk/pkg/branchfilter"
)

var (
	lastTriggerTime sync.Map
)

func (s *Service) ReceiveWebhook(ctx context.Context, repoKey string, req *http.Request) error {
	repo, err := s.repoDAO.FindByKey(repoKey)
	if err != nil {
		return err
	}
	if repo == nil {
		return fmt.Errorf("repo not found")
	}

	prov, err := s.providerMgr.GetProvider(repo)
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
	rules, err := s.ruleDAO.FindByRepoKey(repoKey)
	if err != nil {
		slog.Error("find rules failed", "repoKey", repoKey, "error", err)
		return
	}

	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}

		if !matchEventType(rule.EventType, event.EventType) {
			continue
		}

		if !branchfilter.New(rule.BranchPattern).Match(event.Branch) {
			continue
		}

		if rule.Action == "sync" && rule.SyncTaskKeys != "" {
			for _, taskKey := range strings.Split(rule.SyncTaskKeys, ",") {
				taskKey = strings.TrimSpace(taskKey)
				if taskKey == "" {
					continue
				}
				if rule.MinInterval > 0 {
					key := fmt.Sprintf("%s:%s", rule.RepoKey, taskKey)
					if lastTime, ok := lastTriggerTime.Load(key); ok {
						if time.Since(lastTime.(time.Time)) < time.Duration(rule.MinInterval)*time.Second {
							slog.Warn("skipping due to min interval", "taskKey", taskKey, "minInterval", rule.MinInterval)
							continue
						}
					}
					lastTriggerTime.Store(key, time.Now())
				}
				if err := s.RunTaskWithTrigger(ctx, taskKey, "webhook"); err != nil {
					slog.Error("run task failed", "taskKey", taskKey, "error", err)
				}
			}
		}
	}
}

func (s *Service) ListRules(ctx context.Context, repoKey string) ([]*model.WebhookRule, error) {
	return s.ruleDAO.FindByRepoKey(repoKey)
}

func (s *Service) GetRule(ctx context.Context, id uint) (*model.WebhookRule, error) {
	return s.ruleDAO.FindByID(id)
}

func (s *Service) CreateRule(ctx context.Context, req *model.CreateRuleRequest) (*model.WebhookRule, error) {
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

func (s *Service) UpdateRule(ctx context.Context, req *model.UpdateRuleRequest) (*model.WebhookRule, error) {
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

func (s *Service) DeleteRule(ctx context.Context, id uint) error {
	return s.ruleDAO.Delete(id)
}

func (s *Service) ListEvents(ctx context.Context, repoKey string, offset, limit int) ([]*model.WebhookEvent, int64, error) {
	page := dao.DefaultPagination(offset, limit)
	return s.eventDAO.FindByRepoKey(repoKey, page)
}

func (s *Service) RetryEvent(ctx context.Context, eventID uint) error {
	event, err := s.eventDAO.FindByID(eventID)
	if err != nil {
		return err
	}
	if event == nil {
		return fmt.Errorf("event not found")
	}

	go s.safeApplyRules(context.Background(), event.RepoKey, event)

	now := time.Now()
	event.ProcessedAt = &now
	event.Status = "processed"
	return s.eventDAO.Update(event)
}

func matchEventType(pattern, actual string) bool {
	if pattern == "" || pattern == "*" {
		return true
	}
	return pattern == actual
}
