// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package dict

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type RefreshDictionaryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshDictionaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshDictionaryLogic {
	return &RefreshDictionaryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshDictionaryLogic) RefreshDictionary() (resp *types.StatusResponse, err error) {
	return system.Read(l.ctx, l.svcCtx, "system:dict:remove", func(s *system.Session) (*types.StatusResponse, error) { return s.RefreshDictionary() })
}
