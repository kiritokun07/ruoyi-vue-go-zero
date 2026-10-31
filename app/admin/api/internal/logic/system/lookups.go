package system

import (
	"encoding/json"

	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func (s *Session) Dictionary(r *types.LookupKey) (*types.DictResult, error) {
	key := s.lookupKey(dictPrefix, r.Key)
	if key != "" {
		if raw, e := s.services.Redis.GetCtx(s.ctx, key); e == nil && raw != "" {
			var rows []types.DictRecord
			if json.Unmarshal([]byte(raw), &rows) == nil {
				return &types.DictResult{StatusResponse: ok(), Data: rows}, nil
			}
		}
	}
	rows, err := many[types.DictRecord](s, `SELECT dict_code,dict_sort,dict_label,dict_value,dict_type,COALESCE(css_class,'') css_class,COALESCE(list_class,'') list_class,is_default,status FROM sys_dict_data WHERE dict_type=? AND status='0' ORDER BY dict_sort,dict_code`, r.Key)
	if err == nil && key != "" {
		if raw, e := json.Marshal(rows); e == nil {
			_ = s.services.Redis.SetexCtx(s.ctx, key, string(raw), 300)
		}
	}
	return &types.DictResult{StatusResponse: ok(), Data: rows}, err
}

func (s *Session) ConfigValue(r *types.LookupKey) (*types.StatusResponse, error) {
	key := s.lookupKey(configPrefix, r.Key)
	if key != "" {
		if raw, e := s.services.Redis.GetCtx(s.ctx, key); e == nil && raw != "" {
			return &types.StatusResponse{Code: 200, Msg: raw}, nil
		}
	}
	v, err := one[string](s, "SELECT COALESCE((SELECT config_value FROM sys_config WHERE config_key=? ORDER BY config_id LIMIT 1),'')", r.Key)
	if err == nil && key != "" && v != "" {
		_ = s.services.Redis.SetexCtx(s.ctx, key, v, 300)
	}
	return &types.StatusResponse{Code: 200, Msg: v}, err
}

func (s *Session) lookupKey(prefix, key string) string {
	generation, e := s.services.Redis.GetCtx(s.ctx, "ruoyi-go:generation:"+prefix)
	if e != nil {
		return ""
	}
	return prefix + generation + ":" + key
}
