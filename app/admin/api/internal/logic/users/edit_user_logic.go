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

type EditUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditUserLogic {
	return &EditUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditUserLogic) EditUser(req *types.UserInput) (resp *types.StatusResponse, err error) {
	return system.Write(l.ctx, l.svcCtx, "system:user:edit", req, func(s *system.Session) (*types.StatusResponse, error) { return s.EditUser(req) })
}
