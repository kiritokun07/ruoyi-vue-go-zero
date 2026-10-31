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

type ExportPostsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExportPostsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportPostsLogic {
	return &ExportPostsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExportPostsLogic) ExportPosts(req *types.ListFilter) (*system.File, error) {
	return system.Read(l.ctx, l.svcCtx, "system:post:export", func(s *system.Session) (*system.File, error) { return s.ExportPosts(req) })
}
