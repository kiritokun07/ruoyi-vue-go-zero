package notice

import "ruoyi-vue-go-zero/app/admin/api/internal/logic/system"

import (
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func (s *catalogService) ListBulletin(r *types.AdminFilter) (*types.PageResult, error) {
	rows, n, e := s.BulletinRows(r, true)
	return &types.PageResult{StatusResponse: system.OK(), Rows: rows, Total: n}, e
}

func (s *catalogService) CreateBulletin(r *types.BulletinEntry) (*types.StatusResponse, error) {
	if err := s.ValidateBulletin(r, true); err != nil {
		return nil, err
	}
	return system.Saved(s.Models().CreateBulletin(s.Context(), r, s.Actor()))
}

func (s *catalogService) EditBulletin(r *types.BulletinEntry) (*types.StatusResponse, error) {
	if err := s.ValidateBulletin(r, false); err != nil {
		return nil, err
	}
	if e := system.RequireCount(s.Models().EditBulletinExists(s.Context(), r.NoticeId)); e != nil {
		return nil, e
	}
	return system.Saved(s.Models().EditBulletin(s.Context(), r, s.Actor()))
}

func (s *catalogService) DeleteBulletin(r *types.RecordIds) (*types.StatusResponse, error) {
	ids, e := system.ParseIDs(r.Ids)
	if e != nil {
		return nil, e
	}
	if e = s.BeforeDelete("sys_notice", ids); e != nil {
		return nil, e
	}
	return system.Saved(s.Models().DeleteBulletin(s.Context(), ids))
}
