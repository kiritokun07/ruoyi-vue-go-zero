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

type GenColumnsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenColumnsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenColumnsLogic {
	return &GenColumnsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenColumnsLogic) GenColumns(req *types.RecordId) (resp *types.PageResult, err error) {
	return system.Read(l.ctx, l.svcCtx, "tool:gen:list", func(s *system.Session) (*types.PageResult, error) { return s.GenColumns(req) })
}
