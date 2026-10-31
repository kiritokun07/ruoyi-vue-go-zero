package job_log

import "ruoyi-vue-go-zero/app/admin/api/internal/logic/system"

import (
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func (s *catalogService) ListJobLog(r *types.AdminFilter) (*types.PageResult, error) {
	rows, n, e := s.JobLogRows(r, true)
	return &types.PageResult{StatusResponse: system.OK(), Rows: rows, Total: n}, e
}

func (s *catalogService) GetJobLog(r *types.RecordId) (*types.Result, error) {
	v, e := s.Models().GetJobLog(s.Context(), r.Id)
	return system.Result(v, e)
}

func (s *catalogService) DeleteJobLog(r *types.RecordIds) (*types.StatusResponse, error) {
	ids, e := system.ParseIDs(r.Ids)
	if e != nil {
		return nil, e
	}
	return system.Saved(s.Models().DeleteJobLog(s.Context(), ids))
}

func (s *catalogService) CleanJobLog() (*types.StatusResponse, error) {
	return system.Saved(s.Models().CleanJobLog(s.Context()))
}
