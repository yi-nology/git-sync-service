package service

import (
	"context"
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
	ctx := context.Background()

	if err := _ = svc.Record(ctx, &model.OperationLog{Action: "create", ResourceType: "repo", Resource: "创建仓库 repo-x", Actor: "admin"}); err != nil {
		t.Fatalf("record failed: %v", err)
	}

	// 缺省 status 应被填为 success
	if err := _ = svc.Record(ctx, &model.OperationLog{Action: "delete", ResourceType: "task", Resource: "删除任务 t1", Actor: "ops"}); err != nil {
		t.Fatalf("record failed: %v", err)
	}

	filter := &dao.OperationLogFilter{}
	got, total, err := svc.List(ctx, 0, 50, filter)
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
}

func TestOperationLogService_Stats(t *testing.T) {
	svc := setupOperationLogTestService(t)
	ctx := context.Background()

	// 添加一些测试数据
	_ = svc.Record(ctx, &model.OperationLog{Action: "create", ResourceType: "repo", Resource: "test1", Actor: "admin"})
	_ = svc.Record(ctx, &model.OperationLog{Action: "update", ResourceType: "task", Resource: "test2", Actor: "user"})

	today, week, total, err := svc.Stats(ctx)
	if err != nil {
		t.Fatalf("stats failed: %v", err)
	}

	if total < 2 {
		t.Errorf("expected total >= 2, got %d", total)
	}
	t.Logf("Stats: today=%d, week=%d, total=%d", today, week, total)
}
