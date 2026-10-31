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

type DemoSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDemoSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DemoSaveLogic {
	return &DemoSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DemoSaveLogic) DemoSave(req *types.DemoUser) (resp *types.Result, err error) {
	return system.Write(l.ctx, l.svcCtx, "", req, func(s *system.Session) (*types.Result, error) { return s.DemoSave(req) })
}
