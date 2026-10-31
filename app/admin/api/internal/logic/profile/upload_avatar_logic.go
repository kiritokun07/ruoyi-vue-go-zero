// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package profile

import (
	"context"
	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"
	"io"

	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadAvatarLogic {
	return &UploadAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadAvatarLogic) UploadAvatar(file io.Reader) (resp *types.AvatarResult, err error) {
	return system.SaveAvatar(l.ctx, l.svcCtx, file)
}
