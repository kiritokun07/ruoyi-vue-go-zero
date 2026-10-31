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

type ListParameterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListParameterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListParameterLogic {
	return &ListParameterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListParameterLogic) ListParameter(req *types.AdminFilter) (resp *types.PageResult, err error) {
	return system.Read(l.ctx, l.svcCtx, "system:config:list", func(s *system.Session) (*types.PageResult, error) {
		return (&catalogService{Session: s}).ListParameter(req)
	})
}
