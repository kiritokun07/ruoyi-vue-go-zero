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

type AllocatedUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAllocatedUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllocatedUsersLogic {
	return &AllocatedUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AllocatedUsersLogic) AllocatedUsers(req *types.ListFilter) (resp *types.UserPage, err error) {
	return system.Read(l.ctx, l.svcCtx, "system:role:list", func(s *system.Session) (*types.UserPage, error) { return s.AllocatedUsers(req) })
}
