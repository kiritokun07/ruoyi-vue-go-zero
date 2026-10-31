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

type ListPostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPostLogic {
	return &ListPostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPostLogic) ListPost(req *types.ListFilter) (resp *types.PostList, err error) {
	return system.Read(l.ctx, l.svcCtx, "system:post:list", func(s *system.Session) (*types.PostList, error) { return s.ListPost(req) })
}
