// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package role

import (
	"context"
	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"

	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleScopeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleScopeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleScopeLogic {
	return &RoleScopeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleScopeLogic) RoleScope(req *types.RoleRecord) (resp *types.StatusResponse, err error) {
	return system.Write(l.ctx, l.svcCtx, "system:role:edit", req, func(s *system.Session) (*types.StatusResponse, error) { return s.RoleScope(req) })
}
