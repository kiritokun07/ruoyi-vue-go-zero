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

type GenGenerateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenGenerateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenGenerateLogic {
	return &GenGenerateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenGenerateLogic) GenGenerate(req *types.LookupKey) (resp *types.StatusResponse, err error) {
	return system.Read(l.ctx, l.svcCtx, "tool:gen:code", func(s *system.Session) (*types.StatusResponse, error) { return s.GenGenerate(req) })
}
