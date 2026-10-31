// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package auth

import (
	"net/http"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/auth"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := auth.NewGetInfoLogic(r.Context(), svcCtx)
		resp, err := l.GetInfo()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
