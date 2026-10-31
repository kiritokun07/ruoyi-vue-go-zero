// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package online

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type OnlineUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOnlineUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OnlineUsersLogic {
	return &OnlineUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OnlineUsersLogic) OnlineUsers(req *types.AdminFilter) (resp *types.PageResult, err error) {
	return system.Read(l.ctx, l.svcCtx, "monitor:online:list", func(s *system.Session) (*types.PageResult, error) { return s.OnlineUsers(req) })
}
