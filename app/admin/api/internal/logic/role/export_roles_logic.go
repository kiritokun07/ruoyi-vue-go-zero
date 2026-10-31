// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package role

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type ExportRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExportRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportRolesLogic {
	return &ExportRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExportRolesLogic) ExportRoles(req *types.ListFilter) (*system.File, error) {
	return system.Read(l.ctx, l.svcCtx, "system:role:export", func(s *system.Session) (*system.File, error) { return s.ExportRoles(req) })
}
