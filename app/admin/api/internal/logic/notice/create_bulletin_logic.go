// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package notice

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type CreateBulletinLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateBulletinLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBulletinLogic {
	return &CreateBulletinLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateBulletinLogic) CreateBulletin(req *types.BulletinEntry) (resp *types.StatusResponse, err error) {
	return system.Write(l.ctx, l.svcCtx, "system:notice:add", req, func(s *system.Session) (*types.StatusResponse, error) {
		return (&catalogService{Session: s}).CreateBulletin(req)
	})
}
