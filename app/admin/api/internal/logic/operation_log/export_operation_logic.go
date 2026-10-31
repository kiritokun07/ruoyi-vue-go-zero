// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operation_log

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type ExportOperationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExportOperationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportOperationLogic {
	return &ExportOperationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExportOperationLogic) ExportOperation(req *types.AdminFilter) (*system.File, error) {
	return system.Read(l.ctx, l.svcCtx, "monitor:operlog:export", func(s *system.Session) (*system.File, error) { return s.ExportOperation(req) })
}
