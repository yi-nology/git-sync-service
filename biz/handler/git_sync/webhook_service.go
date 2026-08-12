package git_sync

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/yi-nology/git-sync-service/biz/model/webhook"
	"github.com/yi-nology/git-sync-service/internal/converter"
	"github.com/yi-nology/git-sync-service/internal/pkg/response"
	"github.com/yi-nology/git-sync-service/internal/service"
	syncmodel "github.com/yi-nology/git-sync-service/sync/model"
)

func RuleList(ctx context.Context, c *app.RequestContext) {
	var req webhook.ListRulesReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.RepoKey == "" {
		response.BadRequest(c, "repoKey is required")
		return
	}

	list, err := GetSyncService().ListRules(ctx, req.RepoKey)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, &webhook.ListRulesResp{Rules: converter.ToRuleInfoList(list)})
}

func RuleGet(ctx context.Context, c *app.RequestContext) {
	var req webhook.GetRuleReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.ID <= 0 {
		response.BadRequest(c, "id is required")
		return
	}
	r, err := GetSyncService().GetRule(ctx, uint(req.ID))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	if r == nil {
		response.NotFound(c, "rule not found")
		return
	}
	response.Success(c, &webhook.GetRuleResp{Rule: converter.ToRuleInfo(r)})
}

func RuleCreate(ctx context.Context, c *app.RequestContext) {
	var req webhook.CreateRuleReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Name == "" || req.RepoKey == "" {
		response.BadRequest(c, "name and repoKey are required")
		return
	}

	var taskKeys []string
	if req.SyncTaskKeys != "" {
		for _, key := range strings.Split(req.SyncTaskKeys, ",") {
			key = strings.TrimSpace(key)
			if key != "" {
				taskKeys = append(taskKeys, key)
			}
		}
	}

	r, err := GetSyncService().CreateRule(ctx, &syncmodel.CreateRuleRequest{
		Name: req.Name, RepoKey: req.RepoKey, EventType: req.EventType,
		BranchPattern: req.BranchPattern, Action: req.Action,
		TaskKeys: taskKeys, MinInterval: int(req.MinInterval),
		Enabled: req.Enabled, Description: req.Description,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	recordAudit(ctx, c, "create", "rule", r.Name, "创建 webhook 规则 "+r.Name)
	response.Created(c, &webhook.CreateRuleResp{Rule: converter.ToRuleInfo(r)})
}

func RuleUpdate(ctx context.Context, c *app.RequestContext) {
	var req webhook.UpdateRuleReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.ID <= 0 {
		response.BadRequest(c, "id is required")
		return
	}

	var taskKeys []string
	if req.SyncTaskKeys != "" {
		for _, key := range strings.Split(req.SyncTaskKeys, ",") {
			key = strings.TrimSpace(key)
			if key != "" {
				taskKeys = append(taskKeys, key)
			}
		}
	}

	r, err := GetSyncService().UpdateRule(ctx, &syncmodel.UpdateRuleRequest{
		ID: uint(req.ID), Name: req.Name, EventType: req.EventType,
		BranchPattern: req.BranchPattern, Action: req.Action,
		TaskKeys: taskKeys, MinInterval: int(req.MinInterval),
		Enabled: req.Enabled, Description: req.Description,
	})
	if err != nil {
		if errors.Is(err, service.ErrRuleNotFound) {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	recordAudit(ctx, c, "update", "rule", r.Name, "更新 webhook 规则 "+r.Name)
	response.Success(c, &webhook.UpdateRuleResp{Rule: converter.ToRuleInfo(r)})
}

func RuleDelete(ctx context.Context, c *app.RequestContext) {
	var req webhook.DeleteRuleReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.ID <= 0 {
		response.BadRequest(c, "id is required")
		return
	}
	if err := GetSyncService().DeleteRule(ctx, uint(req.ID)); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	recordAudit(ctx, c, "delete", "rule", strconv.FormatUint(uint64(req.ID), 10), "删除 webhook 规则 #"+strconv.FormatUint(uint64(req.ID), 10))
	response.NoContent(c)
}

func ListEvents(ctx context.Context, c *app.RequestContext) {
	var req webhook.ListEventsReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.RepoKey == "" {
		response.BadRequest(c, "repoKey is required")
		return
	}

	limit := int(req.Limit)
	if limit <= 0 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.Query("offset"))
	if offset < 0 {
		offset = 0
	}
	events, _, err := GetSyncService().ListEvents(ctx, req.RepoKey, offset, limit)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, &webhook.ListEventsResp{Events: converter.ToEventInfoList(events)})
}

func RetryEvent(ctx context.Context, c *app.RequestContext) {
	var req webhook.RetryEventReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.ID <= 0 {
		response.BadRequest(c, "id is required")
		return
	}
	if err := GetSyncService().RetryEvent(ctx, uint(req.ID)); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	recordAudit(ctx, c, "retry", "event", strconv.FormatUint(uint64(req.ID), 10), "重试 webhook 事件 #"+strconv.FormatUint(uint64(req.ID), 10))
	response.Success(c, &webhook.RetryEventResp{Success: true, Message: "event retried"})
}
