// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/kiritokun07/ruoyi-vue-go-zero/app/admin/api/internal/logic"
	"github.com/kiritokun07/ruoyi-vue-go-zero/app/admin/api/internal/svc"
)

func CaptchaImageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewCaptchaImageLogic(r.Context(), svcCtx)
		resp, err := l.CaptchaImage()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
