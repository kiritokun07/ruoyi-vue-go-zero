package model

import (
	"context"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

const roleColumns = `r.role_id,r.role_name,r.role_key,r.role_sort,r.data_scope,r.status,IFNULL(r.remark,'') remark,
 r.menu_check_strictly,r.dept_check_strictly,DATE_FORMAT(r.create_time,'%Y-%m-%d %H:%i:%s') create_time,r.role_id=1 admin`

func (m *AdminModel) RoleRows(ctx context.Context, r *types.ListFilter, count bool, scope string, args []any) ([]types.RoleRecord, int64, error) {
	f := filter{}
	f.add("r.del_flag='0'")
	if scope != "" {
		f.add(scope, args...)
	}
	f.text("r.role_name", r.RoleName)
	f.text("r.role_key", r.RoleKey)
	f.equal("r.status", r.Status)
	f.dates("r.create_time", r)
	from := " FROM sys_role r" + f.where()
	var n int64
	if count {
		var e error
		n, e = queryOne[int64](ctx, m.db, "SELECT COUNT(*)"+from, f.args...)
		if e != nil {
			return nil, 0, e
		}
	}
	rows, e := queryRows[types.RoleRecord](ctx, m.db, "SELECT "+roleColumns+from+" ORDER BY r.role_sort,r.role_id"+page(r), f.args...)
	return rows, n, e
}
