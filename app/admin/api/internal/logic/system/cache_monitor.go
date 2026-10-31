package system

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"

	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"ruoyi-vue-go-zero/app/admin/api/internal/online"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

const configPrefix = "ruoyi-go:config:"
const dictPrefix = "ruoyi-go:dict:"

func (s *Session) keys(prefix string) ([]string, error) {
	keys := map[string]bool{}
	var cursor uint64
	for {
		batch, next, e := s.services.Redis.ScanCtx(s.ctx, cursor, prefix+"*", 128)
		if e != nil {
			return nil, e
		}
		for _, k := range batch {
			keys[k] = true
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	out := make([]string, 0, len(keys))
	for k := range keys {
		out = append(out, k)
	}
	sort.Strings(out)
	return out, nil
}

func (s *Session) OnlineUsers(r *types.AdminFilter) (*types.PageResult, error) {
	keys, e := s.keys(online.Prefix)
	if e != nil {
		return nil, e
	}
	rows := []online.Session{}
	for _, k := range keys {
		v, e := s.services.Redis.GetCtx(s.ctx, k)
		if e != nil {
			return nil, e
		}
		if v == "" {
			continue
		}
		var o online.Session
		if e = json.Unmarshal([]byte(v), &o); e != nil {
			return nil, e
		}
		if (r.UserName == "" || r.UserName == o.UserName) && (r.Ipaddr == "" || r.Ipaddr == o.Ipaddr) {
			rows = append(rows, o)
		}
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].LoginTime > rows[j].LoginTime })
	return &types.PageResult{StatusResponse: ok(), Rows: rows, Total: int64(len(rows))}, nil
}

func (s *Session) ForceLogout(r *types.LookupKey) (*types.StatusResponse, error) {
	_, err := online.Revoke(s.ctx, s.services.Redis, r.Key, 0)
	return saved(err)
}

func (s *Session) UnlockLogin(r *types.LookupKey) (*types.StatusResponse, error) {
	ids, e := many[int64](s, "SELECT user_id FROM sys_user WHERE user_name=? AND del_flag='0'", r.Key)
	if e != nil {
		return nil, e
	}
	for _, id := range ids {
		if _, e = s.services.Redis.DelCtx(s.ctx, "ruoyi-go:pwd-errors:"+strconv.FormatInt(id, 10)); e != nil {
			return nil, e
		}
	}
	return saved(nil)
}

func validCache(prefix string) bool { return prefix == dictPrefix || prefix == configPrefix }
func (s *Session) CacheNames() (*types.Result, error) {
	return result([]map[string]string{{"cacheName": dictPrefix, "remark": "系统字典"}, {"cacheName": configPrefix, "remark": "参数设置"}}, nil)
}

func (s *Session) CacheKeys(r *types.LookupKey) (*types.Result, error) {
	if !validCache(r.Key) {
		return nil, apierror.Forbidden
	}
	v, e := s.keys(r.Key)
	return result(v, e)
}

func (s *Session) CacheValue(r *types.KeyPair) (*types.Result, error) {
	if !validCache(r.Name) || !strings.HasPrefix(r.Key, r.Name) {
		return nil, apierror.Forbidden
	}
	v, e := s.services.Redis.GetCtx(s.ctx, r.Key)
	return result(map[string]string{"cacheName": r.Name, "cacheKey": r.Key, "cacheValue": v}, e)
}

func (s *Session) ClearCacheKey(r *types.LookupKey) (*types.StatusResponse, error) {
	if !strings.HasPrefix(r.Key, dictPrefix) && !strings.HasPrefix(r.Key, configPrefix) {
		return nil, apierror.Forbidden
	}
	_, e := s.services.Redis.DelCtx(s.ctx, r.Key)
	return saved(e)
}

func (s *Session) clear(prefix string) error {
	if !validCache(prefix) {
		return apierror.Forbidden
	}
	if _, e := s.services.Redis.IncrCtx(s.ctx, "ruoyi-go:generation:"+prefix); e != nil {
		return e
	}
	keys, e := s.keys(prefix)
	if e != nil || len(keys) == 0 {
		return e
	}
	_, e = s.services.Redis.DelCtx(s.ctx, keys...)
	return e
}

func (s *Session) ClearCacheName(r *types.LookupKey) (*types.StatusResponse, error) {
	return saved(s.clear(r.Key))
}

func (s *Session) ClearCaches() (*types.StatusResponse, error) {
	if e := s.clear(configPrefix); e != nil {
		return nil, e
	}
	return saved(s.clear(dictPrefix))
}

func (s *Session) RefreshConfig() (*types.StatusResponse, error) { return saved(s.clear(configPrefix)) }
func (s *Session) RefreshDictionary() (*types.StatusResponse, error) {
	return saved(s.clear(dictPrefix))
}

func (s *Session) CacheInfo() (*types.Result, error) {
	raw, e := s.services.Redis.EvalCtx(s.ctx, "return redis.call('INFO','all')", nil)
	if e != nil {
		return nil, e
	}
	size, e := s.services.Redis.EvalCtx(s.ctx, "return redis.call('DBSIZE')", nil)
	if e != nil {
		return nil, e
	}
	info := map[string]string{}
	stats := []map[string]string{}
	for _, line := range strings.Split(raw.(string), "\r\n") {
		k, v, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		info[k] = v
		if strings.HasPrefix(k, "cmdstat_") {
			calls, _, _ := strings.Cut(strings.TrimPrefix(v, "calls="), ",")
			stats = append(stats, map[string]string{"name": strings.TrimPrefix(k, "cmdstat_"), "value": calls})
		}
	}
	return result(map[string]any{"info": info, "dbSize": size, "commandStats": stats}, nil)
}
