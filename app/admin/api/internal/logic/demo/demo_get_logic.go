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

type DemoGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDemoGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DemoGetLogic {
	return &DemoGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DemoGetLogic) DemoGet(req *types.RecordId) (resp *types.Result, err error) {
	return system.Read(l.ctx, l.svcCtx, "", func(s *system.Session) (*types.Result, error) { return s.DemoGet(req) })
}
