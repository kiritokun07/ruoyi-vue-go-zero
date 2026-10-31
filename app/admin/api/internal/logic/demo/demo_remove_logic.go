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

type DemoRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDemoRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DemoRemoveLogic {
	return &DemoRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DemoRemoveLogic) DemoRemove(req *types.RecordId) (resp *types.Result, err error) {
	return system.Write(l.ctx, l.svcCtx, "", req, func(s *system.Session) (*types.Result, error) { return s.DemoRemove(req) })
}
