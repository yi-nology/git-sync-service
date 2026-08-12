package dao

import (
	"testing"
	"time"

	"github.com/yi-nology/git-sync-service/sync/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupOperationLogTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&model.OperationLog{}); err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}
	return db
}

func TestOperationLogDAO_CreateAndList(t *testing.T) {
	db := setupOperationLogTestDB(t)
	d := NewOperationLogDAO(db)

	entries := []*model.OperationLog{
		{Action: "create", ResourceType: "repo", ResourceKey: "repo-a", Resource: "创建仓库 repo-a", Actor: "admin", IP: "127.0.0.1", Status: "success"},
		{Action: "update", ResourceType: "task", ResourceKey: "task-1", Resource: "更新同步任务 task-1", Actor: "ops", IP: "10.0.0.1", Status: "success"},
		{Action: "delete", ResourceType: "rule", ResourceKey: "3", Resource: "删除 webhook 规则 #3", Actor: "admin", IP: "127.0.0.1", Status: "success"},
	}
	for _, e := range entries {
		if err := d.Create(e); err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// List all (newest first)
	got, total, err := d.List(DefaultPagination(0, 50), &OperationLogFilter{})
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if total != 3 || len(got) != 3 {
		t.Fatalf("expected 3 rows, got total=%d len=%d", total, len(got))
	}
	if got[0].ID != entries[2].ID {
		t.Fatalf("expected newest first, got id=%d", got[0].ID)
	}

	// Filter by action
	got, total, err = d.List(DefaultPagination(0, 50), &OperationLogFilter{Action: "create"})
	if err != nil {
		t.Fatalf("list by action failed: %v", err)
	}
	if total != 1 || len(got) != 1 || got[0].ResourceKey != "repo-a" {
		t.Fatalf("action filter wrong: total=%d got=%+v", total, got)
	}

	// Filter by actor
	got, total, err = d.List(DefaultPagination(0, 50), &OperationLogFilter{Actor: "admin"})
	if err != nil {
		t.Fatalf("list by actor failed: %v", err)
	}
	if total != 2 {
		t.Fatalf("actor filter expected 2, got %d", total)
	}

	// Search on resource/detail
	got, total, err = d.List(DefaultPagination(0, 50), &OperationLogFilter{Search: "webhook"})
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if total != 1 || got[0].ResourceKey != "3" {
		t.Fatalf("search expected 1 webhook row, got total=%d", total)
	}

	// Pagination
	got, total, err = d.List(DefaultPagination(0, 2), &OperationLogFilter{})
	if err != nil {
		t.Fatalf("paginated list failed: %v", err)
	}
	if total != 3 || len(got) != 2 {
		t.Fatalf("pagination expected total=3 len=2, got total=%d len=%d", total, len(got))
	}
}

func TestOperationLogDAO_DateFilterAndCount(t *testing.T) {
	db := setupOperationLogTestDB(t)
	d := NewOperationLogDAO(db)

	now := time.Now()
	old := now.AddDate(0, 0, -10)

	if err := db.Create(&model.OperationLog{Action: "create", ResourceType: "repo", Resource: "旧记录", Actor: "admin", Status: "success", CreatedAt: old}).Error; err != nil {
		t.Fatalf("seed old failed: %v", err)
	}
	if err := db.Create(&model.OperationLog{Action: "create", ResourceType: "repo", Resource: "新记录", Actor: "admin", Status: "success", CreatedAt: now}).Error; err != nil {
		t.Fatalf("seed new failed: %v", err)
	}

	today := now.Format("2006-01-02")
	got, total, err := d.List(DefaultPagination(0, 50), &OperationLogFilter{StartDate: today, EndDate: today})
	if err != nil {
		t.Fatalf("date filter list failed: %v", err)
	}
	// 当天记录的 created_at 秒级精度可能略晚于 "today 00:00:00"，但不应包含 10 天前的记录
	if total != 1 {
		t.Fatalf("date filter expected 1 (today only), got %d", total)
	}
	if got[0].Resource != "新记录" {
		t.Fatalf("date filter returned wrong row: %+v", got[0])
	}

	// CountSince: from start of today should be 1
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	c, err := d.CountSince(startOfToday)
	if err != nil {
		t.Fatalf("count since failed: %v", err)
	}
	if c != 1 {
		t.Fatalf("CountSince(today) expected 1, got %d", c)
	}

	// CountAll
	all, err := d.CountAll()
	if err != nil {
		t.Fatalf("count all failed: %v", err)
	}
	if all != 2 {
		t.Fatalf("CountAll expected 2, got %d", all)
	}
}
