package system

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"ruoyi-vue-go-zero/app/admin/api/internal/middleware"
	"ruoyi-vue-go-zero/app/admin/api/internal/model"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type Session struct {
	services    *svc.ServiceContext
	ctx         context.Context
	db          sqlx.Session
	user        int64
	name        string
	department  int64
	all, self   bool
	departments []int64
	grants      map[string]bool
	roleIDs     []int64
}

func Read[T any](ctx context.Context, services *svc.ServiceContext, permission string, action func(*Session) (*T, error)) (*T, error) {
	return withSession(ctx, services, permission, false, action)
}

func Write[T any](ctx context.Context, services *svc.ServiceContext, permission string, input any, action func(*Session) (*T, error)) (*T, error) {
	if err := services.Validator.StructCtx(ctx, input); err != nil {
		return nil, apierror.Business("提交的数据不符合字段要求")
	}
	return withSession(ctx, services, permission, true, action)
}

func withSession[T any](ctx context.Context, services *svc.ServiceContext, permission string, write bool, action func(*Session) (*T, error)) (*T, error) {
	var result *T
	run := func(db sqlx.Session) error {
		s := &Session{services: services, ctx: ctx, db: db, user: middleware.Current(ctx).UserID, grants: map[string]bool{}}
		if err := s.authorize(permission); err != nil {
			return err
		}
		var err error
		result, err = action(s)
		return err
	}
	var err error
	if write {
		err = services.DB.TransactCtx(ctx, func(_ context.Context, tx sqlx.Session) error { return run(tx) })
	} else {
		err = run(services.DB)
	}
	if err == nil && write {
		prefix := ""
		switch {
		case strings.HasPrefix(permission, "system:dict:"):
			prefix = dictPrefix
		case strings.HasPrefix(permission, "system:config:"):
			prefix = configPrefix
		}
		if prefix != "" {
			_, cacheErr := services.Redis.IncrCtx(ctx, "ruoyi-go:generation:"+prefix)
			if cacheErr != nil {
				logx.WithContext(ctx).Error(cacheErr)
			}
		}
	}
	if errors.Is(err, sql.ErrNoRows) {
		err = apierror.Business("记录不存在")
	}
	return result, err
}

func (s *Session) authorize(permission string) error {
	u, err := s.models().ActiveUser(s.ctx, s.user)
	if errors.Is(err, sql.ErrNoRows) {
		return apierror.Unauthorized
	}
	if err != nil {
		return err
	}
	s.name, s.department = u.Name, u.Dept
	if s.user == 1 {
		s.all = true
		return nil
	}
	rows, err := s.models().PermissionGrants(s.ctx, s.user)
	if err != nil {
		return err
	}
	scopes := map[int64]string{}
	for _, r := range rows {
		if !slices.Contains(s.roleIDs, r.Role) {
			s.roleIDs = append(s.roleIDs, r.Role)
		}
		if permission == "" {
			scopes[r.Role] = r.Scope
		}
		for _, k := range strings.Split(r.Keys, ",") {
			k = strings.TrimSpace(k)
			s.grants[k] = true
			if k == permission || k == "*:*:*" {
				scopes[r.Role] = r.Scope
			}
		}
	}
	if permission != "" && len(scopes) == 0 {
		return apierror.Forbidden
	}
	for role, scope := range scopes {
		switch scope {
		case "1":
			s.all = true
			return nil
		case "2":
			ids, e := many[int64](s, "SELECT dept_id FROM sys_role_dept WHERE role_id=?", role)
			if e != nil {
				return e
			}
			s.departments = append(s.departments, ids...)
		case "3":
			s.departments = append(s.departments, s.department)
		case "4":
			ids, e := many[int64](s, "SELECT dept_id FROM sys_dept WHERE del_flag='0' AND (dept_id=? OR FIND_IN_SET(?,ancestors))", s.department, s.department)
			if e != nil {
				return e
			}
			s.departments = append(s.departments, ids...)
		case "5":
			s.self = true
		}
	}
	return nil
}

func many[T any](s *Session, q string, args ...any) ([]T, error) {
	rows := []T{}
	err := s.db.QueryRowsPartialCtx(s.ctx, &rows, q, args...)
	return rows, err
}

func one[T any](s *Session, q string, args ...any) (T, error) {
	var value T
	err := s.db.QueryRowPartialCtx(s.ctx, &value, q, args...)
	return value, err
}

func (s *Session) exec(q string, args ...any) error {
	_, err := s.db.ExecCtx(s.ctx, q, args...)
	return err
}

func (s *Session) rejectIf(q, msg string, args ...any) error {
	n, err := one[int](s, q, args...)
	if err == nil && n > 0 {
		return apierror.Business(msg)
	}
	return err
}

func (s *Session) require(q string, args ...any) error {
	n, err := one[int](s, q, args...)
	if err == nil && n == 0 {
		return apierror.Forbidden
	}
	return err
}

func (s *Session) deptAllowed(id int64) bool { return s.all || slices.Contains(s.departments, id) }
func (s *Session) userScope() string {
	if s.all {
		return "1=1"
	}
	expr := "u.dept_id IN (" + numbers(s.departments) + ")"
	if s.self {
		expr += " OR u.user_id=" + strconv.FormatInt(s.user, 10)
	}
	return "(" + expr + ")"
}

func numbers(ids []int64) string {
	v := make([]string, len(ids))
	for i, id := range ids {
		v[i] = strconv.FormatInt(id, 10)
	}
	if len(v) == 0 {
		return "NULL"
	}
	return strings.Join(v, ",")
}

func parseIDs(raw string) ([]int64, error) {
	ids := []int64{}
	if raw == "" {
		return ids, nil
	}
	for _, v := range strings.Split(raw, ",") {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil || id <= 0 {
			return nil, apierror.Business("编号格式错误")
		}
		if !slices.Contains(ids, id) {
			ids = append(ids, id)
		}
	}
	return ids, nil
}

func (s *Session) userAllowed(id int64, write bool) error {
	if write && id == 1 {
		return apierror.Business("不允许操作超级管理员用户")
	}
	return s.require("SELECT COUNT(*) FROM sys_user u WHERE u.del_flag='0' AND u.user_id=? AND "+s.userScope(), id)
}

func ok() types.StatusResponse                       { return types.StatusResponse{Code: 200, Msg: "操作成功"} }
func saved(err error) (*types.StatusResponse, error) { r := ok(); return &r, err }

type filter struct {
	clauses []string
	args    []any
}

func (f *filter) add(clause string, args ...any) {
	f.clauses = append(f.clauses, clause)
	f.args = append(f.args, args...)
}

func (f *filter) text(column, value string) {
	if value != "" {
		f.add(column+" LIKE ?", "%"+value+"%")
	}
}

func (f *filter) equal(column, value string) {
	if value != "" {
		f.add(column+"=?", value)
	}
}

func (f *filter) dates(col string, r *types.ListFilter) {
	if r.BeginTime != "" {
		f.add(col+">=?", r.BeginTime)
	}
	if r.EndTime != "" {
		f.add(col+"<DATE_ADD(?,INTERVAL 1 DAY)", r.EndTime)
	}
}

func (f filter) where() string {
	if len(f.clauses) == 0 {
		return ""
	}
	return " WHERE " + strings.Join(f.clauses, " AND ")
}

func page(r *types.ListFilter) string {
	size := r.PageSize
	if size <= 0 {
		size = 10
	}
	if size > 1000 {
		size = 1000
	}
	n := r.PageNum
	if n < 1 {
		n = 1
	}
	return fmt.Sprintf(" LIMIT %d OFFSET %d", size, (int64(n)-1)*int64(size))
}

func (s *Session) replaceLinks(table, ownerCol, itemCol string, owner int64, ids []int64) error {
	if err := s.exec("DELETE FROM "+table+" WHERE "+ownerCol+"=?", owner); err != nil {
		return err
	}
	for _, id := range ids {
		if err := s.exec("INSERT INTO "+table+" ("+ownerCol+","+itemCol+") VALUES (?,?)", owner, id); err != nil {
			return err
		}
	}
	return nil
}

// models binds all custom queries to the current transaction, including test sessions.
func (s *Session) models() *model.AdminModel {
	if s.services != nil && s.services.Admin != nil {
		return s.services.Admin.WithSession(s.db)
	}
	return model.NewAdminModel(s.db)
}
func requireCount(n int, err error) error {
	if err == nil && n == 0 {
		return apierror.Forbidden
	}
	return err
}

func (s *Session) Models() *model.AdminModel             { return s.models() }
func (s *Session) Context() context.Context              { return s.ctx }
func (s *Session) Actor() string                         { return s.name }
func OK() types.StatusResponse                           { return ok() }
func Saved(err error) (*types.StatusResponse, error)     { return saved(err) }
func Result(value any, err error) (*types.Result, error) { return result(value, err) }
func ParseIDs(raw string) ([]int64, error)               { return parseIDs(raw) }
func RequireCount(n int, err error) error                { return requireCount(n, err) }
func SetNextFire(job *types.ScheduleEntry)               { setNextFire(job) }
