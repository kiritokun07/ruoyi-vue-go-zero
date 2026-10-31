package middleware

import (
	"context"
	"encoding/json"
	"ruoyi-vue-go-zero/app/admin/api/internal/online"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/stores/redis/redistest"
	"github.com/zeromicro/go-zero/rest/httpx"

	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
)

func TestSession(t *testing.T) {
	httpx.SetErrorHandlerCtx(apierror.HTTP)
	cache := redistest.CreateRedis(t)
	active := true
	m := NewSessionMiddleware("test-key", cache, func(context.Context, int64) (online.Identity, error) {
		if active {
			return online.Identity{UserName: "member"}, nil
		}
		return online.Identity{}, nil
	})
	token := func(id string, expiry time.Time) string {
		c := Claims{UserID: 7, RegisteredClaims: jwt.RegisteredClaims{ID: id, ExpiresAt: jwt.NewNumericDate(expiry)}}
		raw, err := jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(m.secret)
		if err != nil {
			t.Fatal(err)
		}
		return raw
	}
	request := func(raw string) int {
		r := httptest.NewRequest("GET", "/getInfo", nil)
		if raw != "" {
			r.Header.Set("Authorization", "Bearer "+raw)
		}
		w := httptest.NewRecorder()
		m.Handle(func(w http.ResponseWriter, r *http.Request) {
			if Current(r.Context()).UserID != 7 {
				t.Fatal("wrong user")
			}
			w.WriteHeader(http.StatusNoContent)
		})(w, r)
		if w.Code == http.StatusNoContent {
			return 200
		}
		var result struct{ Code int }
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		return result.Code
	}
	now := time.Now()
	first := token("one", now.Add(time.Hour))
	for _, raw := range []string{"", first + "x", token("expired", now.Add(-time.Minute)), token("", now.Add(time.Hour))} {
		if request(raw) != 401 {
			t.Fatal("invalid token accepted")
		}
	}
	if request(first) != 200 {
		t.Fatal("valid token rejected")
	}
	if err := cache.SetexCtx(context.Background(), (&Claims{RegisteredClaims: jwt.RegisteredClaims{ID: "one"}}).RevocationKey(), "1", 3600); err != nil {
		t.Fatal(err)
	}
	other := token("two", now.Add(time.Hour))
	if request(first) != 401 || request(other) != 200 {
		t.Fatal("revocation affected wrong session")
	}
	active = false
	if request(other) != 401 {
		t.Fatal("inactive user accepted")
	}
}
