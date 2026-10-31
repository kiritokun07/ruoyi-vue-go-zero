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

type HeaderNoticesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHeaderNoticesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HeaderNoticesLogic {
	return &HeaderNoticesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HeaderNoticesLogic) HeaderNotices() (resp *types.NoticeList, err error) {
	return system.Read(l.ctx, l.svcCtx, "", func(s *system.Session) (*types.NoticeList, error) { return s.HeaderNotices() })
}
