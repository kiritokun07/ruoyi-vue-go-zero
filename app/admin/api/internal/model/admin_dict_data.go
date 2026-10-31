package model

import (
	"context"

	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

const DictEntryColumns = `IFNULL(dict_code,0) dict_code,IFNULL(dict_sort,0) dict_sort,IFNULL(dict_label,'') dict_label,IFNULL(dict_value,'') dict_value,IFNULL(dict_type,'') dict_type,IFNULL(css_class,'') css_class,IFNULL(list_class,'') list_class,IFNULL(is_default,'') is_default,IFNULL(status,'') status,IFNULL(create_by,'') create_by,IFNULL(DATE_FORMAT(create_time,'%Y-%m-%d %H:%i:%s'),'') create_time,IFNULL(update_by,'') update_by,IFNULL(DATE_FORMAT(update_time,'%Y-%m-%d %H:%i:%s'),'') update_time,IFNULL(remark,'') remark`

func (m *AdminModel) DictEntryRows(ctx context.Context, r *types.AdminFilter, paged bool) ([]types.DictEntryEntry, int64, error) {
	return m.dictEntryRows(ctx, r, paged, true)
}
func (m *AdminModel) DictEntryExportBatch(ctx context.Context, r *types.AdminFilter) ([]types.DictEntryEntry, error) {
	rows, _, err := m.dictEntryRows(ctx, r, true, false)
	return rows, err
}
func (m *AdminModel) dictEntryRows(ctx context.Context, r *types.AdminFilter, paged, count bool) ([]types.DictEntryEntry, int64, error) {
	f := filter{}
	f.equal("dict_type", r.DictType)
	f.text("dict_label", r.DictLabel)
	f.equal("status", r.Status)
	f.dates("create_time", &types.ListFilter{BeginTime: r.BeginTime, EndTime: r.EndTime})
	from := " FROM sys_dict_data" + f.where()
	var n int64
	if count {
		var err error
		n, err = queryOne[int64](ctx, m.db, "SELECT COUNT(*)"+from, f.args...)
		if err != nil {
			return nil, 0, err
		}
	}
	order := " ORDER BY dict_code DESC"
	if paged {
		order += page(&types.ListFilter{PageNum: r.PageNum, PageSize: r.PageSize})
	}
	rows, err := queryRows[types.DictEntryEntry](ctx, m.db, "SELECT "+DictEntryColumns+from+order, f.args...)
	return rows, n, err
}

func (m *AdminModel) DeleteDictEntry(ctx context.Context, ids []int64) error {
	return m.exec(ctx, "DELETE FROM sys_dict_data WHERE dict_code IN ("+numbers(ids)+")")
}

func (m *AdminModel) EditDictEntry(ctx context.Context, r *types.DictEntryEntry, actor string) error {
	return m.exec(ctx, "UPDATE sys_dict_data SET dict_sort=?,dict_label=?,dict_value=?,dict_type=?,css_class=?,list_class=?,is_default=?,status=?,remark=?,update_by=?,update_time=NOW() WHERE dict_code=?", r.DictSort, r.DictLabel, r.DictValue, r.DictType, r.CssClass, r.ListClass, r.IsDefault, r.Status, r.Remark, actor, r.DictCode)
}

func (m *AdminModel) EditDictEntryExists(ctx context.Context, id int64) (int, error) {
	return queryOne[int](ctx, m.db, "SELECT COUNT(*) FROM sys_dict_data WHERE dict_code=?", id)
}

func (m *AdminModel) CreateDictEntry(ctx context.Context, r *types.DictEntryEntry, actor string) error {
	return m.exec(ctx, "INSERT INTO sys_dict_data(dict_sort,dict_label,dict_value,dict_type,css_class,list_class,is_default,status,remark,create_by,create_time) VALUES(?,?,?,?,?,?,?,?,?,?,NOW())", r.DictSort, r.DictLabel, r.DictValue, r.DictType, r.CssClass, r.ListClass, r.IsDefault, r.Status, r.Remark, actor)
}

func (m *AdminModel) GetDictEntry(ctx context.Context, id int64) (types.DictEntryEntry, error) {
	return queryOne[types.DictEntryEntry](ctx, m.db, "SELECT "+DictEntryColumns+" FROM sys_dict_data WHERE dict_code=?", id)
}
