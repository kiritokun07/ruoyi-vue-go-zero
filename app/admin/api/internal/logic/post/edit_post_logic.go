// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package post

import (
	"context"
	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"

	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EditPostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditPostLogic {
	return &EditPostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditPostLogic) EditPost(req *types.PostRecord) (resp *types.StatusResponse, err error) {
	return system.Write(l.ctx, l.svcCtx, "system:post:edit", req, func(s *system.Session) (*types.StatusResponse, error) { return s.EditPost(req) })
}
