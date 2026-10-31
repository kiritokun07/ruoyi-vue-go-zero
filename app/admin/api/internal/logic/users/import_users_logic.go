// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package users

import (
	"context"
	"io"

	"github.com/zeromicro/go-zero/core/logx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type ImportUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImportUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImportUsersLogic {
	return &ImportUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImportUsersLogic) ImportUsers(file io.Reader, update bool) (resp *types.Result, err error) {
	return system.ImportUserFile(l.ctx, l.svcCtx, file, update)
}
