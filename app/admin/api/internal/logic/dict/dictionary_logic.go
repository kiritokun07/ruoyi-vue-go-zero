// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package dict

import (
	"context"
	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"

	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictionaryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDictionaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictionaryLogic {
	return &DictionaryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DictionaryLogic) Dictionary(req *types.LookupKey) (resp *types.DictResult, err error) {
	return system.Read(l.ctx, l.svcCtx, "", func(s *system.Session) (*types.DictResult, error) { return s.Dictionary(req) })
}
