package system

import (
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func (s *Session) DictEntryRows(r *types.AdminFilter, paged bool) ([]types.DictEntryEntry, int64, error) {
	return s.models().DictEntryRows(s.ctx, r, paged)
}
