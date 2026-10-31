package dict

import "ruoyi-vue-go-zero/app/admin/api/internal/logic/system"

import (
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func (s *catalogService) ListDictEntry(r *types.AdminFilter) (*types.PageResult, error) {
	rows, n, e := s.DictEntryRows(r, true)
	return &types.PageResult{StatusResponse: system.OK(), Rows: rows, Total: n}, e
}

func (s *catalogService) GetDictEntry(r *types.RecordId) (*types.Result, error) {
	v, e := s.Models().GetDictEntry(s.Context(), r.Id)
	return system.Result(v, e)
}

func (s *catalogService) CreateDictEntry(r *types.DictEntryEntry) (*types.StatusResponse, error) {
	if err := s.ValidateDictEntry(r, true); err != nil {
		return nil, err
	}
	return system.Saved(s.Models().CreateDictEntry(s.Context(), r, s.Actor()))
}

func (s *catalogService) EditDictEntry(r *types.DictEntryEntry) (*types.StatusResponse, error) {
	if err := s.ValidateDictEntry(r, false); err != nil {
		return nil, err
	}
	if e := system.RequireCount(s.Models().EditDictEntryExists(s.Context(), r.DictCode)); e != nil {
		return nil, e
	}
	return system.Saved(s.Models().EditDictEntry(s.Context(), r, s.Actor()))
}

func (s *catalogService) DeleteDictEntry(r *types.RecordIds) (*types.StatusResponse, error) {
	ids, e := system.ParseIDs(r.Ids)
	if e != nil {
		return nil, e
	}
	return system.Saved(s.Models().DeleteDictEntry(s.Context(), ids))
}
