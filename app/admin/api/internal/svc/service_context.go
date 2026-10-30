package svc

import (
	"context"
	"strings"

	"github.com/go-playground/validator/v10"
	bc "github.com/mojocn/base64Captcha"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"github.com/kiritokun07/ruoyi-vue-go-zero/app/admin/api/internal/captcha"
	"github.com/kiritokun07/ruoyi-vue-go-zero/app/admin/api/internal/config"
	"github.com/kiritokun07/ruoyi-vue-go-zero/app/admin/api/internal/model"
)

type ServiceContext struct {
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
	return &ServiceContext{
		Config:    c,
		Users:     model.NewSysUserModel(conn),
		Settings:  model.NewSysConfigModel(conn),
		LoginLog:  model.NewSysLogininforModel(conn),
		Redis:     redis.MustNewRedis(c.Redis),
		Validator: validator.New(),
		Captcha:   captcha.NewDriver(c.Auth.CaptchaType),
	}
}

func (s *ServiceContext) CaptchaEnabled(ctx context.Context) (bool, error) {
	value, err := s.Settings.Value(ctx, "sys.account.captchaEnabled")
	return !strings.EqualFold(value, "false"), err
}
