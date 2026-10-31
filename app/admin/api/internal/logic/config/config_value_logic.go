// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package config

import (
	"context"
	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"

	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigValueLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfigValueLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigValueLogic {
	return &ConfigValueLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigValueLogic) ConfigValue(req *types.LookupKey) (resp *types.StatusResponse, err error) {
	return system.Read(l.ctx, l.svcCtx, "", func(s *system.Session) (*types.StatusResponse, error) { return s.ConfigValue(req) })
}
