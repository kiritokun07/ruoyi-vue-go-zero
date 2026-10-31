package svc

import (
	"context"
	"strings"

	"github.com/go-playground/validator/v10"
	bc "github.com/mojocn/base64Captcha"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"

	"ruoyi-vue-go-zero/app/admin/api/internal/captcha"
	"ruoyi-vue-go-zero/app/admin/api/internal/config"
	"ruoyi-vue-go-zero/app/admin/api/internal/middleware"
	"ruoyi-vue-go-zero/app/admin/api/internal/model"
)

type ServiceContext struct {
	Session   rest.Middleware
	Admin     *model.AdminModel
	Access    *model.AccessModel
	DB        sqlx.SqlConn
	Config    config.Config
	Users     model.SysUserModel
	Settings  model.SysConfigModel
	LoginLog  model.SysLogininforModel
	Redis     *redis.Redis
	Validator *validator.Validate
	Captcha   bc.Driver
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MySQLDSN)
	cache := redis.MustNewRedis(c.Redis)
	access := model.NewAccessModel(conn)
	return &ServiceContext{
		Session:   middleware.NewSessionMiddleware(c.Auth.AccessSecret, cache, c.Auth.AccessExpire, access.ActiveIdentity).Handle,
		Access:    access,
		Admin:     model.NewAdminModel(conn),
		DB:        conn,
		Config:    c,
		Users:     model.NewSysUserModel(conn),
		Settings:  model.NewSysConfigModel(conn),
		LoginLog:  model.NewSysLogininforModel(conn),
		Redis:     cache,
		Validator: validator.New(),
		Captcha:   captcha.NewDriver(c.Auth.CaptchaType),
	}
}

func (s *ServiceContext) CaptchaEnabled(ctx context.Context) (bool, error) {
	value, err := s.Settings.Value(ctx, "sys.account.captchaEnabled")
	return !strings.EqualFold(value, "false"), err
}
