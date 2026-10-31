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

type GenDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenDeleteLogic {
	return &GenDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenDeleteLogic) GenDelete(req *types.RecordIds) (resp *types.StatusResponse, err error) {
	return system.Write(l.ctx, l.svcCtx, "tool:gen:remove", req, func(s *system.Session) (*types.StatusResponse, error) { return s.GenDelete(req) })
}
