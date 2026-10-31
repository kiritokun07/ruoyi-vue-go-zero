package system

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"ruoyi-vue-go-zero/app/admin/api/internal/model"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

// Acquire before beginning a transaction and release only after commit/rollback.
// Different account names never share a lock. A busy account fails fast.
func withUserWrite[T any](ctx context.Context, sc *svc.ServiceContext, name string, write func(context.Context) (*T, error)) (*T, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	key, err := model.NewAdminModel(sc.DB).UserLockKey(ctx, name)
	if err != nil {
		return nil, err
	}
	// SET NX avoids concurrent EVALSHA cache misses tripping the Redis breaker.
	// The lease exceeds the bounded transaction lifetime, including cleanup time.
	owner := uuid.NewString()
	acquired, err := sc.Redis.SetnxExCtx(ctx, key, owner, 60)
	if err != nil {
		return nil, err
	}
	if !acquired {
		return nil, apierror.Business("该账号正在处理，请稍后重试")
	}
	defer func() {
		releaseCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if _, err := sc.Redis.EvalCtx(releaseCtx, `if redis.call("GET",KEYS[1])==ARGV[1] then return redis.call("DEL",KEYS[1]) end;return 0`, []string{key}, owner); err != nil {
			logx.WithContext(ctx).Error(err)
		}
	}()
	return write(ctx)
}

func CreateUser(ctx context.Context, sc *svc.ServiceContext, req *types.UserInput) (*types.StatusResponse, error) {
	return withUserWrite(ctx, sc, req.UserName, func(ctx context.Context) (*types.StatusResponse, error) {
		return Write(ctx, sc, "system:user:add", req, func(s *Session) (*types.StatusResponse, error) { return s.CreateUser(req) })
	})
}
