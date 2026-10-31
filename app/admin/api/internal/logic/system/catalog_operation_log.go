package system

import (
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func (s *Session) OperationRows(r *types.AdminFilter, paged bool) ([]types.OperationEntry, int64, error) {
	return s.models().OperationRows(s.ctx, r, paged)
}
