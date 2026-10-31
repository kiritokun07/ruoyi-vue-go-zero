package model

import (
	"context"

	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

const ParameterColumns = `IFNULL(config_id,0) config_id,IFNULL(config_name,'') config_name,IFNULL(config_key,'') config_key,IFNULL(config_value,'') config_value,IFNULL(config_type,'') config_type,IFNULL(create_by,'') create_by,IFNULL(DATE_FORMAT(create_time,'%Y-%m-%d %H:%i:%s'),'') create_time,IFNULL(update_by,'') update_by,IFNULL(DATE_FORMAT(update_time,'%Y-%m-%d %H:%i:%s'),'') update_time,IFNULL(remark,'') remark`

func (m *AdminModel) ParameterRows(ctx context.Context, r *types.AdminFilter, paged bool) ([]types.ParameterEntry, int64, error) {
	return m.parameterRows(ctx, r, paged, true)
}
func (m *AdminModel) ParameterExportBatch(ctx context.Context, r *types.AdminFilter) ([]types.ParameterEntry, error) {
	rows, _, err := m.parameterRows(ctx, r, true, false)
	return rows, err
}
func (m *AdminModel) parameterRows(ctx context.Context, r *types.AdminFilter, paged, count bool) ([]types.ParameterEntry, int64, error) {
	f := filter{}
	f.text("config_name", r.ConfigName)
	f.text("config_key", r.ConfigKey)
	f.equal("config_type", r.ConfigType)
	f.dates("create_time", &types.ListFilter{BeginTime: r.BeginTime, EndTime: r.EndTime})
	from := " FROM sys_config" + f.where()
	var n int64
	if count {
		var err error
		n, err = queryOne[int64](ctx, m.db, "SELECT COUNT(*)"+from, f.args...)
		if err != nil {
			return nil, 0, err
		}
	}
	order := " ORDER BY config_id DESC"
	if paged {
		order += page(&types.ListFilter{PageNum: r.PageNum, PageSize: r.PageSize})
	}
	rows, err := queryRows[types.ParameterEntry](ctx, m.db, "SELECT "+ParameterColumns+from+order, f.args...)
	return rows, n, err
}

func (m *AdminModel) DeleteParameter(ctx context.Context, ids []int64) error {
	return m.exec(ctx, "DELETE FROM sys_config WHERE config_id IN ("+numbers(ids)+")")
}

func (m *AdminModel) EditParameter(ctx context.Context, r *types.ParameterEntry, actor string) error {
	return m.exec(ctx, "UPDATE sys_config SET config_name=?,config_key=?,config_value=?,config_type=?,remark=?,update_by=?,update_time=NOW() WHERE config_id=?", r.ConfigName, r.ConfigKey, r.ConfigValue, r.ConfigType, r.Remark, actor, r.ConfigId)
}

func (m *AdminModel) EditParameterConflict(ctx context.Context, key string, id int64) (int, error) {
	return queryOne[int](ctx, m.db, "SELECT COUNT(*) FROM sys_config WHERE config_key=? AND config_id<>?", key, id)
}

func (m *AdminModel) EditParameterExists(ctx context.Context, id int64) (int, error) {
	return queryOne[int](ctx, m.db, "SELECT COUNT(*) FROM sys_config WHERE config_id=?", id)
}

func (m *AdminModel) CreateParameter(ctx context.Context, r *types.ParameterEntry, actor string) error {
	return m.exec(ctx, "INSERT INTO sys_config(config_name,config_key,config_value,config_type,remark,create_by,create_time) VALUES(?,?,?,?,?,?,NOW())", r.ConfigName, r.ConfigKey, r.ConfigValue, r.ConfigType, r.Remark, actor)
}

func (m *AdminModel) CreateParameterConflict(ctx context.Context, key string, id int64) (int, error) {
	return queryOne[int](ctx, m.db, "SELECT COUNT(*) FROM sys_config WHERE config_key=? AND config_id<>?", key, id)
}

func (m *AdminModel) GetParameter(ctx context.Context, id int64) (types.ParameterEntry, error) {
	return queryOne[types.ParameterEntry](ctx, m.db, "SELECT "+ParameterColumns+" FROM sys_config WHERE config_id=?", id)
}

func (m *AdminModel) BuiltinParameterCount(ctx context.Context, ids []int64) (int, error) {
	return queryOne[int](ctx, m.db, "SELECT COUNT(*) FROM sys_config WHERE config_type='Y' AND config_id IN ("+numbers(ids)+")")
}
