package system

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"ruoyi-vue-go-zero/app/admin/api/internal/model"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func TestRedisUsernameWrites(t *testing.T) {
	dsn := os.Getenv("RUOYI_SYSTEM_TEST_DSN")
	if dsn == "" {
		t.Skip("set RUOYI_SYSTEM_TEST_DSN")
	}
	db := sqlx.NewMysql(dsn)
	ctx := context.Background()
	server := miniredis.RunT(t)
	sc := &svc.ServiceContext{DB: db, Redis: redis.New(server.Addr())}
	name := "lock_" + uuid.NewString()[:8]
	defer db.ExecCtx(ctx, "DELETE FROM sys_user WHERE user_name=?", name)
	var count int
	if err := db.QueryRowCtx(ctx, &count, "SELECT COUNT(*) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='sys_user' AND INDEX_NAME='uk_sys_user_active_name'"); err != nil || count != 0 {
		t.Fatal("test requires original schema without added unique constraint", err)
	}
	insert := func(ctx context.Context) (*types.StatusResponse, error) {
		err := db.TransactCtx(ctx, func(ctx context.Context, tx sqlx.Session) error {
			var n int
			if err := tx.QueryRowCtx(ctx, &n, "SELECT COUNT(*) FROM sys_user WHERE user_name=? AND del_flag='0'", name); err != nil {
				return err
			}
			if n > 0 {
				return apierror.Business("用户账号已存在")
			}
			_, err := tx.ExecCtx(ctx, "INSERT INTO sys_user(user_name,nick_name) VALUES(?,?)", name, name)
			return err
		})
		return saved(err)
	}
	var succeeded atomic.Int32
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_, err := withUserWrite(ctx, sc, name, insert)
			var business apierror.Business
			if err == nil {
				succeeded.Add(1)
			} else if !errors.As(err, &business) {
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}
	close(start)
	wg.Wait()
	if succeeded.Load() != 1 {
		t.Fatalf("created %d accounts", succeeded.Load())
	}
	if err := db.QueryRowCtx(ctx, &count, "SELECT COUNT(*) FROM sys_user WHERE user_name=?", name); err != nil || count != 1 {
		t.Fatalf("count=%d err=%v", count, err)
	}
	// A held same-account lock blocks other instances, but not another account.
	key, err := model.NewAdminModel(db).UserLockKey(ctx, name)
	if err != nil {
		t.Fatal(err)
	}
	owner := redis.NewRedisLock(sc.Redis, key)
	owner.SetExpire(60)
	if ok, err := owner.AcquireCtx(ctx); err != nil || !ok {
		t.Fatal(ok, err)
	}
	if _, err := withUserWrite(ctx, sc, name, insert); err == nil {
		t.Fatal("concurrent owner accepted")
	}
	if _, err := withUserWrite(ctx, sc, name+"x", func(context.Context) (*types.StatusResponse, error) { return nil, nil }); err != nil {
		t.Fatal("different account blocked", err)
	}
	if _, err := owner.ReleaseCtx(ctx); err != nil {
		t.Fatal(err)
	}
	// A cancelled request still releases its own lock.
	cancelled, cancel := context.WithCancel(ctx)
	_, err = withUserWrite(cancelled, sc, name, func(context.Context) (*types.StatusResponse, error) { cancel(); return nil, context.Canceled })
	if !errors.Is(err, context.Canceled) || server.Exists(key) {
		t.Fatal("lock leaked after cancellation", err)
	}
	if _, err := db.ExecCtx(ctx, "UPDATE sys_user SET del_flag='2' WHERE user_name=?", name); err != nil {
		t.Fatal(err)
	}
	if _, err := withUserWrite(ctx, sc, name, insert); err != nil {
		t.Fatal("soft-deleted name cannot be reused", err)
	}
	// Case variants use the same lock whenever the database treats them as equal.
	var equal int
	if err := db.QueryRowCtx(ctx, &equal, "SELECT COUNT(*) FROM sys_user WHERE user_name=?", strings.ToUpper(name)); err != nil {
		t.Fatal(err)
	}
	if equal > 0 {
		other, err := model.NewAdminModel(db).UserLockKey(ctx, strings.ToUpper(name))
		if err != nil || key != other {
			t.Fatal("case-equivalent names have different locks", err)
		}
	}
	// Redis outage fails closed: do not execute the database callback.
	server.Close()
	if _, err := withUserWrite(ctx, sc, name, func(context.Context) (*types.StatusResponse, error) {
		t.Error("write proceeded without Redis")
		return nil, nil
	}); err == nil {
		t.Fatal("Redis failure ignored")
	}
}
