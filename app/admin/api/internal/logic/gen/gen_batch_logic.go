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

type GenBatchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenBatchLogic {
	return &GenBatchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenBatchLogic) GenBatch(req *types.TablesForm) (*system.File, error) {
	return system.Read(l.ctx, l.svcCtx, "tool:gen:code", func(s *system.Session) (*system.File, error) { return s.GenBatch(req) })
}
