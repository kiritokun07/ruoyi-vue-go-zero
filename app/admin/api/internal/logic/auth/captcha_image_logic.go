// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package auth

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"ruoyi-vue-go-zero/app/admin/api/internal/captcha"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type CaptchaImageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCaptchaImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaptchaImageLogic {
	return &CaptchaImageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CaptchaImageLogic) CaptchaImage() (resp *types.CaptchaResponse, err error) {
	enabled, err := l.svcCtx.CaptchaEnabled(l.ctx)
	if err != nil {
		return nil, err
	}
	resp = &types.CaptchaResponse{Code: 200, Msg: "操作成功", CaptchaEnabled: enabled}
	if !enabled {
		return resp, nil
	}
	id, image, answer, err := captcha.Render(l.svcCtx.Captcha)
	if err != nil {
		return nil, err
	}
	if err = l.svcCtx.Redis.SetexCtx(l.ctx, captchaKey+id, answer, 120); err != nil {
		return nil, err
	}
	resp.Uuid = id
	resp.Img = image
	return resp, nil
}
