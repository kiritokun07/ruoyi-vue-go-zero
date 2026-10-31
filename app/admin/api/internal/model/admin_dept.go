package model

import (
	"context"

	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

const DepartmentColumns = `dept_id,parent_id,dept_name,order_num,ancestors,COALESCE(leader,'') leader,COALESCE(phone,'') phone,COALESCE(email,'') email,status,DATE_FORMAT(create_time,'%Y-%m-%d %H:%i:%s') create_time`

func (m *AdminModel) DepartmentRows(ctx context.Context, r *types.ListFilter, exclude int64, all bool, departments []int64) ([]types.DeptRecord, error) {
	f := filter{}
	f.add("del_flag='0'")
	if !all {
		f.add("dept_id IN (" + numbers(departments) + ")")
	}
	f.text("dept_name", r.DeptName)
	f.equal("status", r.Status)
	if exclude != 0 {
		f.add("dept_id<>? AND NOT FIND_IN_SET(?,ancestors)", exclude, exclude)
	}
	return queryRows[types.DeptRecord](ctx, m.db, "SELECT "+DepartmentColumns+" FROM sys_dept"+f.where()+" ORDER BY parent_id,order_num,dept_id", f.args...)
}
