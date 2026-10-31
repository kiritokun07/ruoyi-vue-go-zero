// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package serverinfo

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/serverinfo"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
)

func ServerInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := serverinfo.NewServerInfoLogic(r.Context(), svcCtx)
		resp, err := l.ServerInfo()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
