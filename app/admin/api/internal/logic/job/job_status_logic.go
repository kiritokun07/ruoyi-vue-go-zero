// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package job

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type JobStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJobStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JobStatusLogic {
	return &JobStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JobStatusLogic) JobStatus(req *types.ScheduleEntry) (resp *types.StatusResponse, err error) {
	return system.Write(l.ctx, l.svcCtx, "monitor:job:changeStatus", req, func(s *system.Session) (*types.StatusResponse, error) { return s.JobStatus(req) })
}
