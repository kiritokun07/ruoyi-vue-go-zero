// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operation_log

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/operation_log"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
)

func CleanOperationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := operation_log.NewCleanOperationLogic(r.Context(), svcCtx)
		resp, err := l.CleanOperation()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
