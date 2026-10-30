package logic

import (
	"context"
	"testing"

	"github.com/zeromicro/go-zero/core/stores/redis/redistest"
	"golang.org/x/crypto/bcrypt"

	"github.com/kiritokun07/ruoyi-vue-go-zero/app/admin/api/internal/model"
	"github.com/kiritokun07/ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"github.com/kiritokun07/ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type captchaSettings struct{ model.SysConfigModel }

func (captchaSettings) Value(context.Context, string) (string, error) { return "true", nil }

func TestWrongCaptchaConsumesChallenge(t *testing.T) {
	ctx := context.Background()
	services := &svc.ServiceContext{Redis: redistest.CreateRedis(t), Settings: captchaSettings{}}
	l := NewLoginLogic(ctx, services)
	if err := services.Redis.SetexCtx(ctx, captchaKey+"sample", "abc", 120); err != nil {
		t.Fatal(err)
	}
	req := &types.LoginRequest{Uuid: "sample", Code: "wrong"}
	if _, err := l.authenticate(req, "127.0.0.1"); err == nil || err.Error() != "验证码错误" {
		t.Fatalf("first attempt: %v", err)
	}
	req.Code = "abc"
	if _, err := l.authenticate(req, "127.0.0.1"); err == nil || err.Error() != "验证码已失效" {
		t.Fatalf("replay: %v", err)
	}
}

func TestPasswordLock(t *testing.T) {
	services := &svc.ServiceContext{Redis: redistest.CreateRedis(t)}
	services.Config.Auth.MaxRetryCount = 2
	services.Config.Auth.LockMinutes = 10
	l := NewLoginLogic(context.Background(), services)
	hash, err := bcrypt.GenerateFromPassword([]byte("valid-password"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	user := &model.LoginUser{UserID: 1, Password: string(hash)}
	for _, tc := range []struct{ password, want string }{
		{"incorrect", badCredentials.Error()},
		{"valid-password", ""},
		{"incorrect", badCredentials.Error()},
		{"incorrect", badCredentials.Error()},
		{"valid-password", "密码输入错误2次，帐户锁定10分钟"},
	} {
		err := l.checkPassword(user, tc.password)
		if err == nil && tc.want != "" || err != nil && err.Error() != tc.want {
			t.Fatalf("got %v, want %q", err, tc.want)
		}
	}
}
