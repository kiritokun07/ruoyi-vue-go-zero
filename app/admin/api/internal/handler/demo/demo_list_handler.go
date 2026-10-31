// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package demo

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/demo"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
)

func DemoListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := demo.NewDemoListLogic(r.Context(), svcCtx)
		resp, err := l.DemoList()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
