// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package auth

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"ruoyi-vue-go-zero/app/admin/api/internal/middleware"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type GetRoutersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRoutersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoutersLogic {
	return &GetRoutersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRoutersLogic) GetRouters() (resp *types.RoutersResponse, err error) {
	menus, err := l.svcCtx.Access.Menus(l.ctx, middleware.Current(l.ctx).UserID)
	if err != nil {
		return nil, err
	}
	return &types.RoutersResponse{StatusResponse: types.StatusResponse{Code: 200, Msg: "操作成功"}, Data: buildRoutes(menus)}, nil
}
