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

type CreateDictTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDictTypeLogic {
	return &CreateDictTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateDictTypeLogic) CreateDictType(req *types.DictTypeEntry) (resp *types.StatusResponse, err error) {
	return system.Write(l.ctx, l.svcCtx, "system:dict:add", req, func(s *system.Session) (*types.StatusResponse, error) {
		return (&catalogService{Session: s}).CreateDictType(req)
	})
}
