package system

import (
	"context"
	"errors"
	"os"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

// The fixtures live within one transaction, including when the test fails.
func TestRoleDataScope(t *testing.T) {
	dsn := os.Getenv("RUOYI_SYSTEM_TEST_DSN")
	if dsn == "" {
		t.Skip("set RUOYI_SYSTEM_TEST_DSN to a development database initialized with the project SQL")
	}
	ctx := context.Background()
	rolledBack := errors.New("rollback test fixtures")
	err := sqlx.NewMysql(dsn).TransactCtx(ctx, func(ctx context.Context, tx sqlx.Session) error {
		admin := &Session{ctx: ctx, db: tx, user: 1, grants: map[string]bool{}}
		if err := admin.authorize("system:user:list"); err != nil {
			return err
		}
		suffix := uuid.NewString()[:8]
		// Parent, child and sibling give different results for department and subtree scopes.
		dept := func(parent int64) int64 {
			name := uuid.NewString()
			_, err := admin.CreateDept(&types.DeptRecord{ParentId: parent, DeptName: name[:20], Status: "0"})
			if err != nil {
				t.Fatal(err)
			}
			id, err := one[int64](admin, "SELECT dept_id FROM sys_dept WHERE dept_name=?", name[:20])
			if err != nil {
				t.Fatal(err)
			}
			return id
		}
		parent := dept(100)
		child := dept(parent)
		sibling := dept(100)
		getMenu := func(key string) int64 {
			id, err := one[int64](admin, "SELECT menu_id FROM sys_menu WHERE perms=? LIMIT 1", key)
			if err != nil {
				t.Fatal(err)
			}
			return id
		}
		role := func(key string, menu int64) int64 {
			_, e := admin.CreateRole(&types.RoleRecord{RoleName: key, RoleKey: key, MenuIds: []int64{menu}, Status: "0"})
			if e != nil {
				t.Fatal(e)
			}
			id, e := one[int64](admin, "SELECT role_id FROM sys_role WHERE role_key=?", key)
			if e != nil {
				t.Fatal(e)
			}
			return id
		}
		listRole := role("list"+suffix, getMenu("system:user:list"))
		detailRole := role("detail"+suffix, getMenu("system:user:query"))
		makeUser := func(d int64, roles ...int64) int64 {
			name := "t" + uuid.NewString()[:8]
			_, err := admin.CreateUser(&types.UserInput{UserName: name, NickName: name, DeptId: &d, Password: "test-pass", Status: "0", RoleIds: roles})
			if err != nil {
				t.Fatal(err)
			}
			id, err := one[int64](admin, "SELECT user_id FROM sys_user WHERE user_name=?", name)
			if err != nil {
				t.Fatal(err)
			}
			return id
		}
		actor := makeUser(parent, listRole, detailRole)
		peer := makeUser(parent)
		sub := makeUser(child)
		other := makeUser(sibling)
		for _, tc := range []struct {
			scope   string
			visible []int64
		}{
			{"1", []int64{actor, peer, sub, other}}, {"2", []int64{other}}, {"3", []int64{actor, peer}}, {"4", []int64{actor, peer, sub}}, {"5", []int64{actor}},
		} {
			_, err := admin.RoleScope(&types.RoleRecord{RoleId: listRole, DataScope: tc.scope, DeptIds: []int64{sibling}})
			if err != nil {
				return err
			}
			session := &Session{ctx: ctx, db: tx, user: actor, grants: map[string]bool{}}
			if err := session.authorize("system:user:list"); err != nil {
				return err
			}
			for _, id := range []int64{actor, peer, sub, other} {
				e := session.userAllowed(id, false)
				if (e == nil) != slices.Contains(tc.visible, id) || (e != nil && !errors.Is(e, apierror.Forbidden)) {
					t.Fatalf("scope %s, user %d: %v", tc.scope, id, e)
				}
			}
		}
		denied := &Session{ctx: ctx, db: tx, user: actor, grants: map[string]bool{}}
		if e := denied.authorize("system:user:remove"); !errors.Is(e, apierror.Forbidden) {
			t.Fatalf("missing operation permission: %v", e)
		}
		return rolledBack
	})
	if !errors.Is(err, rolledBack) {
		t.Fatal(err)
	}
}
