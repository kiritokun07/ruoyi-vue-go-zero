package model

import (
	"context"

	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

const OperationColumns = `IFNULL(oper_id,0) oper_id,IFNULL(title,'') title,IFNULL(business_type,0) business_type,IFNULL(method,'') method,IFNULL(request_method,'') request_method,IFNULL(operator_type,0) operator_type,IFNULL(oper_name,'') oper_name,IFNULL(dept_name,'') dept_name,IFNULL(oper_url,'') oper_url,IFNULL(oper_ip,'') oper_ip,IFNULL(oper_location,'') oper_location,IFNULL(oper_param,'') oper_param,IFNULL(json_result,'') json_result,IFNULL(status,0) status,IFNULL(error_msg,'') error_msg,IFNULL(DATE_FORMAT(oper_time,'%Y-%m-%d %H:%i:%s'),'') oper_time,IFNULL(cost_time,0) cost_time`

func (m *AdminModel) OperationRows(ctx context.Context, r *types.AdminFilter, paged bool) ([]types.OperationEntry, int64, error) {
	return m.operationRows(ctx, r, paged, true)
}
func (m *AdminModel) OperationExportBatch(ctx context.Context, r *types.AdminFilter) ([]types.OperationEntry, error) {
	rows, _, err := m.operationRows(ctx, r, true, false)
	return rows, err
}
func (m *AdminModel) operationRows(ctx context.Context, r *types.AdminFilter, paged, count bool) ([]types.OperationEntry, int64, error) {
	f := filter{}
	f.text("title", r.Title)
	f.text("oper_name", r.OperName)
	f.text("oper_ip", r.OperIp)
	f.equal("status", r.Status)
	f.equal("business_type", r.BusinessType)
	f.dates("oper_time", &types.ListFilter{BeginTime: r.BeginTime, EndTime: r.EndTime})
	from := " FROM sys_oper_log" + f.where()
	var n int64
	if count {
		var err error
		n, err = queryOne[int64](ctx, m.db, "SELECT COUNT(*)"+from, f.args...)
		if err != nil {
			return nil, 0, err
		}
	}
	order := eventOrder(r, "oper_id", map[string]string{"operName": "oper_name", "operTime": "oper_time", "costTime": "cost_time"})
	if !count {
		order += ", oper_id DESC"
	}
	if paged {
		order += page(&types.ListFilter{PageNum: r.PageNum, PageSize: r.PageSize})
	}
	rows, err := queryRows[types.OperationEntry](ctx, m.db, "SELECT "+OperationColumns+from+order, f.args...)
	return rows, n, err
}

func (m *AdminModel) CleanOperation(ctx context.Context) error {
	return m.exec(ctx, "DELETE FROM sys_oper_log")
}

func (m *AdminModel) DeleteOperation(ctx context.Context, ids []int64) error {
	return m.exec(ctx, "DELETE FROM sys_oper_log WHERE oper_id IN ("+numbers(ids)+")")
}
