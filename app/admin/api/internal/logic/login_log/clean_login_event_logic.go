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

type CleanLoginEventLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCleanLoginEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CleanLoginEventLogic {
	return &CleanLoginEventLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CleanLoginEventLogic) CleanLoginEvent() (resp *types.StatusResponse, err error) {
	return system.Write(l.ctx, l.svcCtx, "monitor:logininfor:remove", &types.StatusResponse{}, func(s *system.Session) (*types.StatusResponse, error) {
		return (&catalogService{Session: s}).CleanLoginEvent()
	})
}
