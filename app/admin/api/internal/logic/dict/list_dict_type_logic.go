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

type ListDictTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDictTypeLogic {
	return &ListDictTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListDictTypeLogic) ListDictType(req *types.AdminFilter) (resp *types.PageResult, err error) {
	return system.Read(l.ctx, l.svcCtx, "system:dict:list", func(s *system.Session) (*types.PageResult, error) {
		return (&catalogService{Session: s}).ListDictType(req)
	})
}
