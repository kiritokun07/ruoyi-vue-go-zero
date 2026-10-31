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

type DeleteDictEntryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteDictEntryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDictEntryLogic {
	return &DeleteDictEntryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteDictEntryLogic) DeleteDictEntry(req *types.RecordIds) (resp *types.StatusResponse, err error) {
	return system.Write(l.ctx, l.svcCtx, "system:dict:remove", req, func(s *system.Session) (*types.StatusResponse, error) {
		return (&catalogService{Session: s}).DeleteDictEntry(req)
	})
}
