package system

import (
	"bytes"
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func TestStreamingImportPartialSuccess(t *testing.T) {
	dsn := os.Getenv("RUOYI_SYSTEM_TEST_DSN")
	if dsn == "" {
		t.Skip("set RUOYI_SYSTEM_TEST_DSN")
	}
	ctx := context.Background()
	db := sqlx.NewMysql(dsn)
	name := "stream_" + uuid.NewString()[:8]
	file, err := workbook("导入", importHeaders, [][]string{{"103", name, "测试"}, {"103", name, "重复"}})
	if err != nil {
		t.Fatal(err)
	}
	prepared, err := prepareUsers(ctx, bytes.NewReader(readWorkbookTest(t, file)))
	if err != nil {
		t.Fatal(err)
	}
	defer prepared.Close()
	defer db.ExecCtx(ctx, "DELETE FROM sys_user WHERE user_name=?", name)
	_, err = importUsers(ctx, prepared, func(record importedUser) error {
		return db.TransactCtx(ctx, func(ctx context.Context, tx sqlx.Session) error {
			s := &Session{ctx: ctx, db: tx, user: 1, grants: map[string]bool{}}
			if err := s.authorize("system:user:import"); err != nil {
				return err
			}
			return s.importUser(record, "admin123", false)
		})
	})
	if err == nil || !strings.Contains(err.Error(), "第3行") {
		t.Fatalf("expected duplicate at row 3: %v", err)
	}
	var count int
	if err = db.QueryRowCtx(ctx, &count, "SELECT COUNT(*) FROM sys_user WHERE user_name=?", name); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatal("successful row was not preserved")
	}
}

func TestStreamingImportUpdatePreservesAssignments(t *testing.T) {
	dsn := os.Getenv("RUOYI_SYSTEM_TEST_DSN")
	if dsn == "" {
		t.Skip("set RUOYI_SYSTEM_TEST_DSN")
	}
	ctx := context.Background()
	rollback := errors.New("rollback fixtures")
	name := "stream_" + uuid.NewString()[:8]
	file, err := workbook("导入", importHeaders, [][]string{{"103", name, "修改后"}})
	if err != nil {
		t.Fatal(err)
	}
	prepared, err := prepareUsers(ctx, bytes.NewReader(readWorkbookTest(t, file)))
	if err != nil {
		t.Fatal(err)
	}
	defer prepared.Close()
	err = sqlx.NewMysql(dsn).TransactCtx(ctx, func(ctx context.Context, tx sqlx.Session) error {
		s := &Session{ctx: ctx, db: tx, user: 1, grants: map[string]bool{}}
		if err := s.authorize("system:user:import"); err != nil {
			return err
		}
		dept := int64(103)
		if _, err := s.CreateUser(&types.UserInput{UserName: name, NickName: "原名", Password: "admin123", Status: "0", DeptId: &dept, RoleIds: []int64{2}, PostIds: []int64{1}}); err != nil {
			return err
		}
		if _, err := s.importPreparedUsers(prepared, true); err != nil {
			return err
		}
		ids, err := s.models().UserIDsByName(ctx, name)
		if err != nil {
			return err
		}
		var nick string
		if err = tx.QueryRowCtx(ctx, &nick, "SELECT nick_name FROM sys_user WHERE user_id=?", ids[0]); err != nil {
			return err
		}
		if nick != "修改后" {
			t.Fatal("profile not updated")
		}
		for _, table := range []string{"sys_user_role", "sys_user_post"} {
			var n int
			if err = tx.QueryRowCtx(ctx, &n, "SELECT COUNT(*) FROM "+table+" WHERE user_id=?", ids[0]); err != nil {
				return err
			}
			if n != 1 {
				t.Fatal("assignments removed")
			}
		}
		return rollback
	})
	if !errors.Is(err, rollback) {
		t.Fatal(err)
	}
}
