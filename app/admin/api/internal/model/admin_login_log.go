package model

import (
	"context"

	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

const LoginEventColumns = `IFNULL(info_id,0) info_id,IFNULL(user_name,'') user_name,IFNULL(ipaddr,'') ipaddr,IFNULL(login_location,'') login_location,IFNULL(browser,'') browser,IFNULL(os,'') os,IFNULL(status,'') status,IFNULL(msg,'') msg,IFNULL(DATE_FORMAT(login_time,'%Y-%m-%d %H:%i:%s'),'') login_time`

func (m *AdminModel) LoginEventRows(ctx context.Context, r *types.AdminFilter, paged bool) ([]types.LoginEventEntry, int64, error) {
	f := filter{}
	f.text("ipaddr", r.Ipaddr)
	f.text("user_name", r.UserName)
	f.equal("status", r.Status)
	f.dates("login_time", &types.ListFilter{BeginTime: r.BeginTime, EndTime: r.EndTime})
	from := " FROM sys_logininfor" + f.where()
	n, err := queryOne[int64](ctx, m.db, "SELECT COUNT(*)"+from, f.args...)
	if err != nil {
		return nil, 0, err
	}
	order := eventOrder(r, "info_id", map[string]string{"userName": "user_name", "loginTime": "login_time"})
	if paged {
		order += page(&types.ListFilter{PageNum: r.PageNum, PageSize: r.PageSize})
	}
	rows, err := queryRows[types.LoginEventEntry](ctx, m.db, "SELECT "+LoginEventColumns+from+order, f.args...)
	return rows, n, err
}

func (m *AdminModel) CleanLoginEvent(ctx context.Context) error {
	return m.exec(ctx, "DELETE FROM sys_logininfor")
}

func (m *AdminModel) DeleteLoginEvent(ctx context.Context, ids []int64) error {
	return m.exec(ctx, "DELETE FROM sys_logininfor WHERE info_id IN ("+numbers(ids)+")")
}
