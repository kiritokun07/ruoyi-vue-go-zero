// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package role

import (
	"context"
	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"

	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleDepartmentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleDepartmentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleDepartmentsLogic {
	return &RoleDepartmentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleDepartmentsLogic) RoleDepartments(req *types.RecordId) (resp *types.DeptTreeResult, err error) {
	return system.Read(l.ctx, l.svcCtx, "system:role:query", func(s *system.Session) (*types.DeptTreeResult, error) { return s.RoleDepartments(req) })
}
