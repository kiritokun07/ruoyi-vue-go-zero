package system

import (
	"strings"

	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

const postColumns = `post_id,post_name,post_code,post_sort,status,IFNULL(remark,'') remark,DATE_FORMAT(create_time,'%Y-%m-%d %H:%i:%s') create_time`

func (s *Session) ListPost(r *types.ListFilter) (*types.PostList, error) {
	f := filter{}
	f.text("post_code", r.PostCode)
	f.text("post_name", r.PostName)
	f.equal("status", r.Status)
	from := " FROM sys_post" + f.where()
	n, e := one[int64](s, "SELECT COUNT(*)"+from, f.args...)
	if e != nil {
		return nil, e
	}
	rows, e := many[types.PostRecord](s, "SELECT "+postColumns+from+" ORDER BY post_sort,post_id"+page(r), f.args...)
	return &types.PostList{StatusResponse: ok(), Rows: rows, Total: n}, e
}

func (s *Session) GetPost(r *types.RecordId) (*types.PostResult, error) {
	v, e := one[types.PostRecord](s, "SELECT "+postColumns+" FROM sys_post WHERE post_id=?", r.Id)
	return &types.PostResult{StatusResponse: ok(), Data: v}, e
}

func (s *Session) writePost(r *types.PostRecord, create bool) (*types.StatusResponse, error) {
	if strings.TrimSpace(r.PostCode) == "" || strings.TrimSpace(r.PostName) == "" || !validStatus(r.Status) || r.PostSort < 0 {
		return nil, apierror.Business("岗位编码、名称、排序或状态无效")
	}
	id := r.PostId
	if create {
		id = 0
	} else if err := s.require("SELECT COUNT(*) FROM sys_post WHERE post_id=?", id); err != nil {
		return nil, err
	}
	if err := s.rejectIf("SELECT COUNT(*) FROM sys_post WHERE (post_name=? OR post_code=?) AND post_id<>?", "岗位编码或名称已存在", r.PostName, r.PostCode, id); err != nil {
		return nil, err
	}
	if create {
		return saved(s.exec("INSERT INTO sys_post(post_code,post_name,post_sort,status,remark,create_by,create_time) VALUES(?,?,?,?,?,?,NOW())", r.PostCode, r.PostName, r.PostSort, r.Status, r.Remark, s.name))
	}
	return saved(s.exec("UPDATE sys_post SET post_code=?,post_name=?,post_sort=?,status=?,remark=?,update_by=?,update_time=NOW() WHERE post_id=?", r.PostCode, r.PostName, r.PostSort, r.Status, r.Remark, s.name, id))
}

func (s *Session) CreatePost(r *types.PostRecord) (*types.StatusResponse, error) {
	return s.writePost(r, true)
}

func (s *Session) EditPost(r *types.PostRecord) (*types.StatusResponse, error) {
	return s.writePost(r, false)
}

func (s *Session) DeletePost(r *types.RecordIds) (*types.StatusResponse, error) {
	ids, e := parseIDs(r.Ids)
	if e != nil {
		return nil, e
	}
	for _, id := range ids {
		if e := s.rejectIf("SELECT COUNT(*) FROM sys_user_post WHERE post_id=?", "岗位已分配，不能删除", id); e != nil {
			return nil, e
		}
		if e := s.exec("DELETE FROM sys_post WHERE post_id=?", id); e != nil {
			return nil, e
		}
	}
	return saved(nil)
}
