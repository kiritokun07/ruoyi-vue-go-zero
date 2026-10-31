// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package login_log

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type ListLoginEventLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLoginEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLoginEventLogic {
	return &ListLoginEventLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLoginEventLogic) ListLoginEvent(req *types.AdminFilter) (resp *types.PageResult, err error) {
	return system.Read(l.ctx, l.svcCtx, "monitor:logininfor:list", func(s *system.Session) (*types.PageResult, error) {
		return (&catalogService{Session: s}).ListLoginEvent(req)
	})
}
