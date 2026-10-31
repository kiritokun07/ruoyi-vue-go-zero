// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package demo

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type DemoEditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDemoEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DemoEditLogic {
	return &DemoEditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DemoEditLogic) DemoEdit(req *types.DemoUser) (resp *types.Result, err error) {
	return system.Write(l.ctx, l.svcCtx, "", req, func(s *system.Session) (*types.Result, error) { return s.DemoEdit(req) })
}
