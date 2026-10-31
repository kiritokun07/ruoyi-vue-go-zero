// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package cache

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type ClearCacheNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClearCacheNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearCacheNameLogic {
	return &ClearCacheNameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClearCacheNameLogic) ClearCacheName(req *types.LookupKey) (resp *types.StatusResponse, err error) {
	return system.Read(l.ctx, l.svcCtx, "monitor:cache:list", func(s *system.Session) (*types.StatusResponse, error) { return s.ClearCacheName(req) })
}
