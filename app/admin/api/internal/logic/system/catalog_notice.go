package system

import (
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func (s *Session) BulletinRows(r *types.AdminFilter, paged bool) ([]types.BulletinEntry, int64, error) {
	return s.models().BulletinRows(s.ctx, r, paged)
}
