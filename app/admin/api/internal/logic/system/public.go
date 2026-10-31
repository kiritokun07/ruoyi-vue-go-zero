package system

import (
	"context"
	"sort"
	"strings"
	"sync"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func Register(ctx context.Context, sc *svc.ServiceContext, r *types.LoginRequest) (*types.StatusResponse, error) {
	v, e := sc.Settings.Value(ctx, "sys.account.registerUser")
	if e != nil {
		return nil, e
	}
	if !strings.EqualFold(v, "true") {
		return nil, apierror.Business("当前系统没有开启注册功能！")
	}
	enabled, e := sc.CaptchaEnabled(ctx)
	if e != nil {
		return nil, e
	}
	if enabled {
		a, e := sc.Redis.GetDelCtx(ctx, "ruoyi-go:captcha:"+r.Uuid)
		if e != nil {
			return nil, e
		}
		if a == "" {
			return nil, apierror.Business("验证码已失效")
		}
		if !strings.EqualFold(a, r.Code) {
			return nil, apierror.Business("验证码错误")
		}
	}
	if sc.Validator.StructCtx(ctx, r) != nil {
		return nil, apierror.Business("用户名或密码长度无效")
	}
	hash, e := hashPassword(r.Password)
	if e != nil {
		return nil, e
	}
	return withUserWrite(ctx, sc, r.Username, func(ctx context.Context) (*types.StatusResponse, error) {
		e := sc.DB.TransactCtx(ctx, func(ctx context.Context, tx sqlx.Session) error {
			s := Session{ctx: ctx, db: tx}
			if e := s.rejectIf("SELECT COUNT(*) FROM sys_user WHERE user_name=? AND del_flag='0'", "注册账号已存在", r.Username); e != nil {
				return e
			}
			return s.exec("INSERT INTO sys_user(user_name,nick_name,password,create_time,pwd_update_date) VALUES(?,?,?,NOW(),NOW())", r.Username, r.Username, hash)
		})
		return saved(e)
	})
}

func Welcome() *File { return welcome() }

var demo = struct {
	sync.RWMutex
	rows map[int64]types.DemoUser
}{rows: map[int64]types.DemoUser{1: {UserId: 1, Username: "admin", Password: "admin123", Mobile: "15888888888"}, 2: {UserId: 2, Username: "ry", Password: "admin123", Mobile: "15666666666"}}}

func (s *Session) DemoList() (*types.Result, error) {
	demo.RLock()
	defer demo.RUnlock()
	rows := []types.DemoUser{}
	for _, v := range demo.rows {
		rows = append(rows, v)
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].UserId < rows[j].UserId })
	return result(rows, nil)
}

func (s *Session) DemoGet(r *types.RecordId) (*types.Result, error) {
	demo.RLock()
	defer demo.RUnlock()
	v, ok := demo.rows[r.Id]
	if !ok {
		return nil, apierror.Business("用户不存在")
	}
	return result(v, nil)
}

func (s *Session) DemoSave(r *types.DemoUser) (*types.Result, error) {
	if r.UserId <= 0 {
		return nil, apierror.Business("用户ID不能为空")
	}
	demo.Lock()
	defer demo.Unlock()
	demo.rows[r.UserId] = *r
	return result(nil, nil)
}

func (s *Session) DemoEdit(r *types.DemoUser) (*types.Result, error) {
	demo.Lock()
	defer demo.Unlock()
	if _, ok := demo.rows[r.UserId]; !ok {
		return nil, apierror.Business("用户不存在")
	}
	demo.rows[r.UserId] = *r
	return result(nil, nil)
}

func (s *Session) DemoRemove(r *types.RecordId) (*types.Result, error) {
	demo.Lock()
	defer demo.Unlock()
	if _, ok := demo.rows[r.Id]; !ok {
		return nil, apierror.Business("用户不存在")
	}
	delete(demo.rows, r.Id)
	return result(nil, nil)
}
