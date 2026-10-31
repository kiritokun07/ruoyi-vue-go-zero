package model

import (
	"context"

	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

const DictTypeColumns = `IFNULL(dict_id,0) dict_id,IFNULL(dict_name,'') dict_name,IFNULL(dict_type,'') dict_type,IFNULL(status,'') status,IFNULL(create_by,'') create_by,IFNULL(DATE_FORMAT(create_time,'%Y-%m-%d %H:%i:%s'),'') create_time,IFNULL(update_by,'') update_by,IFNULL(DATE_FORMAT(update_time,'%Y-%m-%d %H:%i:%s'),'') update_time,IFNULL(remark,'') remark`

func (m *AdminModel) DictTypeRows(ctx context.Context, r *types.AdminFilter, paged bool) ([]types.DictTypeEntry, int64, error) {
	return m.dictTypeRows(ctx, r, paged, true)
}
func (m *AdminModel) DictTypeExportBatch(ctx context.Context, r *types.AdminFilter) ([]types.DictTypeEntry, error) {
	rows, _, err := m.dictTypeRows(ctx, r, true, false)
	return rows, err
}
func (m *AdminModel) dictTypeRows(ctx context.Context, r *types.AdminFilter, paged, count bool) ([]types.DictTypeEntry, int64, error) {
	f := filter{}
	f.text("dict_name", r.DictName)
	f.text("dict_type", r.DictType)
	f.equal("status", r.Status)
	f.dates("create_time", &types.ListFilter{BeginTime: r.BeginTime, EndTime: r.EndTime})
	from := " FROM sys_dict_type" + f.where()
	var n int64
	if count {
		var err error
		n, err = queryOne[int64](ctx, m.db, "SELECT COUNT(*)"+from, f.args...)
		if err != nil {
			return nil, 0, err
		}
	}
	order := " ORDER BY dict_id DESC"
	if paged {
		order += page(&types.ListFilter{PageNum: r.PageNum, PageSize: r.PageSize})
	}
	rows, err := queryRows[types.DictTypeEntry](ctx, m.db, "SELECT "+DictTypeColumns+from+order, f.args...)
	return rows, n, err
}

func (m *AdminModel) DeleteDictType(ctx context.Context, ids []int64) error {
	return m.exec(ctx, "DELETE FROM sys_dict_type WHERE dict_id IN ("+numbers(ids)+")")
}

func (m *AdminModel) EditDictType(ctx context.Context, r *types.DictTypeEntry, actor string) error {
	return m.exec(ctx, "UPDATE sys_dict_type SET dict_name=?,dict_type=?,status=?,remark=?,update_by=?,update_time=NOW() WHERE dict_id=?", r.DictName, r.DictType, r.Status, r.Remark, actor, r.DictId)
}

func (m *AdminModel) EditDictTypeConflict(ctx context.Context, key string, id int64) (int, error) {
	return queryOne[int](ctx, m.db, "SELECT COUNT(*) FROM sys_dict_type WHERE dict_type=? AND dict_id<>?", key, id)
}

func (m *AdminModel) EditDictTypeExists(ctx context.Context, id int64) (int, error) {
	return queryOne[int](ctx, m.db, "SELECT COUNT(*) FROM sys_dict_type WHERE dict_id=?", id)
}

func (m *AdminModel) CreateDictType(ctx context.Context, r *types.DictTypeEntry, actor string) error {
	return m.exec(ctx, "INSERT INTO sys_dict_type(dict_name,dict_type,status,remark,create_by,create_time) VALUES(?,?,?,?,?,NOW())", r.DictName, r.DictType, r.Status, r.Remark, actor)
}

func (m *AdminModel) CreateDictTypeConflict(ctx context.Context, key string, id int64) (int, error) {
	return queryOne[int](ctx, m.db, "SELECT COUNT(*) FROM sys_dict_type WHERE dict_type=? AND dict_id<>?", key, id)
}

func (m *AdminModel) GetDictType(ctx context.Context, id int64) (types.DictTypeEntry, error) {
	return queryOne[types.DictTypeEntry](ctx, m.db, "SELECT "+DictTypeColumns+" FROM sys_dict_type WHERE dict_id=?", id)
}

func (m *AdminModel) RenameDictDataType(ctx context.Context, id int64, key string) error {
	return m.exec(ctx, "UPDATE sys_dict_data SET dict_type=? WHERE dict_type=(SELECT dict_type FROM sys_dict_type WHERE dict_id=?)", key, id)
}
func (m *AdminModel) DictTypeCount(ctx context.Context, key string) (int, error) {
	return queryOne[int](ctx, m.db, "SELECT COUNT(*) FROM sys_dict_type WHERE dict_type=?", key)
}
func (m *AdminModel) DictTypeReferences(ctx context.Context, ids []int64) (int, error) {
	return queryOne[int](ctx, m.db, "SELECT COUNT(*) FROM sys_dict_data d JOIN sys_dict_type t ON t.dict_type=d.dict_type WHERE t.dict_id IN ("+numbers(ids)+")")
}
func (m *AdminModel) DictOptions(ctx context.Context) ([]types.DictTypeEntry, error) {
	return queryRows[types.DictTypeEntry](ctx, m.db, "SELECT "+DictTypeColumns+" FROM sys_dict_type ORDER BY dict_id")
}
