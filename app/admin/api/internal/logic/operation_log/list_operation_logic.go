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

type ListOperationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListOperationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOperationLogic {
	return &ListOperationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListOperationLogic) ListOperation(req *types.AdminFilter) (resp *types.PageResult, err error) {
	return system.Read(l.ctx, l.svcCtx, "monitor:operlog:list", func(s *system.Session) (*types.PageResult, error) {
		return (&catalogService{Session: s}).ListOperation(req)
	})
}
