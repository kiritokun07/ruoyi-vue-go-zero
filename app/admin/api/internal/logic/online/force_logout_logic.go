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

type ForceLogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewForceLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ForceLogoutLogic {
	return &ForceLogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ForceLogoutLogic) ForceLogout(req *types.LookupKey) (resp *types.StatusResponse, err error) {
	return system.Read(l.ctx, l.svcCtx, "monitor:online:forceLogout", func(s *system.Session) (*types.StatusResponse, error) { return s.ForceLogout(req) })
}
