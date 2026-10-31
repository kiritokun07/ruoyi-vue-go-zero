package system

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"ruoyi-vue-go-zero/app/admin/api/internal/middleware"
	"ruoyi-vue-go-zero/app/admin/api/internal/online"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

// A lock held on the administrator must not serialize unrelated writes.
// Both transactions roll back their changes.
func TestWriteDoesNotLockAdministrator(t *testing.T) {
	dsn := os.Getenv("RUOYI_SYSTEM_TEST_DSN")
	if dsn == "" {
		t.Skip("set RUOYI_SYSTEM_TEST_DSN")
	}
	db := sqlx.NewMysql(dsn)
	cache := miniredis.RunT(t)
	secret := "transaction-regression-test"
	auth := middleware.NewSessionMiddleware(secret, redis.New(cache.Addr()), 1800, func(context.Context, int64) (online.Identity, error) {
		return online.Identity{UserName: "admin"}, nil
	})
	id := uuid.NewString()
	if err := online.Create(context.Background(), redis.New(cache.Addr()), online.Session{UserID: 1, TokenId: id}, 1800); err != nil {
		t.Fatal(err)
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"userId": 1, "jti": id, "exp": time.Now().Add(time.Minute).Unix()}).SignedString([]byte(secret))
	if err != nil {
		t.Fatal(err)
	}
	rollback := errors.New("rollback fixtures")
	key := "lock_" + uuid.NewString()
	err = db.TransactCtx(context.Background(), func(ctx context.Context, tx sqlx.Session) error {
		var id int64
		if err := tx.QueryRowCtx(ctx, &id, "SELECT user_id FROM sys_user WHERE user_id=1 FOR UPDATE"); err != nil {
			return err
		}
		requestCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		req := httptest.NewRequest(http.MethodPost, "/test", nil).WithContext(requestCtx)
		req.Header.Set("Authorization", "Bearer "+token)
		reached := false
		var writeErr error
		auth.Handle(func(w http.ResponseWriter, r *http.Request) {
			reached = true
			_, writeErr = withSession(r.Context(), &svc.ServiceContext{DB: db}, "", true, func(s *Session) (*types.StatusResponse, error) {
				if err := s.exec("INSERT INTO sys_config(config_name,config_key,config_value,config_type,create_time) VALUES(?,?,?,'N',NOW())", key, key, "test"); err != nil {
					return nil, err
				}
				return nil, rollback
			})
		})(httptest.NewRecorder(), req)
		if !reached {
			return errors.New("middleware rejected test request")
		}
		if !errors.Is(writeErr, rollback) {
			return writeErr
		}
		return rollback
	})
	if !errors.Is(err, rollback) {
		t.Fatalf("unrelated write blocked or failed: %v", err)
	}
	var count int
	if err = db.QueryRowCtx(context.Background(), &count, "SELECT COUNT(*) FROM sys_config WHERE config_key=?", key); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatal("write transaction did not roll back")
	}
}
