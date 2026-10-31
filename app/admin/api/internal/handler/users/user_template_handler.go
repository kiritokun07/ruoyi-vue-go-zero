// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package users

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/users"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
)

func UserTemplateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := users.NewUserTemplateLogic(r.Context(), svcCtx)
		file, err := l.UserTemplate()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			file.Send(w)
		}
	}
}
