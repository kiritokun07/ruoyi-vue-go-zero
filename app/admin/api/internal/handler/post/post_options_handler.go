// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package post

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/post"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
)

func PostOptionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := post.NewPostOptionsLogic(r.Context(), svcCtx)
		resp, err := l.PostOptions()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
