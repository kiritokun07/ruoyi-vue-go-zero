package model

import (
	"context"

	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

const JobLogColumns = `IFNULL(job_log_id,0) job_log_id,IFNULL(job_name,'') job_name,IFNULL(job_group,'') job_group,IFNULL(invoke_target,'') invoke_target,IFNULL(job_message,'') job_message,IFNULL(status,'') status,IFNULL(exception_info,'') exception_info,IFNULL(DATE_FORMAT(create_time,'%Y-%m-%d %H:%i:%s'),'') create_time`

func (m *AdminModel) JobLogRows(ctx context.Context, r *types.AdminFilter, paged bool) ([]types.JobLogEntry, int64, error) {
	return m.jobLogRows(ctx, r, paged, true)
}
func (m *AdminModel) JobLogExportBatch(ctx context.Context, r *types.AdminFilter) ([]types.JobLogEntry, error) {
	rows, _, err := m.jobLogRows(ctx, r, true, false)
	return rows, err
}
func (m *AdminModel) jobLogRows(ctx context.Context, r *types.AdminFilter, paged, count bool) ([]types.JobLogEntry, int64, error) {
	f := filter{}
	f.text("job_name", r.JobName)
	f.equal("job_group", r.JobGroup)
	f.equal("status", r.Status)
	f.dates("create_time", &types.ListFilter{BeginTime: r.BeginTime, EndTime: r.EndTime})
	from := " FROM sys_job_log" + f.where()
	var n int64
	if count {
		var err error
		n, err = queryOne[int64](ctx, m.db, "SELECT COUNT(*)"+from, f.args...)
		if err != nil {
			return nil, 0, err
		}
	}
	order := " ORDER BY job_log_id DESC"
	if paged {
		order += page(&types.ListFilter{PageNum: r.PageNum, PageSize: r.PageSize})
	}
	rows, err := queryRows[types.JobLogEntry](ctx, m.db, "SELECT "+JobLogColumns+from+order, f.args...)
	return rows, n, err
}

func (m *AdminModel) CleanJobLog(ctx context.Context) error {
	return m.exec(ctx, "DELETE FROM sys_job_log")
}

func (m *AdminModel) DeleteJobLog(ctx context.Context, ids []int64) error {
	return m.exec(ctx, "DELETE FROM sys_job_log WHERE job_log_id IN ("+numbers(ids)+")")
}

func (m *AdminModel) GetJobLog(ctx context.Context, id int64) (types.JobLogEntry, error) {
	return queryOne[types.JobLogEntry](ctx, m.db, "SELECT "+JobLogColumns+" FROM sys_job_log WHERE job_log_id=?", id)
}
