// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package auth

import (
	"context"
	"database/sql"
	"ruoyi-vue-go-zero/app/admin/api/internal/model"
	"ruoyi-vue-go-zero/app/admin/api/internal/online"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"ruoyi-vue-go-zero/app/admin/api/internal/middleware"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(authorization string) (*types.StatusResponse, error) {
	success := &types.StatusResponse{Code: 200, Msg: "退出成功"}
	claims := middleware.ParseBearer(authorization, []byte(l.svcCtx.Config.Auth.AccessSecret))
	if claims == nil {
		return success, nil
	}
	ttl := claims.ExpiresAt.Unix() - time.Now().Unix()
	if ttl > 0 {
		v, err := online.Revoke(l.ctx, l.svcCtx.Redis, claims.ID, claims.ExpiresAt.Unix())
		if err != nil {
			return nil, err
		}
		if v != nil {
			_, err = l.svcCtx.LoginLog.Insert(l.ctx, &model.SysLogininfor{UserName: v.UserName, Ipaddr: v.Ipaddr, LoginLocation: v.LoginLocation, Browser: v.Browser, Os: v.OS, Status: "0", Msg: "退出成功", LoginTime: sql.NullTime{Time: time.Now(), Valid: true}})
			if err != nil {
				l.Error(err)
			}
		}
	}
	return success, nil
}
