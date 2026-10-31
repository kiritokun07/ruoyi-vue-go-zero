package system

import (
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func (s *Session) ParameterRows(r *types.AdminFilter, paged bool) ([]types.ParameterEntry, int64, error) {
	return s.models().ParameterRows(s.ctx, r, paged)
}
