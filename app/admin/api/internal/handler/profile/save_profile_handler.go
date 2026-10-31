// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package profile

import (
	"net/http"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/profile"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SaveProfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ProfileInput
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := profile.NewSaveProfileLogic(r.Context(), svcCtx)
		resp, err := l.SaveProfile(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
