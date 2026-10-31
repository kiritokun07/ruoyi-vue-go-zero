// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package profile

import (
	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"net/http"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/profile"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadAvatarHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := profile.NewUploadAvatarLogic(r.Context(), svcCtx)
		file, _, err := r.FormFile("avatarfile")
		if r.MultipartForm != nil {
			defer r.MultipartForm.RemoveAll()
		}
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, apierror.Business("请上传头像图片"))
			return
		}
		defer file.Close()
		resp, err := l.UploadAvatar(file)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
