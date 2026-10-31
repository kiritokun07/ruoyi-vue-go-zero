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

type GetDeptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDeptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeptLogic {
	return &GetDeptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDeptLogic) GetDept(req *types.RecordId) (resp *types.DeptResult, err error) {
	return system.Read(l.ctx, l.svcCtx, "system:dept:query", func(s *system.Session) (*types.DeptResult, error) { return s.GetDept(req) })
}
