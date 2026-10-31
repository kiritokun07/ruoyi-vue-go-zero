// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package public

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
)

type WelcomeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWelcomeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WelcomeLogic {
	return &WelcomeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WelcomeLogic) Welcome() (*system.File, error) {
	return system.Welcome(), nil
}
