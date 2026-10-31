// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package gen

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type GenDownloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenDownloadLogic {
	return &GenDownloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenDownloadLogic) GenDownload(req *types.LookupKey) (*system.File, error) {
	return system.Read(l.ctx, l.svcCtx, "tool:gen:code", func(s *system.Session) (*system.File, error) { return s.GenDownload(req) })
}
