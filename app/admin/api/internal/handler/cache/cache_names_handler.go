// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package cache

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/cache"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
)

func CacheNamesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := cache.NewCacheNamesLogic(r.Context(), svcCtx)
		resp, err := l.CacheNames()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
