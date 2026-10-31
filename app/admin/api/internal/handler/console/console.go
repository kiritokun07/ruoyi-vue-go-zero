// Package console serves the two embedded monitoring pages used by the original Vue client.
package console

import (
	"database/sql"
	_ "embed"
	"html/template"
	"io"
	"mime"
	"net/http"
	"path"
	"strings"

	swagger "github.com/swaggo/files/v2"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/pathvar"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
)

//go:embed swagger.json
var specification []byte

const initializer = `window.onload=async function(){
 const spec=await (await fetch('../v3/api-docs')).json();
 spec.basePath=location.pathname.split('/swagger-ui/')[0] || '/';
 spec.schemes=[location.protocol.slice(0,-1)];
 for(const p of Object.values(spec.paths))for(const op of Object.values(p))delete op.schemes;
 window.ui=SwaggerUIBundle({spec,dom_id:'#swagger-ui',validatorUrl:null,requestInterceptor:r=>{
 const cookie=document.cookie.split('; ').find(x=>x.startsWith('Admin-Token='));
 if(cookie && new URL(r.url,location.href).origin===location.origin)r.headers.Authorization='Bearer '+decodeURIComponent(cookie.slice(12));
 return r}});
}`

var dashboard = template.Must(template.New("pool").Parse(`<!doctype html><html lang="zh"><meta charset="utf-8"><title>数据库连接池</title><style>body{font:16px system-ui;margin:3rem;color:#28384c}table{border-collapse:collapse}td,th{padding:10px 24px;border-bottom:1px solid #ddd;text-align:left}</style><h1>MySQL 连接池</h1><p>Go database/sql 实时统计</p><table><tr><th>已打开</th><td>{{.OpenConnections}}</td></tr><tr><th>使用中</th><td>{{.InUse}}</td></tr><tr><th>空闲</th><td>{{.Idle}}</td></tr><tr><th>最大连接数</th><td>{{.MaxOpenConnections}}</td></tr><tr><th>累计等待次数</th><td>{{.WaitCount}}</td></tr><tr><th>累计等待时间</th><td>{{.WaitDuration}}</td></tr></table></html>`))

func Register(server *rest.Server, sc *svc.ServiceContext) {
	server.AddRoute(rest.Route{Method: "GET", Path: "/v3/api-docs", Handler: func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(specification)
	}})
	server.AddRoute(rest.Route{Method: "GET", Path: "/swagger-ui/:asset", Handler: func(w http.ResponseWriter, r *http.Request) {
		asset := pathvar.Vars(r)["asset"]
		if path.Base(asset) != asset {
			http.NotFound(w, r)
			return
		}
		if asset == "swagger-initializer.js" {
			w.Header().Set("Content-Type", "application/javascript")
			w.Write([]byte(initializer))
			return
		}
		f, e := swagger.FS.Open(asset)
		if e != nil {
			http.NotFound(w, r)
			return
		}
		defer f.Close()
		st, e := f.Stat()
		if e != nil || st.IsDir() {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(asset)))
		io.Copy(w, f)
	}})
	protected := sc.Session(func(w http.ResponseWriter, r *http.Request) {
		stats, e := system.Read(r.Context(), sc, "monitor:druid:list", func(_ *system.Session) (*sql.DBStats, error) {
			db, e := sc.DB.RawDB()
			if e != nil {
				return nil, e
			}
			v := db.Stats()
			return &v, nil
		})
		if e != nil {
			httpx.ErrorCtx(r.Context(), w, e)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		dashboard.Execute(w, stats)
	})
	server.AddRoute(rest.Route{Method: "GET", Path: "/druid/login.html", Handler: func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			if c, e := r.Cookie("Admin-Token"); e == nil && !strings.ContainsAny(c.Value, "\r\n") {
				r.Header.Set("Authorization", "Bearer "+c.Value)
			}
		}
		protected(w, r)
	}})
}
