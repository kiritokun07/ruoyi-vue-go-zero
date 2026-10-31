// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package config

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type RefreshConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshConfigLogic {
	return &RefreshConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshConfigLogic) RefreshConfig() (resp *types.StatusResponse, err error) {
	return system.Read(l.ctx, l.svcCtx, "system:config:remove", func(s *system.Session) (*types.StatusResponse, error) { return s.RefreshConfig() })
}
