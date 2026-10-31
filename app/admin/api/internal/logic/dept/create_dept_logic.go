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

type CreateDeptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateDeptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDeptLogic {
	return &CreateDeptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateDeptLogic) CreateDept(req *types.DeptRecord) (resp *types.StatusResponse, err error) {
	return system.Write(l.ctx, l.svcCtx, "system:dept:add", req, func(s *system.Session) (*types.StatusResponse, error) { return s.CreateDept(req) })
}
