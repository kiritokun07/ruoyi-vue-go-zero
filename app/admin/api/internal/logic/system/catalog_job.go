package system

import (
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func (s *Session) ScheduleRows(r *types.AdminFilter, paged bool) ([]types.ScheduleEntry, int64, error) {
	rows, n, err := s.models().ScheduleRows(s.ctx, r, paged)
	for i := range rows {
		setNextFire(&rows[i])
	}
	return rows, n, err
}
