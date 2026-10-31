// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ruoyi-vue-go-zero/app/admin/api/internal/online"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"

	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"ruoyi-vue-go-zero/app/admin/api/internal/model"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

const captchaKey = "ruoyi-go:captcha:"
const badCredentials = apierror.Business("用户不存在/密码错误")

func (l *LoginLogic) Login(req *types.LoginRequest, ip, agent string) (*types.LoginResponse, error) {
	userID, err := l.authenticate(req, ip)
	client := online.Client(ip, agent)
	var token string
	if err == nil {
		now := time.Now().Unix()
		id := uuid.NewString()
		claims := jwt.MapClaims{"userId": userID, "jti": id, "iat": now}
		token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(l.svcCtx.Config.Auth.AccessSecret))
		if err == nil {
			client.Identity, err = l.svcCtx.Access.ActiveIdentity(l.ctx, userID)
			if err == nil && client.UserName == "" {
				err = badCredentials
			}
			if err == nil {
				err = l.svcCtx.Users.RecordLogin(l.ctx, userID, ip)
			}
			if err == nil {
				client.UserID = userID
				client.TokenId = id
				client.LoginTime = now * 1000
				err = online.Create(l.ctx, l.svcCtx.Redis, client, l.svcCtx.Config.Auth.AccessExpire)
			}
		}
	}
	l.audit(req.Username, client, err)
	if err != nil {
		return nil, err
	}
	return &types.LoginResponse{Code: 200, Msg: "操作成功", Token: token}, nil
}

func (l *LoginLogic) authenticate(req *types.LoginRequest, ip string) (int64, error) {
	s := l.svcCtx
	enabled, err := s.CaptchaEnabled(l.ctx)
	if err != nil {
		return 0, err
	}
	if enabled {
		answer, err := s.Redis.GetDelCtx(l.ctx, captchaKey+req.Uuid)
		if err != nil {
			return 0, err
		}
		if answer == "" {
			return 0, apierror.Business("验证码已失效")
		}
		if !strings.EqualFold(answer, req.Code) {
			return 0, apierror.Business("验证码错误")
		}
	}
	if err := s.Validator.StructCtx(l.ctx, req); err != nil {
		return 0, badCredentials
	}
	blocked, err := s.Settings.Value(l.ctx, "sys.login.blackIPList")
	if err != nil {
		return 0, err
	}
	if ipBlocked(ip, blocked) {
		return 0, apierror.Business("很遗憾，访问IP已被列入系统黑名单")
	}
	user, err := s.Users.FindLoginUser(l.ctx, req.Username)
	if errors.Is(err, model.ErrNotFound) {
		return 0, badCredentials
	}
	if err != nil {
		return 0, err
	}
	if user.Status != "0" {
		return 0, apierror.Business("用户已封禁，请联系管理员")
	}
	if err = l.checkPassword(user, req.Password); err != nil {
		return 0, err
	}
	return user.UserID, nil
}

// Apply the check and counter update together so concurrent attempts cannot clear an active lock.
const passwordAttempt = `
local attempts = tonumber(redis.call('GET',KEYS[1]) or '0')
if attempts >= tonumber(ARGV[1]) then return 2 end
if ARGV[3] == '1' then
 redis.call('DEL',KEYS[1])
 return 0
end
redis.call('SET',KEYS[1],attempts+1,'EX',ARGV[2])
return 1`

func (l *LoginLogic) checkPassword(user *model.LoginUser, password string) error {
	s := l.svcCtx
	key := "ruoyi-go:pwd-errors:" + strconv.FormatInt(user.UserID, 10)
	// Avoid expensive BCrypt work for already locked accounts.
	count, err := s.Redis.GetCtx(l.ctx, key)
	if err != nil {
		return err
	}
	attempts, _ := strconv.Atoi(count)
	conf := s.Config.Auth
	locked := apierror.Business(fmt.Sprintf("密码输入错误%d次，帐户锁定%d分钟", conf.MaxRetryCount, conf.LockMinutes))
	if attempts >= conf.MaxRetryCount {
		return locked
	}
	valid := 0
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil {
		valid = 1
	}
	result, err := s.Redis.EvalCtx(l.ctx, passwordAttempt, []string{key}, conf.MaxRetryCount, conf.LockMinutes*60, valid)
	if err != nil {
		return err
	}
	switch result {
	case int64(0):
		return nil
	case int64(1):
		return badCredentials
	case int64(2):
		return locked
	default:
		return fmt.Errorf("unexpected password counter result")
	}
}

func (l *LoginLogic) audit(username string, client online.Session, loginErr error) {
	status, msg := "0", "登录成功"
	if loginErr != nil {
		status, msg = "1", "系统异常"
		var business apierror.Business
		if errors.As(loginErr, &business) {
			msg = business.Error()
		}
	}
	name := []rune(username)
	if len(name) > 50 {
		name = name[:50]
	}
	_, err := l.svcCtx.LoginLog.Insert(l.ctx, &model.SysLogininfor{UserName: string(name), Ipaddr: client.Ipaddr, Browser: client.Browser, Os: client.OS, LoginLocation: client.LoginLocation, Status: status, Msg: msg, LoginTime: sql.NullTime{Time: time.Now(), Valid: true}})
	if err != nil {
		l.Errorf("记录登录日志失败: %v", err)
	}
}
