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

type GetDictTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDictTypeLogic {
	return &GetDictTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDictTypeLogic) GetDictType(req *types.RecordId) (resp *types.Result, err error) {
	return system.Read(l.ctx, l.svcCtx, "system:dict:query", func(s *system.Session) (*types.Result, error) { return (&catalogService{Session: s}).GetDictType(req) })
}
