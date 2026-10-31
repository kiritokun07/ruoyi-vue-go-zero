// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package users

import (
	"context"
	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"

	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDepartmentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserDepartmentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDepartmentsLogic {
	return &UserDepartmentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserDepartmentsLogic) UserDepartments() (resp *types.TreeResult, err error) {
	return system.Read(l.ctx, l.svcCtx, "system:user:list", func(s *system.Session) (*types.TreeResult, error) { return s.UserDepartments() })
}
