// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package auth

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"ruoyi-vue-go-zero/app/admin/api/internal/middleware"
	"ruoyi-vue-go-zero/app/admin/api/internal/model"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type GetInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInfoLogic {
	return &GetInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetInfoLogic) GetInfo() (resp *types.InfoResponse, err error) {
	id := middleware.Current(l.ctx).UserID
	user, err := l.svcCtx.Access.Profile(l.ctx, id)
	if errors.Is(err, model.ErrNotFound) {
		return nil, apierror.Unauthorized
	}
	if err != nil {
		return nil, err
	}
	roles, permissions := []string{"admin"}, []string{"*:*:*"}
	if !user.Admin {
		grants, err := l.svcCtx.Access.Permissions(l.ctx, id)
		if err != nil {
			return nil, err
		}
		byRole := make(map[int64][]string)
		var keys, perms []string
		for _, grant := range grants {
			byRole[grant.RoleID] = append(byRole[grant.RoleID], grant.Keys)
			perms = append(perms, grant.Keys)
		}
		for i := range user.Roles {
			role := &user.Roles[i]
			if role.Status == "0" {
				keys = append(keys, role.RoleKey)
				role.Permissions = splitKeys(byRole[role.RoleId])
			}
		}
		roles, permissions = splitKeys(keys), splitKeys(perms)
	}
	initial, err := l.svcCtx.Settings.Value(l.ctx, "sys.account.initPasswordModify")
	if err != nil {
		return nil, err
	}
	days, err := l.svcCtx.Settings.Value(l.ctx, "sys.account.passwordValidateDays")
	if err != nil {
		return nil, err
	}
	limit, _ := strconv.Atoi(days)
	expired := false
	if limit > 0 {
		expired = user.PwdUpdateDate == nil
		if !expired {
			updated, err := time.ParseInLocation(time.DateTime, *user.PwdUpdateDate, time.Local)
			if err != nil {
				return nil, err
			}
			expired = int(time.Since(updated).Abs().Hours()/24) > limit
		}
	}
	return &types.InfoResponse{
		StatusResponse:     types.StatusResponse{Code: 200, Msg: "操作成功"},
		User:               user,
		Roles:              roles,
		Permissions:        permissions,
		IsDefaultModifyPwd: initial == "1" && user.PwdUpdateDate == nil,
		IsPasswordExpired:  expired,
	}, nil
}

func splitKeys(values []string) []string {
	set := make(map[string]bool)
	for _, v := range values {
		for _, key := range strings.Split(v, ",") {
			key = strings.TrimSpace(key)
			if key != "" {
				set[key] = true
			}
		}
	}
	result := make([]string, 0, len(set))
	for key := range set {
		result = append(result, key)
	}
	sort.Strings(result)
	return result
}
