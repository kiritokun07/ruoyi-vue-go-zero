package model

import (
	"context"

	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func (m *AdminModel) userRows(ctx context.Context, r *types.ListFilter, scope, membership string, count bool) ([]types.UserProfile, int64, error) {
	f := filter{}
	f.add("u.del_flag='0'")
	f.add(scope)
	f.text("u.user_name", r.UserName)
	f.text("u.nick_name", r.NickName)
	f.text("u.phonenumber", r.Phonenumber)
	f.equal("u.status", r.Status)
	f.dates("u.create_time", r)
	if r.UserId != 0 {
		f.add("u.user_id=?", r.UserId)
	}
	if r.DeptId != 0 {
		f.add("(u.dept_id=? OR u.dept_id IN (SELECT dept_id FROM sys_dept WHERE FIND_IN_SET(?,ancestors)))", r.DeptId, r.DeptId)
	}
	if membership != "" {
		f.add(membership+" (SELECT 1 FROM sys_user_role ur WHERE ur.user_id=u.user_id AND ur.role_id=?)", r.RoleId)
	}
	from := " FROM sys_user u" + f.where()
	var total int64
	if count {
		var err error
		total, err = queryOne[int64](ctx, m.db, "SELECT COUNT(*)"+from, f.args...)
		if err != nil {
			return nil, 0, err
		}
	}
	rows, err := queryRows[types.UserProfile](ctx, m.db, "SELECT "+UserProfileColumns+from+" ORDER BY u.user_id"+page(r), f.args...)
	if err != nil {
		return nil, 0, err
	}
	ids := []int64{}
	for _, u := range rows {
		if u.DeptId != nil {
			ids = append(ids, *u.DeptId)
		}
	}
	if len(ids) > 0 {
		depts, e := queryRows[types.Department](ctx, m.db, "SELECT "+DepartmentColumns+" FROM sys_dept WHERE dept_id IN ("+numbers(ids)+")")
		if e != nil {
			return nil, 0, e
		}
		byID := map[int64]*types.Department{}
		for i := range depts {
			byID[depts[i].DeptId] = &depts[i]
		}
		for i := range rows {
			if rows[i].DeptId != nil {
				rows[i].Dept = byID[*rows[i].DeptId]
			}
		}
	}
	return rows, total, nil
}

func (m *AdminModel) UserRows(ctx context.Context, r *types.ListFilter, scope, membership string) ([]types.UserProfile, int64, error) {
	return m.userRows(ctx, r, scope, membership, true)
}
func (m *AdminModel) UserExportBatch(ctx context.Context, r *types.ListFilter, scope string) ([]types.UserProfile, error) {
	rows, _, err := m.userRows(ctx, r, scope, "", false)
	return rows, err
}

func (m *AdminModel) InitialPassword(ctx context.Context) (string, error) {
	return queryOne[string](ctx, m.db, "SELECT config_value FROM sys_config WHERE config_key='sys.user.initPassword' LIMIT 1")
}
func (m *AdminModel) UserIDsByName(ctx context.Context, name string) ([]int64, error) {
	return queryRows[int64](ctx, m.db, "SELECT user_id FROM sys_user WHERE user_name=? AND del_flag='0'", name)
}
func (m *AdminModel) UpdateImportedUser(ctx context.Context, id int64, r *types.UserInput, actor string) error {
	return m.exec(ctx, "UPDATE sys_user SET nick_name=?,email=?,phonenumber=?,sex=?,status=?,dept_id=COALESCE(?,dept_id),update_by=?,update_time=NOW() WHERE user_id=?", r.NickName, r.Email, r.Phonenumber, r.Sex, r.Status, r.DeptId, actor, id)
}
func (m *AdminModel) ProfilePhoneConflicts(ctx context.Context, id int64, phone string) (int, error) {
	return queryOne[int](ctx, m.db, "SELECT COUNT(*) FROM sys_user WHERE del_flag='0' AND phonenumber=? AND user_id<>?", phone, id)
}
func (m *AdminModel) ProfileEmailConflicts(ctx context.Context, id int64, email string) (int, error) {
	return queryOne[int](ctx, m.db, "SELECT COUNT(*) FROM sys_user WHERE del_flag='0' AND email=? AND user_id<>?", email, id)
}
