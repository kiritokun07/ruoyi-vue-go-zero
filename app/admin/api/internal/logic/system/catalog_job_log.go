package system

import (
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func (s *Session) JobLogRows(r *types.AdminFilter, paged bool) ([]types.JobLogEntry, int64, error) {
	return s.models().JobLogRows(s.ctx, r, paged)
}
