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

type EditParameterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditParameterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditParameterLogic {
	return &EditParameterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditParameterLogic) EditParameter(req *types.ParameterEntry) (resp *types.StatusResponse, err error) {
	return system.Write(l.ctx, l.svcCtx, "system:config:edit", req, func(s *system.Session) (*types.StatusResponse, error) {
		return (&catalogService{Session: s}).EditParameter(req)
	})
}
