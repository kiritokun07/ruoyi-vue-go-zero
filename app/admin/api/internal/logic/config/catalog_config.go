package config

import "ruoyi-vue-go-zero/app/admin/api/internal/logic/system"

import (
	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func (s *catalogService) ListParameter(r *types.AdminFilter) (*types.PageResult, error) {
	rows, n, e := s.ParameterRows(r, true)
	return &types.PageResult{StatusResponse: system.OK(), Rows: rows, Total: n}, e
}

func (s *catalogService) GetParameter(r *types.RecordId) (*types.Result, error) {
	v, e := s.Models().GetParameter(s.Context(), r.Id)
	return system.Result(v, e)
}

func (s *catalogService) CreateParameter(r *types.ParameterEntry) (*types.StatusResponse, error) {
	if err := s.ValidateParameter(r, true); err != nil {
		return nil, err
	}
	id := int64(0)
	if n, e := s.Models().CreateParameterConflict(s.Context(), r.ConfigKey, id); e != nil || n > 0 {
		if e != nil {
			return nil, e
		}
		return nil, apierror.Business("名称或键已存在")
	}
	return system.Saved(s.Models().CreateParameter(s.Context(), r, s.Actor()))
}

func (s *catalogService) EditParameter(r *types.ParameterEntry) (*types.StatusResponse, error) {
	if err := s.ValidateParameter(r, false); err != nil {
		return nil, err
	}
	if e := system.RequireCount(s.Models().EditParameterExists(s.Context(), r.ConfigId)); e != nil {
		return nil, e
	}
	id := r.ConfigId
	if n, e := s.Models().EditParameterConflict(s.Context(), r.ConfigKey, id); e != nil || n > 0 {
		if e != nil {
			return nil, e
		}
		return nil, apierror.Business("名称或键已存在")
	}
	return system.Saved(s.Models().EditParameter(s.Context(), r, s.Actor()))
}

func (s *catalogService) DeleteParameter(r *types.RecordIds) (*types.StatusResponse, error) {
	ids, e := system.ParseIDs(r.Ids)
	if e != nil {
		return nil, e
	}
	if e = s.BeforeDelete("sys_config", ids); e != nil {
		return nil, e
	}
	return system.Saved(s.Models().DeleteParameter(s.Context(), ids))
}
