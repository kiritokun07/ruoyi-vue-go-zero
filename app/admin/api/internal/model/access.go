package model

import (
	"context"
	"errors"
	"ruoyi-vue-go-zero/app/admin/api/internal/online"

	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

// AccessModel reads the user, role and menu projections used by the login APIs.
type AccessModel struct{ conn sqlx.SqlConn }

func NewAccessModel(conn sqlx.SqlConn) *AccessModel { return &AccessModel{conn} }
func (m *AccessModel) ActiveIdentity(ctx context.Context, id int64) (online.Identity, error) {
	var v online.Identity
	err := m.conn.QueryRowCtx(ctx, &v, "SELECT u.user_name,IFNULL(d.dept_name,'') dept_name FROM sys_user u LEFT JOIN sys_dept d ON d.dept_id=u.dept_id WHERE u.user_id=? AND u.status='0' AND u.del_flag='0'", id)
	if errors.Is(err, ErrNotFound) {
		return online.Identity{}, nil
	}
	return v, err
}

const UserProfileColumns = `u.user_id,u.dept_id,u.user_name,u.nick_name,IFNULL(u.email,'') email,IFNULL(u.phonenumber,'') phonenumber,IFNULL(u.sex,'') sex,
 IFNULL(u.avatar,'') avatar,u.status,u.del_flag,IFNULL(u.login_ip,'') login_ip,DATE_FORMAT(u.login_date,'%Y-%m-%d %H:%i:%s') login_date,
 DATE_FORMAT(u.pwd_update_date,'%Y-%m-%d %H:%i:%s') pwd_update_date,IFNULL(u.create_by,'') create_by,DATE_FORMAT(u.create_time,'%Y-%m-%d %H:%i:%s') create_time,IFNULL(u.remark,'') remark,u.user_id=1 admin`

func (m *AccessModel) Profile(ctx context.Context, id int64) (types.UserProfile, error) {
	return ReadUserProfile(ctx, m.conn, id, true)
}

func ReadUserProfile(ctx context.Context, db sqlx.Session, id int64, activeOnly bool) (types.UserProfile, error) {
	var u types.UserProfile
	query := "SELECT " + UserProfileColumns + " FROM sys_user u WHERE u.user_id=? AND u.del_flag='0'"
	if activeOnly {
		query += " AND u.status='0'"
	}
	err := db.QueryRowPartialCtx(ctx, &u, query, id)
	if err != nil {
		return u, err
	}
	if u.DeptId != nil {
		var d types.Department
		err = db.QueryRowCtx(ctx, &d, `SELECT dept_id,IFNULL(parent_id,0) parent_id,IFNULL(ancestors,'') ancestors,
  IFNULL(dept_name,'') dept_name,IFNULL(order_num,0) order_num,IFNULL(leader,'') leader,IFNULL(status,'') status FROM sys_dept WHERE dept_id=?`, *u.DeptId)
		if err == nil {
			u.Dept = &d
		} else if !errors.Is(err, ErrNotFound) {
			return u, err
		}
	}
	u.Roles = []types.Role{}
	err = db.QueryRowsPartialCtx(ctx, &u.Roles, `SELECT r.role_id,r.role_name,r.role_key,r.role_sort,IFNULL(r.data_scope,'') data_scope,r.status,r.role_id=1 admin
 FROM sys_role r JOIN sys_user_role ur ON ur.role_id=r.role_id
 WHERE ur.user_id=? AND r.del_flag='0' ORDER BY r.role_sort,r.role_id`, id)
	return u, err
}

const authorizedMenu = ` EXISTS (SELECT 1 FROM sys_role_menu rm
 JOIN sys_role r ON r.role_id=rm.role_id JOIN sys_user_role ur ON ur.role_id=r.role_id
 WHERE rm.menu_id=m.menu_id AND ur.user_id=? AND r.status='0' AND r.del_flag='0')`

type RolePermission struct {
	RoleID int64  `db:"role_id"`
	Keys   string `db:"perms"`
}

func (m *AccessModel) Permissions(ctx context.Context, id int64) ([]RolePermission, error) {
	var perms []RolePermission
	err := m.conn.QueryRowsCtx(ctx, &perms, `SELECT DISTINCT r.role_id,m.perms FROM sys_role r
 JOIN sys_user_role ur ON ur.role_id=r.role_id JOIN sys_role_menu rm ON rm.role_id=r.role_id
 JOIN sys_menu m ON m.menu_id=rm.menu_id
 WHERE ur.user_id=? AND r.role_id<>1 AND r.status='0' AND r.del_flag='0' AND m.status='0' AND m.perms<>''`, id)
	return perms, err
}

type Menu struct {
	ID        int64  `db:"menu_id"`
	Parent    int64  `db:"parent_id"`
	Title     string `db:"menu_name"`
	Path      string `db:"path"`
	Component string `db:"component"`
	Query     string `db:"route_query"`
	Name      string `db:"route_name"`
	Type      string `db:"menu_type"`
	Internal  bool   `db:"internal"`
	Hidden    bool   `db:"hidden"`
	NoCache   bool   `db:"no_cache"`
	Icon      string `db:"icon"`
}

func (m *AccessModel) Menus(ctx context.Context, id int64) ([]Menu, error) {
	query := `SELECT m.menu_id,m.parent_id,m.menu_name,IFNULL(m.path,'') path,IFNULL(m.component,'') component,
 IFNULL(m.query,'') AS route_query,IFNULL(m.route_name,'') route_name,m.menu_type,
 IFNULL(m.is_frame=1,0) AS internal,IFNULL(m.visible='1',0) hidden,IFNULL(m.is_cache=1,0) no_cache,IFNULL(m.icon,'') icon
 FROM sys_menu m WHERE m.status='0' AND m.menu_type IN ('M','C')`
	var args []any
	if id != 1 {
		query += " AND " + authorizedMenu
		args = append(args, id)
	}
	var menus []Menu
	err := m.conn.QueryRowsCtx(ctx, &menus, query+" ORDER BY m.parent_id,m.order_num,m.menu_id", args...)
	return menus, err
}
