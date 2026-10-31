package login_log

import "ruoyi-vue-go-zero/app/admin/api/internal/logic/system"

import (
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func (s *catalogService) ListLoginEvent(r *types.AdminFilter) (*types.PageResult, error) {
	rows, n, e := s.LoginEventRows(r, true)
	return &types.PageResult{StatusResponse: system.OK(), Rows: rows, Total: n}, e
}

func (s *catalogService) DeleteLoginEvent(r *types.RecordIds) (*types.StatusResponse, error) {
	ids, e := system.ParseIDs(r.Ids)
	if e != nil {
		return nil, e
	}
	return system.Saved(s.Models().DeleteLoginEvent(s.Context(), ids))
}

func (s *catalogService) CleanLoginEvent() (*types.StatusResponse, error) {
	return system.Saved(s.Models().CleanLoginEvent(s.Context()))
}
