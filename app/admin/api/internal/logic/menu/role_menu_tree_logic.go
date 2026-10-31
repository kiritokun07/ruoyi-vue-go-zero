// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package menu

import (
	"context"
	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"

	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleMenuTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleMenuTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleMenuTreeLogic {
	return &RoleMenuTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleMenuTreeLogic) RoleMenuTree(req *types.RecordId) (resp *types.MenuTreeResult, err error) {
	return system.Read(l.ctx, l.svcCtx, "system:role:query", func(s *system.Session) (*types.MenuTreeResult, error) { return s.RoleMenuTree(req) })
}
