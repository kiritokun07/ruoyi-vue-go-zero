package operation_log

import "ruoyi-vue-go-zero/app/admin/api/internal/logic/system"

import (
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func (s *catalogService) ListOperation(r *types.AdminFilter) (*types.PageResult, error) {
	rows, n, e := s.OperationRows(r, true)
	return &types.PageResult{StatusResponse: system.OK(), Rows: rows, Total: n}, e
}

func (s *catalogService) DeleteOperation(r *types.RecordIds) (*types.StatusResponse, error) {
	ids, e := system.ParseIDs(r.Ids)
	if e != nil {
		return nil, e
	}
	return system.Saved(s.Models().DeleteOperation(s.Context(), ids))
}

func (s *catalogService) CleanOperation() (*types.StatusResponse, error) {
	return system.Saved(s.Models().CleanOperation(s.Context()))
}
