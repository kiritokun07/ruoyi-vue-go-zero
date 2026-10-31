// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package gen

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type GenPreviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenPreviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenPreviewLogic {
	return &GenPreviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenPreviewLogic) GenPreview(req *types.RecordId) (resp *types.Result, err error) {
	return system.Read(l.ctx, l.svcCtx, "tool:gen:preview", func(s *system.Session) (*types.Result, error) { return s.GenPreview(req) })
}
