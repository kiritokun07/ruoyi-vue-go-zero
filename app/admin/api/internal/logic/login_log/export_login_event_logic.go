// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package login_log

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type ExportLoginEventLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExportLoginEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportLoginEventLogic {
	return &ExportLoginEventLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExportLoginEventLogic) ExportLoginEvent(req *types.AdminFilter) (*system.File, error) {
	return system.Read(l.ctx, l.svcCtx, "monitor:logininfor:export", func(s *system.Session) (*system.File, error) { return s.ExportLoginEvent(req) })
}
