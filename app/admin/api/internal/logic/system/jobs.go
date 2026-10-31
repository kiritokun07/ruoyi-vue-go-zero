package system

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorhill/cronexpr"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"ruoyi-vue-go-zero/app/admin/api/internal/model"
	"ruoyi-vue-go-zero/app/admin/api/internal/tasks"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

var weekdayNumber = regexp.MustCompile(`[0-9]+`)

func schedule(expr string) (*cronexpr.Expression, error) {
	fields := strings.Fields(expr)
	if len(fields) == 6 {
		fields = append(fields, "*")
	}
	if len(fields) != 7 {
		return nil, apierror.Business("无效的 Quartz Cron 表达式")
	}
	// Quartz numbers Sunday as 1. Preserve the occurrence after # and the stride after /.
	var day strings.Builder
	for i := 0; i < len(fields[5]); {
		token := weekdayNumber.FindStringIndex(fields[5][i:])
		if token == nil {
			day.WriteString(fields[5][i:])
			break
		}
		start, end := i+token[0], i+token[1]
		day.WriteString(fields[5][i:start])
		num := fields[5][start:end]
		if start == 0 || (fields[5][start-1] != '#' && fields[5][start-1] != '/') {
			n, _ := strconv.Atoi(num)
			if n < 1 || n > 7 {
				return nil, apierror.Business("星期取值应为1至7")
			}
			num = strconv.Itoa(n - 1)
		}
		day.WriteString(num)
		i = end
	}
	fields[5] = day.String()
	parsed, err := cronexpr.Parse(strings.Join(fields, " "))
	if err != nil {
		return nil, apierror.Business("无效的 Cron 表达式")
	}
	return parsed, nil
}

func (s *Session) ValidateSchedule(r *types.ScheduleEntry, create bool) error {
	if create {
		r.Status = "1"
	}
	if r.MisfirePolicy != "" && r.MisfirePolicy != "0" && r.MisfirePolicy != "1" && r.MisfirePolicy != "2" && r.MisfirePolicy != "3" {
		return apierror.Business("执行策略无效")
	}
	if r.JobName == "" || r.JobGroup == "" || !validStatus(r.Status) || !validStatus(r.Concurrent) {
		return apierror.Business("任务名称、组、并发或状态无效")
	}
	if e := tasks.Default.Validate(r.InvokeTarget); e != nil {
		return apierror.Business(e.Error())
	}
	expr, e := schedule(r.CronExpression)
	if e != nil {
		return e
	}
	if expr.Next(time.Now()).IsZero() {
		return apierror.Business("表达式没有未来的执行时间")
	}
	return nil
}

func (s *Session) JobStatus(r *types.ScheduleEntry) (*types.StatusResponse, error) {
	if !validStatus(r.Status) {
		return nil, apierror.Business("任务状态无效")
	}
	if e := s.require("SELECT COUNT(*) FROM sys_job WHERE job_id=?", r.JobId); e != nil {
		return nil, e
	}
	if e := s.exec("DELETE FROM go_job_trigger WHERE job_id=?", r.JobId); e != nil {
		return nil, e
	}
	return saved(s.exec("UPDATE sys_job SET status=? WHERE job_id=?", r.Status, r.JobId))
}

func executeJob(ctx context.Context, db sqlx.SqlConn, j types.ScheduleEntry) error {
	return db.TransactCtx(ctx, func(ctx context.Context, tx sqlx.Session) error {
		if j.Concurrent == "1" {
			var id int64
			if e := tx.QueryRowCtx(ctx, &id, "SELECT job_id FROM sys_job WHERE job_id=? FOR UPDATE", j.JobId); e != nil {
				return e
			}
		}
		return recordExecution(ctx, tx, j)

	})
}

func (s *Session) RunJob(r *types.ScheduleEntry) (*types.StatusResponse, error) {
	j, e := one[types.ScheduleEntry](s, "SELECT "+model.ScheduleColumns+" FROM sys_job WHERE job_id=? AND job_group=?", r.JobId, r.JobGroup)
	if e != nil {
		return nil, e
	}
	return saved(executeJob(s.ctx, s.services.DB, j))
}

// Claim each occurrence and persist its result in the same transaction as the next deadline.
func StartJobs(ctx context.Context, db sqlx.SqlConn) {
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case now := <-ticker.C:
				var ids []int64
				err := db.QueryRowsCtx(ctx, &ids, "SELECT j.job_id FROM sys_job j LEFT JOIN go_job_trigger t ON t.job_id=j.job_id WHERE j.status='0' AND (t.next_fire IS NULL OR t.next_fire<=?)", now.Unix())
				if err != nil {
					logx.Error(err)
					continue
				}
				for _, id := range ids {
					err := db.TransactCtx(ctx, func(ctx context.Context, tx sqlx.Session) error { return fireDueJob(ctx, tx, id, now) })
					if err != nil {
						logx.Error(err)
					}
				}
			}
		}
	}()
}

func fireDueJob(ctx context.Context, tx sqlx.Session, id int64, now time.Time) error {
	var j types.ScheduleEntry
	err := tx.QueryRowPartialCtx(ctx, &j, "SELECT "+model.ScheduleColumns+" FROM sys_job WHERE job_id=? FOR UPDATE", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil || j.Status != "0" {
		return err
	}
	expr, err := schedule(j.CronExpression)
	if err != nil {
		return err
	}
	var due int64
	err = tx.QueryRowCtx(ctx, &due, "SELECT next_fire FROM go_job_trigger WHERE job_id=? FOR UPDATE", id)
	next := expr.Next(now)
	switch {
	case errors.Is(err, sql.ErrNoRows):
	case err != nil:
		return err
	case due > now.Unix():
		return nil
	default:
		if now.Unix()-due <= 1 || j.MisfirePolicy != "3" {
			if err = recordExecution(ctx, tx, j); err != nil {
				return err
			}
		}
		if j.MisfirePolicy == "1" {
			next = expr.Next(time.Unix(due, 0))
		}
	}
	if next.IsZero() {
		_, err = tx.ExecCtx(ctx, "UPDATE sys_job SET status='1' WHERE job_id=?", id)
		return err
	}
	_, err = tx.ExecCtx(ctx, "INSERT INTO go_job_trigger(job_id,next_fire) VALUES(?,?) ON DUPLICATE KEY UPDATE next_fire=?", id, next.Unix(), next.Unix())
	return err
}

func setNextFire(j *types.ScheduleEntry) {
	if c, e := schedule(j.CronExpression); e == nil {
		if t := c.Next(time.Now()); !t.IsZero() {
			j.NextValidTime = t.Format(time.DateTime)
		}
	}
}

func recordExecution(ctx context.Context, tx sqlx.Session, j types.ScheduleEntry) error {
	message, e := tasks.Default.Invoke(ctx, j.InvokeTarget)
	status, detail := "0", ""
	if e != nil {
		status, detail = "1", e.Error()
	} else {
		logx.WithContext(ctx).Info(message)
	}
	_, err := tx.ExecCtx(ctx, "INSERT INTO sys_job_log(job_name,job_group,invoke_target,job_message,status,exception_info,create_time) VALUES(?,?,?,?,?,?,NOW())", j.JobName, j.JobGroup, j.InvokeTarget, message, status, detail)
	return err
}
