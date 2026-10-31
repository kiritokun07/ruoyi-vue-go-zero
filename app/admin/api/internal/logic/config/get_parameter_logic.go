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

type GetParameterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetParameterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetParameterLogic {
	return &GetParameterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetParameterLogic) GetParameter(req *types.RecordId) (resp *types.Result, err error) {
	return system.Read(l.ctx, l.svcCtx, "system:config:query", func(s *system.Session) (*types.Result, error) { return (&catalogService{Session: s}).GetParameter(req) })
}
