// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package users

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
)

type UserTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserTemplateLogic {
	return &UserTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserTemplateLogic) UserTemplate() (*system.File, error) {
	return system.Read(l.ctx, l.svcCtx, "", func(s *system.Session) (*system.File, error) { return s.UserTemplate() })
}
