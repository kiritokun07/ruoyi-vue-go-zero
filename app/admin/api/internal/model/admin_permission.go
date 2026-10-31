package model

import (
	"context"
)

type ActiveUser struct {
	Name string `db:"user_name"`
	Dept int64  `db:"dept_id"`
}
type PermissionGrant struct {
	Role  int64  `db:"role_id"`
	Scope string `db:"data_scope"`
	Keys  string `db:"perms"`
}

func (m *AdminModel) ActiveUser(ctx context.Context, id int64) (ActiveUser, error) {
	return queryOne[ActiveUser](ctx, m.db, "SELECT user_name,IFNULL(dept_id,0) dept_id FROM sys_user WHERE user_id=? AND status='0' AND del_flag='0'", id)
}
func (m *AdminModel) PermissionGrants(ctx context.Context, id int64) ([]PermissionGrant, error) {
	return queryRows[PermissionGrant](ctx, m.db, `SELECT r.role_id,r.data_scope,COALESCE(m.perms,'') perms FROM sys_user_role ur JOIN sys_role r ON r.role_id=ur.role_id
 LEFT JOIN sys_role_menu rm ON rm.role_id=r.role_id LEFT JOIN sys_menu m ON m.menu_id=rm.menu_id AND m.status='0'
 WHERE ur.user_id=? AND r.role_id<>1 AND r.status='0' AND r.del_flag='0'`, id)
}
