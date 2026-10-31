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

type ListMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMenuLogic {
	return &ListMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMenuLogic) ListMenu(req *types.ListFilter) (resp *types.MenuList, err error) {
	return system.Read(l.ctx, l.svcCtx, "system:menu:list", func(s *system.Session) (*types.MenuList, error) { return s.ListMenu(req) })
}
