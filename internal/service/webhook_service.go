package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/yi-nology/git-sync-service/internal/dao"
	"github.com/yi-nology/git-sync-service/sync/model"
	"github.com/yi-nology/git-platform-sdk/pkg/branchfilter"
	"gorm.io/gorm"
)

// WebhookService handles webhook-related operations.
type WebhookService struct {
	ruleDAO  *dao.WebhookRuleDAO
	eventDAO *dao.WebhookEventDAO
	repoDAO  *dao.RepoDAO
}

// NewWebhookService creates a new WebhookService instance.
func NewWebhookService(ruleDAO *dao.WebhookRuleDAO, eventDAO *dao.WebhookEventDAO, repoDAO *dao.RepoDAO) *WebhookService {
	return &WebhookService{
		ruleDAO:  ruleDAO,
		eventDAO: eventDAO,
		repoDAO:  repoDAO,
	}
}

// ReceiveWebhook processes an incoming webhook request.
func (ws *WebhookService) ReceiveWebhook(ctx context.Context, repoKey string, req *http.Request, providerMgr interface{ GetByURL(url, token string) (WebhookParser, error) }) error {
	repo, err := ws.repoDAO.FindByKey(repoKey)
	if err != nil {
		return err
	}
	if repo == nil {
		return ErrRepoNotFound
	}

	prov, err := providerMgr.GetByURL(repo.CloneURL, repo.AccessToken)
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

	existing, err := ws.eventDAO.FindByEventID(event.ID)
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

	if err := ws.eventDAO.Create(whEvent); err != nil {
		return err
	}

	return nil
}

// ListRules returns webhook rules for a repository.
func (ws *WebhookService) ListRules(ctx context.Context, repoKey string) ([]*model.WebhookRule, error) {
	return ws.ruleDAO.FindByRepoKey(repoKey)
}

// GetRule returns a webhook rule by ID.
func (ws *WebhookService) GetRule(ctx context.Context, id uint) (*model.WebhookRule, error) {
	return ws.ruleDAO.FindByID(id)
}

// CreateRule creates a new webhook rule.
func (ws *WebhookService) CreateRule(ctx context.Context, req *model.CreateRuleRequest) (*model.WebhookRule, error) {
	rule := &model.WebhookRule{
		Name:          req.Name,
		RepoKey:       req.RepoKey,
		EventType:     req.EventType,
		BranchPattern: req.BranchPattern,
		Action:        req.Action,
		MinInterval:   req.MinInterval,
		Enabled:       req.Enabled,
		Description:   req.Description,
	}
	rule.SetTaskKeys(req.TaskKeys)

	if err := ws.ruleDAO.Create(rule); err != nil {
		return nil, err
	}

	return rule, nil
}

// UpdateRule updates an existing webhook rule.
func (ws *WebhookService) UpdateRule(ctx context.Context, req *model.UpdateRuleRequest) (*model.WebhookRule, error) {
	rule, err := ws.ruleDAO.FindByID(req.ID)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return nil, ErrRuleNotFound
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
	if req.TaskKeys != nil {
		rule.SetTaskKeys(req.TaskKeys)
	}
	if req.MinInterval > 0 {
		rule.MinInterval = req.MinInterval
	}
	rule.Enabled = req.Enabled
	rule.Description = req.Description

	if err := ws.ruleDAO.Update(rule); err != nil {
		return nil, err
	}

	return rule, nil
}

// DeleteRule deletes a webhook rule by ID.
func (ws *WebhookService) DeleteRule(ctx context.Context, id uint) error {
	return ws.ruleDAO.Delete(id)
}

// ListEvents returns webhook events for a repository.
func (ws *WebhookService) ListEvents(ctx context.Context, repoKey string, offset, limit int) ([]*model.WebhookEvent, int64, error) {
	page := dao.DefaultPagination(offset, limit)
	return ws.eventDAO.FindByRepoKey(repoKey, page)
}

// RetryEvent retries processing a webhook event.
func (ws *WebhookService) RetryEvent(ctx context.Context, eventID uint, applyRulesFn func(ctx context.Context, repoKey string, event *model.WebhookEvent)) error {
	event, err := ws.eventDAO.FindByID(eventID)
	if err != nil {
		return err
	}
	if event == nil {
		return ErrEventNotFound
	}

	go applyRulesFn(context.Background(), event.RepoKey, event)

	now := time.Now()
	event.ProcessedAt = &now
	event.Status = "processed"
	return ws.eventDAO.Update(event)
}

// FindRulesByRepoKey returns webhook rules for a repository (internal use).
func (ws *WebhookService) FindRulesByRepoKey(repoKey string) ([]*model.WebhookRule, error) {
	return ws.ruleDAO.FindByRepoKey(repoKey)
}

// FindEventByEventID returns a webhook event by event ID (internal use).
func (ws *WebhookService) FindEventByEventID(eventID string) (*model.WebhookEvent, error) {
	return ws.eventDAO.FindByEventID(eventID)
}

// FindEventByID returns a webhook event by ID (internal use).
func (ws *WebhookService) FindEventByID(id uint) (*model.WebhookEvent, error) {
	return ws.eventDAO.FindByID(id)
}

// ApplyRules applies webhook rules to an event.
func (ws *WebhookService) ApplyRules(ctx context.Context, repoKey string, event *model.WebhookEvent, lastTriggerTime *sync.Map, runTaskFn func(ctx context.Context, taskKey, trigger string) error) {
	rules, err := ws.ruleDAO.FindByRepoKey(repoKey)
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

		if rule.Action == "sync" {
			taskKeys := rule.GetTaskKeys()
			for _, taskKey := range taskKeys {
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
				if err := runTaskFn(ctx, taskKey, "webhook"); err != nil {
					slog.Error("run task failed", "taskKey", taskKey, "error", err)
				}
			}
		}
	}
}

// CreateWebhookEvent creates a new webhook event.
func (ws *WebhookService) CreateWebhookEvent(event *model.WebhookEvent) error {
	return ws.eventDAO.Create(event)
}

// CleanupOldEvents removes webhook events older than the specified duration.
func (ws *WebhookService) CleanupOldEvents(maxAge time.Duration) (int64, error) {
	return ws.eventDAO.CleanupOlderThan(maxAge)
}

// WebhookParser is an interface for webhook parsing operations.
type WebhookParser interface {
	ValidateWebhookSignature(req *http.Request, secret string) error
	ParseWebhookEvent(req *http.Request, secret string) (*NormalizedEvent, error)
}

// NormalizedEvent represents a normalized webhook event from the SDK.
type NormalizedEvent struct {
	ID         string
	Type       string
	Source     string
	Actor      *EventActor
	Branch     string
	CommitSHA  string
	RawPayload []byte
}

// EventActor represents the actor of a webhook event.
type EventActor struct {
	Name string
}
