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

type ExcludeDepartmentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExcludeDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExcludeDepartmentLogic {
	return &ExcludeDepartmentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExcludeDepartmentLogic) ExcludeDepartment(req *types.RecordId) (resp *types.DeptList, err error) {
	return system.Read(l.ctx, l.svcCtx, "system:dept:list", func(s *system.Session) (*types.DeptList, error) { return s.ExcludeDepartment(req) })
}
