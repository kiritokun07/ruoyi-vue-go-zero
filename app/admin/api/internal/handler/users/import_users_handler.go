// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package users

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"ruoyi-vue-go-zero/app/admin/api/internal/logic/users"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
)

func ImportUsersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := users.NewImportUsersLogic(r.Context(), svcCtx)
		file, _, err := r.FormFile("file")
		if r.MultipartForm != nil {
			defer r.MultipartForm.RemoveAll()
		}
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, apierror.Business("请选择导入文件"))
			return
		}
		defer file.Close()
		resp, err := l.ImportUsers(file, r.FormValue("updateSupport") == "true" || r.FormValue("updateSupport") == "1")
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
