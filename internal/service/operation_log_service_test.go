package service

import (
	"testing"

	"github.com/yi-nology/git-sync-service/internal/dao"
	"github.com/yi-nology/git-sync-service/sync/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupOperationLogTestService(t *testing.T) *OperationLogService {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&model.OperationLog{}); err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}
	return NewOperationLogService(dao.NewOperationLogDAO(db))
}

func TestOperationLogService_RecordAndList(t *testing.T) {
	svc := setupOperationLogTestService(t)

	if err := svc.Record(nil, &model.OperationLog{Action: "create", ResourceType: "repo", Resource: "创建仓库 repo-x", Actor: "admin"}); err != nil {
		t.Fatalf("record failed: %v", err)
	}

	// 缺省 status 应被填为 success
	if err := svc.Record(nil, &model.OperationLog{Action: "delete", ResourceType: "task", Resource: "删除任务 t1", Actor: "ops"}); err != nil {
		t.Fatalf("record failed: %v", err)
	}

	got, total, err := svc.List(nil, 0, 50, dao.OperationLogFilter{})
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if total != 2 || len(got) != 2 {
		t.Fatalf("expected 2, got total=%d len=%d", total, len(got))
	}
	// newest first → 删除任务 t1
	if got[0].Resource != "删除任务 t1" {
		t.Fatalf("expected newest first, got %+v", got[0])
	}
	if got[0].Status != model.StatusSuccess {
		t.Fatalf("expected default status success, got %q", got[0].Status)
	}

	// filter by actor
	_, total, err = svc.List(nil, 0, 50, dao.OperationLogFilter{Actor: "ops"})
	if err != nil {
		t.Fatalf("list filter failed: %v", err)
	}
	if total != 1 {
		t.Fatalf("actor filter expected 1, got %d", total)
	}
}

func TestOperationLogService_Stats(t *testing.T) {
	svc := setupOperationLogTestService(t)

	for i := 0; i < 3; i++ {
		if err := svc.Record(nil, &model.OperationLog{Action: "create", ResourceType: "repo", Resource: "创建仓库", Actor: "admin"}); err != nil {
			t.Fatalf("record failed: %v", err)
		}
	}

	today, week, total, err := svc.Stats(nil)
	if err != nil {
		t.Fatalf("stats failed: %v", err)
	}
	if total != 3 {
		t.Fatalf("total expected 3, got %d", total)
	}
	if today != 3 {
		t.Fatalf("today expected 3, got %d", today)
	}
	if week != 3 {
		t.Fatalf("week expected 3, got %d", week)
	}
}
