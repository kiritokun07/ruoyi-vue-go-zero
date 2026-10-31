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

type SortMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSortMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SortMenusLogic {
	return &SortMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SortMenusLogic) SortMenus(req *types.SortInput) (resp *types.StatusResponse, err error) {
	return system.Write(l.ctx, l.svcCtx, "system:menu:edit", req, func(s *system.Session) (*types.StatusResponse, error) { return s.SortMenus(req) })
}
