// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package notice

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type MarkNoticeReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMarkNoticeReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkNoticeReadLogic {
	return &MarkNoticeReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarkNoticeReadLogic) MarkNoticeRead(req *types.ReadNoticeRequest) (resp *types.StatusResponse, err error) {
	return system.Write(l.ctx, l.svcCtx, "", req, func(s *system.Session) (*types.StatusResponse, error) { return s.MarkNoticeRead(req) })
}
