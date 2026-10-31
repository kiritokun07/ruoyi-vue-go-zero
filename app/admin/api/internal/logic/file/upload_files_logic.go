// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package file

import (
	"context"
	"mime/multipart"

	"github.com/zeromicro/go-zero/core/logx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type UploadFilesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadFilesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadFilesLogic {
	return &UploadFilesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadFilesLogic) UploadFiles(files []*multipart.FileHeader) (resp *types.UploadResult, err error) {
	return system.Read(l.ctx, l.svcCtx, "", func(s *system.Session) (*types.UploadResult, error) { return s.Upload(files, true, "") })
}
