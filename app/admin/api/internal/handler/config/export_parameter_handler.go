// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package config

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/config"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func ExportParameterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminFilter
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := config.NewExportParameterLogic(r.Context(), svcCtx)
		file, err := l.ExportParameter(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			file.Send(w)
		}
	}
}
