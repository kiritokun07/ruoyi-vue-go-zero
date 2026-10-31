package model

import (
	"context"

	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

const ScheduleColumns = `IFNULL(job_id,0) job_id,IFNULL(job_name,'') job_name,IFNULL(job_group,'') job_group,IFNULL(invoke_target,'') invoke_target,IFNULL(cron_expression,'') cron_expression,IFNULL(misfire_policy,'') misfire_policy,IFNULL(concurrent,'') concurrent,IFNULL(status,'') status,IFNULL(create_by,'') create_by,IFNULL(DATE_FORMAT(create_time,'%Y-%m-%d %H:%i:%s'),'') create_time,IFNULL(update_by,'') update_by,IFNULL(DATE_FORMAT(update_time,'%Y-%m-%d %H:%i:%s'),'') update_time,IFNULL(remark,'') remark`

func (m *AdminModel) ScheduleRows(ctx context.Context, r *types.AdminFilter, paged bool) ([]types.ScheduleEntry, int64, error) {
	f := filter{}
	f.text("job_name", r.JobName)
	f.equal("job_group", r.JobGroup)
	f.equal("status", r.Status)
	f.dates("create_time", &types.ListFilter{BeginTime: r.BeginTime, EndTime: r.EndTime})
	from := " FROM sys_job" + f.where()
	n, err := queryOne[int64](ctx, m.db, "SELECT COUNT(*)"+from, f.args...)
	if err != nil {
		return nil, 0, err
	}
	order := " ORDER BY job_id DESC"
	if paged {
		order += page(&types.ListFilter{PageNum: r.PageNum, PageSize: r.PageSize})
	}
	rows, err := queryRows[types.ScheduleEntry](ctx, m.db, "SELECT "+ScheduleColumns+from+order, f.args...)
	return rows, n, err
}

func (m *AdminModel) DeleteSchedule(ctx context.Context, ids []int64) error {
	return m.exec(ctx, "DELETE FROM sys_job WHERE job_id IN ("+numbers(ids)+")")
}

func (m *AdminModel) DeleteScheduleTriggers(ctx context.Context, ids []int64) error {
	return m.exec(ctx, "DELETE FROM go_job_trigger WHERE job_id IN ("+numbers(ids)+")")
}

func (m *AdminModel) EditSchedule(ctx context.Context, r *types.ScheduleEntry, actor string) error {
	return m.exec(ctx, "UPDATE sys_job SET job_name=?,job_group=?,invoke_target=?,cron_expression=?,misfire_policy=?,concurrent=?,status=?,remark=?,update_by=?,update_time=NOW() WHERE job_id=?", r.JobName, r.JobGroup, r.InvokeTarget, r.CronExpression, r.MisfirePolicy, r.Concurrent, r.Status, r.Remark, actor, r.JobId)
}

func (m *AdminModel) EditScheduleExists(ctx context.Context, id int64) (int, error) {
	return queryOne[int](ctx, m.db, "SELECT COUNT(*) FROM sys_job WHERE job_id=?", id)
}

func (m *AdminModel) CreateSchedule(ctx context.Context, r *types.ScheduleEntry, actor string) error {
	return m.exec(ctx, "INSERT INTO sys_job(job_name,job_group,invoke_target,cron_expression,misfire_policy,concurrent,status,remark,create_by,create_time) VALUES(?,?,?,?,?,?,?,?,?,NOW())", r.JobName, r.JobGroup, r.InvokeTarget, r.CronExpression, r.MisfirePolicy, r.Concurrent, r.Status, r.Remark, actor)
}

func (m *AdminModel) GetSchedule(ctx context.Context, id int64) (types.ScheduleEntry, error) {
	return queryOne[types.ScheduleEntry](ctx, m.db, "SELECT "+ScheduleColumns+" FROM sys_job WHERE job_id=?", id)
}
