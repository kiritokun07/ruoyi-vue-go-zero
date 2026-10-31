// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package auth

import (
	"net"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"ruoyi-vue-go-zero/app/admin/api/internal/logic/auth"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, apierror.Business("请求参数错误"))
			return
		}

		l := auth.NewLoginLogic(r.Context(), svcCtx)
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		resp, err := l.Login(&req, ip, r.UserAgent())
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
