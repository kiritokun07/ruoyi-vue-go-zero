package avatar

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

func File(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := pathvar.Vars(r)["name"]
		id, err := uuid.Parse(strings.TrimSuffix(name, ".jpg"))
		if err != nil || name != id.String()+".jpg" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "image/jpeg")
		http.ServeFile(w, r, filepath.Join(dir, name))
	}
}
