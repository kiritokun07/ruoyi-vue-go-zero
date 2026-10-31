package model

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

// AdminModel contains custom administration queries and can bind to the caller's transaction.
type AdminModel struct{ db sqlx.Session }

func NewAdminModel(db sqlx.Session) *AdminModel               { return &AdminModel{db: db} }
func (m *AdminModel) WithSession(db sqlx.Session) *AdminModel { return NewAdminModel(db) }

func queryRows[T any](ctx context.Context, db sqlx.Session, q string, args ...any) ([]T, error) {
	rows := []T{}
	err := db.QueryRowsPartialCtx(ctx, &rows, q, args...)
	return rows, err
}
func queryOne[T any](ctx context.Context, db sqlx.Session, q string, args ...any) (T, error) {
	var row T
	err := db.QueryRowPartialCtx(ctx, &row, q, args...)
	return row, err
}
func (m *AdminModel) exec(ctx context.Context, q string, args ...any) error {
	_, err := m.db.ExecCtx(ctx, q, args...)
	return err
}

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
func eventOrder(r *types.AdminFilter, def string, allowed map[string]string) string {
	field, ok := allowed[r.OrderByColumn]
	if !ok {
		return " ORDER BY " + def + " DESC"
	}
	direction := " DESC"
	if r.IsAsc == "asc" || r.IsAsc == "ascending" {
		direction = " ASC"
	}
	return " ORDER BY " + field + direction
}
