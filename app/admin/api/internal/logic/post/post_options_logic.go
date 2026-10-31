// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package post

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type PostOptionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPostOptionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostOptionsLogic {
	return &PostOptionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PostOptionsLogic) PostOptions() (resp *types.Result, err error) {
	return system.Read(l.ctx, l.svcCtx, "", func(s *system.Session) (*types.Result, error) { return s.PostOptions() })
}
