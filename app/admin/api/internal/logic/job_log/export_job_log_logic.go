// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package job_log

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type ExportJobLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExportJobLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportJobLogLogic {
	return &ExportJobLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExportJobLogLogic) ExportJobLog(req *types.AdminFilter) (*system.File, error) {
	return system.Read(l.ctx, l.svcCtx, "monitor:job:export", func(s *system.Session) (*system.File, error) { return s.ExportJobLog(req) })
}
