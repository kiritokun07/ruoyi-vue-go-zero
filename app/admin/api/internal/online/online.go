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
	UserID        int64  `json:"userId"`
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

// Create is only called after successful authentication. Requests never recreate sessions.
func Create(ctx context.Context, r *redis.Redis, v Session, ttl int64) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return r.SetexCtx(ctx, Prefix+v.TokenId, string(b), int(ttl))
}

// Touch atomically checks membership and refreshes sessions close to expiry.
// Refresh and revoke cannot race to resurrect a deleted session.
func Touch(ctx context.Context, r *redis.Redis, id string, userID, ttl int64) (bool, error) {
	threshold := min(int64(20*60), ttl)
	value, err := r.EvalCtx(ctx, `local raw=redis.call('GET',KEYS[1]);if not raw or redis.call('EXISTS',KEYS[2])==1 then return 0 end;local s=cjson.decode(raw);if s.userId~=tonumber(ARGV[1]) then return 0 end;local remaining=redis.call('TTL',KEYS[1]);if remaining<0 then return 0 end;if remaining<=tonumber(ARGV[3]) then redis.call('EXPIRE',KEYS[1],ARGV[2]) end;return 1`, []string{Prefix + id, "ruoyi-go:revoked:" + id}, userID, ttl, threshold)
	return value == int64(1), err
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
