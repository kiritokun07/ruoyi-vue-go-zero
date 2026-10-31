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

type EditDictEntryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditDictEntryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditDictEntryLogic {
	return &EditDictEntryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditDictEntryLogic) EditDictEntry(req *types.DictEntryEntry) (resp *types.StatusResponse, err error) {
	return system.Write(l.ctx, l.svcCtx, "system:dict:edit", req, func(s *system.Session) (*types.StatusResponse, error) {
		return (&catalogService{Session: s}).EditDictEntry(req)
	})
}
