// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/kiritokun07/ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"github.com/kiritokun07/ruoyi-vue-go-zero/app/admin/api/internal/config"
	"github.com/kiritokun07/ruoyi-vue-go-zero/app/admin/api/internal/handler"
	"github.com/kiritokun07/ruoyi-vue-go-zero/app/admin/api/internal/svc"
)

var configFile = flag.String("f", "etc/admin-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	if len(c.Auth.AccessSecret) < 32 || c.MySQLDSN == "" || c.Redis.Host == "" {
		log.Fatal("configure MySQLDSN, Redis.Host and an access secret of at least 32 bytes")
	}
	httpx.SetErrorHandlerCtx(apierror.HTTP)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
