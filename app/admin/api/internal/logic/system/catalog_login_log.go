package system

import (
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func (s *Session) LoginEventRows(r *types.AdminFilter, paged bool) ([]types.LoginEventEntry, int64, error) {
	return s.models().LoginEventRows(s.ctx, r, paged)
}
