// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package profile

import (
	"context"
	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"

	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveProfileLogic {
	return &SaveProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveProfileLogic) SaveProfile(req *types.ProfileInput) (resp *types.StatusResponse, err error) {
	return system.Write(l.ctx, l.svcCtx, "", req, func(s *system.Session) (*types.StatusResponse, error) { return s.SaveProfile(req) })
}
