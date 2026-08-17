package model

import (
	"testing"
	"time"
)

func TestOperationLog_TableName(t *testing.T) {
	l := OperationLog{}
	expected := "operation_logs"

	if l.TableName() != expected {
		t.Errorf("expected table name '%s', got '%s'", expected, l.TableName())
	}
}

func TestOperationLog_Fields(t *testing.T) {
	now := time.Now()
	l := OperationLog{
		ID:           1,
		Actor:        "admin",
		Action:       "create",
		ResourceType: "repo",
		ResourceKey:  "test-repo",
		Resource:     "创建仓库 test-repo",
		Detail:       "详细信息",
		IP:           "192.168.1.1",
		Status:       "success",
		CreatedAt:    now,
	}

	if l.ID != 1 {
		t.Errorf("expected ID 1, got %d", l.ID)
	}

	if l.Actor != "admin" {
		t.Errorf("expected actor 'admin', got '%s'", l.Actor)
	}

	if l.Action != "create" {
		t.Errorf("expected action 'create', got '%s'", l.Action)
	}

	if l.ResourceType != "repo" {
		t.Errorf("expected resource type 'repo', got '%s'", l.ResourceType)
	}

	if l.ResourceKey != "test-repo" {
		t.Errorf("expected resource key 'test-repo', got '%s'", l.ResourceKey)
	}

	if l.Resource != "创建仓库 test-repo" {
		t.Errorf("expected resource '创建仓库 test-repo', got '%s'", l.Resource)
	}

	if l.Detail != "详细信息" {
		t.Errorf("expected detail '详细信息', got '%s'", l.Detail)
	}

	if l.IP != "192.168.1.1" {
		t.Errorf("expected IP '192.168.1.1', got '%s'", l.IP)
	}

	if l.Status != "success" {
		t.Errorf("expected status 'success', got '%s'", l.Status)
	}
}
