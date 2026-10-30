package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SysUserModel = (*customSysUserModel)(nil)

type (
	// SysUserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysUserModel.
	SysUserModel interface {
		sysUserModel
		FindLoginUser(context.Context, string) (*LoginUser, error)
		RecordLogin(context.Context, int64, string) error
		withSession(session sqlx.Session) SysUserModel
	}

	customSysUserModel struct {
		*defaultSysUserModel
	}
)

// NewSysUserModel returns a model for the database table.
func NewSysUserModel(conn sqlx.SqlConn) SysUserModel {
	return &customSysUserModel{
		defaultSysUserModel: newSysUserModel(conn),
	}
}

func (m *customSysUserModel) withSession(session sqlx.Session) SysUserModel {
	return NewSysUserModel(sqlx.NewSqlConnFromSession(session))
}

type LoginUser struct {
	UserID   int64  `db:"user_id"`
	Password string `db:"password"`
	Status   string `db:"status"`
}

func (m *customSysUserModel) FindLoginUser(ctx context.Context, username string) (*LoginUser, error) {
	var user LoginUser
	err := m.conn.QueryRowCtx(ctx, &user, "SELECT user_id,password,status FROM sys_user WHERE user_name=? AND del_flag='0' LIMIT 1", username)
	return &user, err
}

func (m *customSysUserModel) RecordLogin(ctx context.Context, id int64, ip string) error {
	_, err := m.conn.ExecCtx(ctx, "UPDATE sys_user SET login_ip=?, login_date=NOW() WHERE user_id=?", ip, id)
	return err
}
