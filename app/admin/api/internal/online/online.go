package online

import (
	"context"
	"encoding/json"
	"net"

	"github.com/mileusna/useragent"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const Prefix = "ruoyi-go:online:"

type Identity struct {
	UserName string `db:"user_name" json:"userName"`
	DeptName string `db:"dept_name" json:"deptName"`
}
type Session struct {
	Identity
	TokenId       string `json:"tokenId"`
	Ipaddr        string `json:"ipaddr"`
	LoginLocation string `json:"loginLocation"`
	Browser       string `json:"browser"`
	OS            string `json:"os"`
	LoginTime     int64  `json:"loginTime"`
}

func Client(ip, agent string) Session {
	ua := useragent.Parse(agent)
	v := Session{Ipaddr: ip, Browser: ua.Name, OS: ua.OS}
	addr := net.ParseIP(ip)
	if addr.IsLoopback() || addr.IsPrivate() {
		v.LoginLocation = "内网IP"
	}
	return v
}

// Admit also backfills active sessions issued before the online directory was introduced.
func Admit(ctx context.Context, r *redis.Redis, v Session, expires int64) (bool, error) {
	b, e := json.Marshal(v)
	if e != nil {
		return false, e
	}
	res, e := r.EvalCtx(ctx, `if redis.call('EXISTS',KEYS[2])==1 then return 0 end;redis.call('SET',KEYS[1],ARGV[1],'EXAT',ARGV[2],'NX');return 1`, []string{Prefix + v.TokenId, "ruoyi-go:revoked:" + v.TokenId}, string(b), expires)
	return res == int64(1), e
}

// Keep the deny-list marker until the original expiry even if less than a second remains.
func Revoke(ctx context.Context, r *redis.Redis, id string, expires int64) (*Session, error) {
	raw, e := r.EvalCtx(ctx, `local v=redis.call('GET',KEYS[1]) or '';local deadline=tonumber(ARGV[1]);if deadline>0 then redis.call('SET',KEYS[2],'1','EXAT',deadline) else local ttl=redis.call('PTTL',KEYS[1]);if ttl>=0 then redis.call('SET',KEYS[2],'1','PX',math.max(1,ttl)) end end;redis.call('DEL',KEYS[1]);return v`, []string{Prefix + id, "ruoyi-go:revoked:" + id}, expires)
	if e != nil || raw == "" {
		return nil, e
	}
	var v Session
	if e = json.Unmarshal([]byte(raw.(string)), &v); e != nil {
		return nil, e
	}
	return &v, nil
}
