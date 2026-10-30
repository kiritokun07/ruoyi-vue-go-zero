package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	MySQLDSN string
	Redis    redis.RedisConf
	Auth     struct {
		AccessSecret  string
		AccessExpire  int64  `json:",default=1800,range=[1:604800]"`
		CaptchaType   string `json:",default=math,options=math|char"`
		MaxRetryCount int    `json:",default=5,range=[1:100]"`
		LockMinutes   int    `json:",default=10,range=[1:1440]"`
	}
}
