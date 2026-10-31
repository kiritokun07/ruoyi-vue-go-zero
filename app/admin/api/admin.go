// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package main

import (
	"context"
	"flag"
	"fmt"
	"ruoyi-vue-go-zero/app/admin/api/internal/handler/console"
	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"
	"ruoyi-vue-go-zero/app/admin/api/internal/middleware"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest/pathvar"
	"log"
	"net/http"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"

	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"ruoyi-vue-go-zero/app/admin/api/internal/config"
	"ruoyi-vue-go-zero/app/admin/api/internal/handler"
	"ruoyi-vue-go-zero/app/admin/api/internal/handler/avatar"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
)

var configFile = flag.String("f", "etc/admin-api.yaml", "the config file")

func main() {
	flag.Parse()
	sqlx.DisableStmtLog()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	if len(c.Auth.AccessSecret) < 32 || c.MySQLDSN == "" || c.Redis.Host == "" {
		log.Fatal("configure MySQLDSN, Redis.Host and an access secret of at least 32 bytes")
	}
	httpx.SetErrorHandlerCtx(apierror.HTTP)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	schedulerCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	system.StartJobs(schedulerCtx, ctx.DB)
	server.Use(middleware.OperationLog(ctx.DB, c.Auth.AccessSecret))
	handler.RegisterHandlers(server, ctx)
	console.Register(server, ctx)
	server.AddRoute(rest.Route{Method: http.MethodGet, Path: "/profile/avatar/:name", Handler: avatar.File(c.AvatarDir)})

	server.AddRoute(rest.Route{Method: http.MethodGet, Path: "/profile/upload/:name", Handler: func(w http.ResponseWriter, r *http.Request) {
		system.ServeUpload(c.UploadDir, w, pathvar.Vars(r)["name"])
	}})
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
