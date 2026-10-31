// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operation_log

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type DeleteOperationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteOperationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteOperationLogic {
	return &DeleteOperationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteOperationLogic) DeleteOperation(req *types.RecordIds) (resp *types.StatusResponse, err error) {
	return system.Write(l.ctx, l.svcCtx, "monitor:operlog:remove", req, func(s *system.Session) (*types.StatusResponse, error) {
		return (&catalogService{Session: s}).DeleteOperation(req)
	})
}
