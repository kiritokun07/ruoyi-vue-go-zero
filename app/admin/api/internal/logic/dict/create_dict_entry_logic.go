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

type CreateDictEntryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateDictEntryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDictEntryLogic {
	return &CreateDictEntryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateDictEntryLogic) CreateDictEntry(req *types.DictEntryEntry) (resp *types.StatusResponse, err error) {
	return system.Write(l.ctx, l.svcCtx, "system:dict:add", req, func(s *system.Session) (*types.StatusResponse, error) {
		return (&catalogService{Session: s}).CreateDictEntry(req)
	})
}
