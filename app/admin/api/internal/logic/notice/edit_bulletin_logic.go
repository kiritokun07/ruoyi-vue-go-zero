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

type EditBulletinLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditBulletinLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditBulletinLogic {
	return &EditBulletinLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditBulletinLogic) EditBulletin(req *types.BulletinEntry) (resp *types.StatusResponse, err error) {
	return system.Write(l.ctx, l.svcCtx, "system:notice:edit", req, func(s *system.Session) (*types.StatusResponse, error) {
		return (&catalogService{Session: s}).EditBulletin(req)
	})
}
