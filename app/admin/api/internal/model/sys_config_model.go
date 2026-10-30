package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SysConfigModel = (*customSysConfigModel)(nil)

type (
	// SysConfigModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysConfigModel.
	SysConfigModel interface {
		sysConfigModel
		Value(context.Context, string) (string, error)
		withSession(session sqlx.Session) SysConfigModel
	}

	customSysConfigModel struct {
		*defaultSysConfigModel
	}
)

// NewSysConfigModel returns a model for the database table.
func NewSysConfigModel(conn sqlx.SqlConn) SysConfigModel {
	return &customSysConfigModel{
		defaultSysConfigModel: newSysConfigModel(conn),
	}
}

func (m *customSysConfigModel) withSession(session sqlx.Session) SysConfigModel {
	return NewSysConfigModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customSysConfigModel) Value(ctx context.Context, key string) (string, error) {
	var value string
	err := m.conn.QueryRowCtx(ctx, &value, "SELECT COALESCE(config_value,'') FROM sys_config WHERE config_key=? ORDER BY config_id LIMIT 1", key)
	if err == ErrNotFound {
		return "", nil
	}
	return value, err
}
