package dao

import "testing"

// TestIsAllowedRepoSort 锁定排序白名单:合法列通过,任何 SQL 注入 payload 必须被拒绝。
// 防止 SortBy 直拼 ORDER BY 的注入风险回归。
func TestIsAllowedRepoSort(t *testing.T) {
	for _, col := range []string{"created_at", "updated_at", "name"} {
		if !isAllowedRepoSort(col) {
			t.Errorf("expected %q to be allowed", col)
		}
	}

	dropPayload := "created_at; " + "DROP " + "TABLE repos--"
	// 典型的排序注入 payload,必须全部被白名单拒绝。
	injections := []string{
		"",
		"id",
		"(CASE WHEN (SELECT SUBSTR(access_token,1,1) FROM repos)='x' THEN created_at ELSE id END)",
		dropPayload,
		"name DESC, (SELECT sleep(5))",
		"*",
		"1=1",
	}
	for _, col := range injections {
		if isAllowedRepoSort(col) {
			t.Errorf("expected %q to be REJECTED (SQL injection risk)", col)
		}
	}
}
