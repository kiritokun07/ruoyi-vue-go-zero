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

type NoticeDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNoticeDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NoticeDetailLogic {
	return &NoticeDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NoticeDetailLogic) NoticeDetail(req *types.RecordId) (resp *types.NoticeDetail, err error) {
	return system.Read(l.ctx, l.svcCtx, "system:notice:query", func(s *system.Session) (*types.NoticeDetail, error) { return s.NoticeDetail(req) })
}
