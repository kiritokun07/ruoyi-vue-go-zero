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

type UserGrantsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserGrantsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserGrantsLogic {
	return &UserGrantsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserGrantsLogic) UserGrants(req *types.RecordId) (resp *types.UserGrants, err error) {
	return system.Read(l.ctx, l.svcCtx, "system:user:query", func(s *system.Session) (*types.UserGrants, error) { return s.UserGrants(req) })
}
