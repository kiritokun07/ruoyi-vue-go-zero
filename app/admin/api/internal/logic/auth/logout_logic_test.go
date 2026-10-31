package auth

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/stores/redis/redistest"

	"ruoyi-vue-go-zero/app/admin/api/internal/middleware"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
)

func TestLogoutOnlyRevokesValidToken(t *testing.T) {
	ctx := context.Background()
	s := &svc.ServiceContext{Redis: redistest.CreateRedis(t)}
	s.Config.Auth.AccessSecret = "test-key"
	l := NewLogoutLogic(ctx, s)
	for _, tc := range []struct {
		name, key string
		ttl       time.Duration
		revoked   bool
	}{
		{"expired", "test-key", -time.Minute, false},
		{"forged", "other-key", time.Hour, false},
		{"valid", "test-key", time.Hour, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := middleware.Claims{UserID: 1, RegisteredClaims: jwt.RegisteredClaims{ID: tc.name, ExpiresAt: jwt.NewNumericDate(time.Now().Add(tc.ttl))}}
			raw, err := jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString([]byte(tc.key))
			if err != nil {
				t.Fatal(err)
			}
			for _, h := range []string{"", "Bearer broken", "Bearer " + raw, "Bearer " + raw} {
				resp, err := l.Logout(h)
				if err != nil || resp.Code != 200 || resp.Msg != "退出成功" {
					t.Fatalf("logout: %v, %v", resp, err)
				}
			}
			got, err := s.Redis.ExistsCtx(ctx, c.RevocationKey())
			if err != nil || got != tc.revoked {
				t.Fatalf("revoked=%t err=%v", got, err)
			}
		})
	}
}
