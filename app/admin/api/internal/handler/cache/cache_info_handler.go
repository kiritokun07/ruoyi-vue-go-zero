// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package cache

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/cache"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
)

func CacheInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := cache.NewCacheInfoLogic(r.Context(), svcCtx)
		resp, err := l.CacheInfo()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
