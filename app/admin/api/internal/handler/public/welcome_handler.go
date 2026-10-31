// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package public

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/public"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
)

func WelcomeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewWelcomeLogic(r.Context(), svcCtx)
		file, err := l.Welcome()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			file.Send(w)
		}
	}
}
