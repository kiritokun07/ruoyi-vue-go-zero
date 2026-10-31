// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package notice

import (
	"net/http"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/notice"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func MarkNoticesReadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ReadNoticesRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := notice.NewMarkNoticesReadLogic(r.Context(), svcCtx)
		resp, err := l.MarkNoticesRead(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
