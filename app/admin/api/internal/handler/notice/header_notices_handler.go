// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package notice

import (
	"net/http"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/notice"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func HeaderNoticesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := notice.NewHeaderNoticesLogic(r.Context(), svcCtx)
		resp, err := l.HeaderNotices()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
