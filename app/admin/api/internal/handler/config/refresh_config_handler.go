// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package config

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/config"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
)

func RefreshConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := config.NewRefreshConfigLogic(r.Context(), svcCtx)
		resp, err := l.RefreshConfig()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
