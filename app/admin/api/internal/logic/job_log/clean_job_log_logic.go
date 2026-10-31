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

type CleanJobLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCleanJobLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CleanJobLogLogic {
	return &CleanJobLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CleanJobLogLogic) CleanJobLog() (resp *types.StatusResponse, err error) {
	return system.Write(l.ctx, l.svcCtx, "monitor:job:remove", &types.StatusResponse{}, func(s *system.Session) (*types.StatusResponse, error) {
		return (&catalogService{Session: s}).CleanJobLog()
	})
}
