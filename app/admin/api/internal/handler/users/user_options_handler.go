// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package users

import (
	"net/http"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/users"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UserOptionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := users.NewUserOptionsLogic(r.Context(), svcCtx)
		resp, err := l.UserOptions()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
