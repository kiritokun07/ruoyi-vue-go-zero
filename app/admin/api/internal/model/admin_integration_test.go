package model

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

// Check reads and writes share the caller's transaction, then roll back every fixture.
func TestAdminTransactionQueries(t *testing.T) {
	dsn := os.Getenv("RUOYI_SYSTEM_TEST_DSN")
	if dsn == "" {
		t.Skip("RUOYI_SYSTEM_TEST_DSN is not configured")
	}
	ctx := context.Background()
	conn := sqlx.NewMysql(dsn)
	root := NewAdminModel(conn)
	rollback := errors.New("rollback fixtures")
	key := "refactor_" + uuid.NewString()
	var configID int64
	err := conn.TransactCtx(ctx, func(ctx context.Context, tx sqlx.Session) error {
		m := root.WithSession(tx)
		config := &types.ParameterEntry{ConfigName: key, ConfigKey: key, ConfigValue: "00123", ConfigType: "N"}
		if err := m.CreateParameter(ctx, config, "test"); err != nil {
			return err
		}
		rows, total, err := m.ParameterRows(ctx, &types.AdminFilter{ConfigKey: key}, true)
		if err != nil {
			return err
		}
		if total != 1 || len(rows) != 1 || rows[0].ConfigValue != "00123" {
			t.Fatal("configuration not visible in transaction")
		}
		configID = rows[0].ConfigId
		config.ConfigId = configID
		config.ConfigValue = "updated"
		if err := m.EditParameter(ctx, config, "test"); err != nil {
			return err
		}
		got, err := m.GetParameter(ctx, configID)
		if err != nil {
			return err
		}
		if got.ConfigValue != "updated" {
			t.Fatal("configuration update lost")
		}
		dict := &types.DictTypeEntry{DictName: key, DictType: key, Status: "0"}
		if err := m.CreateDictType(ctx, dict, "test"); err != nil {
			return err
		}
		dictionaries, total, err := m.DictTypeRows(ctx, &types.AdminFilter{DictType: key}, false)
		if err != nil {
			return err
		}
		if total != 1 || len(dictionaries) != 1 {
			t.Fatal("dictionary not visible")
		}
		entry := &types.DictEntryEntry{DictLabel: "测试", DictValue: "0", DictType: key, Status: "0", IsDefault: "N"}
		if err := m.CreateDictEntry(ctx, entry, "test"); err != nil {
			return err
		}
		entries, total, err := m.DictEntryRows(ctx, &types.AdminFilter{DictType: key}, true)
		if err != nil {
			return err
		}
		if total != 1 || len(entries) != 1 {
			t.Fatal("dictionary data not visible")
		}
		entry.DictCode = entries[0].DictCode
		entry.DictLabel = "修改"
		if err := m.EditDictEntry(ctx, entry, "test"); err != nil {
			return err
		}
		if err := m.CreateBulletin(ctx, &types.BulletinEntry{NoticeTitle: key, NoticeType: "1", NoticeContent: "测试", Status: "0"}, "test"); err != nil {
			return err
		}
		notices, total, err := m.BulletinRows(ctx, &types.AdminFilter{NoticeTitle: key}, true)
		if err != nil {
			return err
		}
		if total != 1 || len(notices) != 1 {
			t.Fatal("notice not visible")
		}
		job := &types.ScheduleEntry{JobName: key, JobGroup: "DEFAULT", InvokeTarget: "ryTask.ryNoParams", CronExpression: "0 0 0 * * ?", Status: "1", MisfirePolicy: "3", Concurrent: "1"}
		if err := m.CreateSchedule(ctx, job, "test"); err != nil {
			return err
		}
		jobs, total, err := m.ScheduleRows(ctx, &types.AdminFilter{JobName: key}, true)
		if err != nil {
			return err
		}
		if total != 1 || len(jobs) != 1 {
			t.Fatal("job not visible")
		}
		job.JobId = jobs[0].JobId
		if err := m.EditSchedule(ctx, job, "test"); err != nil {
			return err
		}
		if err := m.DeleteScheduleTriggers(ctx, []int64{job.JobId}); err != nil {
			return err
		}
		if err := m.DeleteSchedule(ctx, []int64{job.JobId}); err != nil {
			return err
		}
		if _, _, err := m.JobLogRows(ctx, &types.AdminFilter{JobName: key}, true); err != nil {
			return err
		}
		if _, _, err := m.LoginEventRows(ctx, &types.AdminFilter{UserName: key}, true); err != nil {
			return err
		}
		if _, _, err := m.OperationRows(ctx, &types.AdminFilter{Title: key}, true); err != nil {
			return err
		}
		return rollback
	})
	if !errors.Is(err, rollback) {
		t.Fatal(err)
	}
	if _, err := root.GetParameter(ctx, configID); err == nil {
		t.Fatal("transaction fixture escaped rollback")
	}
}
