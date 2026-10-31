package system

import (
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func (s *Session) DictTypeRows(r *types.AdminFilter, paged bool) ([]types.DictTypeEntry, int64, error) {
	return s.models().DictTypeRows(s.ctx, r, paged)
}
