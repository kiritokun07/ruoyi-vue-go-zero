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

type CacheKeysLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCacheKeysLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CacheKeysLogic {
	return &CacheKeysLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CacheKeysLogic) CacheKeys(req *types.LookupKey) (resp *types.Result, err error) {
	return system.Read(l.ctx, l.svcCtx, "monitor:cache:list", func(s *system.Session) (*types.Result, error) { return s.CacheKeys(req) })
}
