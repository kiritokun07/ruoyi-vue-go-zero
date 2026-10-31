package model

import (
	"context"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

const postColumns = `post_id,post_name,post_code,post_sort,status,IFNULL(remark,'') remark,DATE_FORMAT(create_time,'%Y-%m-%d %H:%i:%s') create_time`

func (m *AdminModel) PostRows(ctx context.Context, r *types.ListFilter, count bool) ([]types.PostRecord, int64, error) {
	f := filter{}
	f.text("post_code", r.PostCode)
	f.text("post_name", r.PostName)
	f.equal("status", r.Status)
	from := " FROM sys_post" + f.where()
	var n int64
	if count {
		var e error
		n, e = queryOne[int64](ctx, m.db, "SELECT COUNT(*)"+from, f.args...)
		if e != nil {
			return nil, 0, e
		}
	}
	rows, e := queryRows[types.PostRecord](ctx, m.db, "SELECT "+postColumns+from+" ORDER BY post_sort,post_id"+page(r), f.args...)
	return rows, n, e
}
