// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package notice

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type ReadUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadUsersLogic {
	return &ReadUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadUsersLogic) ReadUsers(req *types.AdminFilter) (resp *types.PageResult, err error) {
	return system.Read(l.ctx, l.svcCtx, "system:notice:list", func(s *system.Session) (*types.PageResult, error) { return s.ReadUsers(req) })
}
