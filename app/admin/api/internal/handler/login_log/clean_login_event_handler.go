// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package login_log

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/login_log"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
)

func CleanLoginEventHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := login_log.NewCleanLoginEventLogic(r.Context(), svcCtx)
		resp, err := l.CleanLoginEvent()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
