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

type EditScheduleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditScheduleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditScheduleLogic {
	return &EditScheduleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditScheduleLogic) EditSchedule(req *types.ScheduleEntry) (resp *types.StatusResponse, err error) {
	return system.Write(l.ctx, l.svcCtx, "monitor:job:edit", req, func(s *system.Session) (*types.StatusResponse, error) {
		return (&catalogService{Session: s}).EditSchedule(req)
	})
}
