package middleware

import (
	"context"
	"net/http"
	"ruoyi-vue-go-zero/app/admin/api/internal/online"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest/httpx"

	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
)

type Claims struct {
	UserID int64 `json:"userId"`
	jwt.RegisteredClaims
}

type sessionKey struct{}

func Current(ctx context.Context) *Claims { return ctx.Value(sessionKey{}).(*Claims) }
func (c *Claims) RevocationKey() string   { return "ruoyi-go:revoked:" + c.ID }

type SessionMiddleware struct {
	secret []byte
	redis  *redis.Redis
	ttl    int64
	active func(context.Context, int64) (online.Identity, error)
}

func NewSessionMiddleware(secret string, cache *redis.Redis, ttl int64, active func(context.Context, int64) (online.Identity, error)) *SessionMiddleware {
	return &SessionMiddleware{secret: []byte(secret), redis: cache, ttl: ttl, active: active}
}

// ParseBearer returns nil when the request has no valid session token.
func ParseBearer(header string, secret []byte) *Claims {
	fields := strings.Fields(header)
	if len(fields) != 2 || !strings.EqualFold(fields[0], "Bearer") {
		return nil
	}
	var claims Claims
	_, err := jwt.ParseWithClaims(fields[1], &claims, func(*jwt.Token) (any, error) { return secret, nil }, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil || claims.UserID <= 0 || claims.ID == "" {
		return nil
	}
	return &claims
}

func (m *SessionMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims := ParseBearer(r.Header.Get("Authorization"), m.secret)
		if claims == nil {
			httpx.ErrorCtx(ctx, w, apierror.Unauthorized)
			return
		}
		identity, err := m.active(ctx, claims.UserID)
		if err != nil {
			httpx.ErrorCtx(ctx, w, err)
			return
		}
		if identity.UserName == "" {
			httpx.ErrorCtx(ctx, w, apierror.Unauthorized)
			return
		}
		allowed, err := online.Touch(ctx, m.redis, claims.ID, claims.UserID, m.ttl)
		if err != nil {
			httpx.ErrorCtx(ctx, w, err)
			return
		}
		if !allowed {
			httpx.ErrorCtx(ctx, w, apierror.Unauthorized)
			return
		}
		next(w, r.WithContext(context.WithValue(ctx, sessionKey{}, claims)))
	}
}
