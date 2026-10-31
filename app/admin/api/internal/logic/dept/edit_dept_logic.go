// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package dept

import (
	"context"
	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"

	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EditDeptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditDeptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditDeptLogic {
	return &EditDeptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditDeptLogic) EditDept(req *types.DeptRecord) (resp *types.StatusResponse, err error) {
	return system.Write(l.ctx, l.svcCtx, "system:dept:edit", req, func(s *system.Session) (*types.StatusResponse, error) { return s.EditDept(req) })
}
