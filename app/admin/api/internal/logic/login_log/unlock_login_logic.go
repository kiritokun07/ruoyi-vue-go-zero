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

type UnlockLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnlockLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnlockLoginLogic {
	return &UnlockLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnlockLoginLogic) UnlockLogin(req *types.LookupKey) (resp *types.StatusResponse, err error) {
	return system.Read(l.ctx, l.svcCtx, "monitor:logininfor:unlock", func(s *system.Session) (*types.StatusResponse, error) { return s.UnlockLogin(req) })
}
