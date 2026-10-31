// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package cache

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/cache"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func ClearCacheNameHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LookupKey
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := cache.NewClearCacheNameLogic(r.Context(), svcCtx)
		resp, err := l.ClearCacheName(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
