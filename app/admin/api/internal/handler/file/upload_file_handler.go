// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package file

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"ruoyi-vue-go-zero/app/admin/api/internal/logic/file"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
)

func UploadFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := file.NewUploadFileLogic(r.Context(), svcCtx)
		err := r.ParseMultipartForm(4 << 20)
		if r.MultipartForm != nil {
			defer r.MultipartForm.RemoveAll()
		}
		if err != nil || r.MultipartForm == nil {
			httpx.ErrorCtx(r.Context(), w, apierror.Business("上传格式无效"))
			return
		}
		resp, err := l.UploadFile(r.MultipartForm.File["file"])
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
