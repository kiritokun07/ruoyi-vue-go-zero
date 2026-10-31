// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package job

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type ExportScheduleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExportScheduleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportScheduleLogic {
	return &ExportScheduleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExportScheduleLogic) ExportSchedule(req *types.AdminFilter) (*system.File, error) {
	return system.Read(l.ctx, l.svcCtx, "monitor:job:export", func(s *system.Session) (*system.File, error) { return s.ExportSchedule(req) })
}
