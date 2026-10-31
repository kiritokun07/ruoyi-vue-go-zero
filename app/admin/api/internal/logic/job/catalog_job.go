package job

import "ruoyi-vue-go-zero/app/admin/api/internal/logic/system"

import (
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func (s *catalogService) ListSchedule(r *types.AdminFilter) (*types.PageResult, error) {
	rows, n, e := s.ScheduleRows(r, true)
	return &types.PageResult{StatusResponse: system.OK(), Rows: rows, Total: n}, e
}

func (s *catalogService) GetSchedule(r *types.RecordId) (*types.Result, error) {
	v, e := s.Models().GetSchedule(s.Context(), r.Id)
	system.SetNextFire(&v)
	return system.Result(v, e)
}

func (s *catalogService) CreateSchedule(r *types.ScheduleEntry) (*types.StatusResponse, error) {
	if err := s.ValidateSchedule(r, true); err != nil {
		return nil, err
	}
	return system.Saved(s.Models().CreateSchedule(s.Context(), r, s.Actor()))
}

func (s *catalogService) EditSchedule(r *types.ScheduleEntry) (*types.StatusResponse, error) {
	if err := s.ValidateSchedule(r, false); err != nil {
		return nil, err
	}
	if e := system.RequireCount(s.Models().EditScheduleExists(s.Context(), r.JobId)); e != nil {
		return nil, e
	}
	return system.Saved(s.Models().EditSchedule(s.Context(), r, s.Actor()))
}

func (s *catalogService) DeleteSchedule(r *types.RecordIds) (*types.StatusResponse, error) {
	ids, e := system.ParseIDs(r.Ids)
	if e != nil {
		return nil, e
	}
	if e = s.Models().DeleteScheduleTriggers(s.Context(), ids); e != nil {
		return nil, e
	}
	return system.Saved(s.Models().DeleteSchedule(s.Context(), ids))
}
