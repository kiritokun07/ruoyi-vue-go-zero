package dict

import "ruoyi-vue-go-zero/app/admin/api/internal/logic/system"

import (
	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func (s *catalogService) ListDictType(r *types.AdminFilter) (*types.PageResult, error) {
	rows, n, e := s.DictTypeRows(r, true)
	return &types.PageResult{StatusResponse: system.OK(), Rows: rows, Total: n}, e
}

func (s *catalogService) GetDictType(r *types.RecordId) (*types.Result, error) {
	v, e := s.Models().GetDictType(s.Context(), r.Id)
	return system.Result(v, e)
}

func (s *catalogService) CreateDictType(r *types.DictTypeEntry) (*types.StatusResponse, error) {
	if err := s.ValidateDictType(r, true); err != nil {
		return nil, err
	}
	id := int64(0)
	if n, e := s.Models().CreateDictTypeConflict(s.Context(), r.DictType, id); e != nil || n > 0 {
		if e != nil {
			return nil, e
		}
		return nil, apierror.Business("名称或键已存在")
	}
	return system.Saved(s.Models().CreateDictType(s.Context(), r, s.Actor()))
}

func (s *catalogService) EditDictType(r *types.DictTypeEntry) (*types.StatusResponse, error) {
	if err := s.ValidateDictType(r, false); err != nil {
		return nil, err
	}
	if e := system.RequireCount(s.Models().EditDictTypeExists(s.Context(), r.DictId)); e != nil {
		return nil, e
	}
	id := r.DictId
	if n, e := s.Models().EditDictTypeConflict(s.Context(), r.DictType, id); e != nil || n > 0 {
		if e != nil {
			return nil, e
		}
		return nil, apierror.Business("名称或键已存在")
	}
	if err := s.Models().RenameDictDataType(s.Context(), r.DictId, r.DictType); err != nil {
		return nil, err
	}
	return system.Saved(s.Models().EditDictType(s.Context(), r, s.Actor()))
}

func (s *catalogService) DeleteDictType(r *types.RecordIds) (*types.StatusResponse, error) {
	ids, e := system.ParseIDs(r.Ids)
	if e != nil {
		return nil, e
	}
	if e = s.BeforeDelete("sys_dict_type", ids); e != nil {
		return nil, e
	}
	return system.Saved(s.Models().DeleteDictType(s.Context(), ids))
}
