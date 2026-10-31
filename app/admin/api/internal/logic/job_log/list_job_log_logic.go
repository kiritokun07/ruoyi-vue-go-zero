// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package job_log

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type ListJobLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListJobLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListJobLogLogic {
	return &ListJobLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListJobLogLogic) ListJobLog(req *types.AdminFilter) (resp *types.PageResult, err error) {
	return system.Read(l.ctx, l.svcCtx, "monitor:job:list", func(s *system.Session) (*types.PageResult, error) {
		return (&catalogService{Session: s}).ListJobLog(req)
	})
}
