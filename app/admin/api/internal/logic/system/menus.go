package system

import (
	"net/url"
	"strings"

	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

const menuColumns = `m.menu_id,m.parent_id,m.menu_name,m.order_num,m.path,COALESCE(m.component,'') component,
 COALESCE(m.query,'') AS route_query,COALESCE(m.route_name,'') route_name,m.is_frame,m.is_cache,m.menu_type,m.visible,m.status,
 COALESCE(m.perms,'') perms,COALESCE(m.icon,'') icon,DATE_FORMAT(m.create_time,'%Y-%m-%d %H:%i:%s') create_time`

func (s *Session) menuFilter() filter {
	f := filter{}
	if s.user != 1 {
		f.add("EXISTS(SELECT 1 FROM sys_role_menu rm WHERE rm.menu_id=m.menu_id AND rm.role_id IN (" + numbers(s.roleIDs) + "))")
	}
	return f
}

func (s *Session) checkMenu(id int64) error {
	f := s.menuFilter()
	f.add("m.menu_id=?", id)
	return s.require("SELECT COUNT(*) FROM sys_menu m"+f.where(), f.args...)
}

func (s *Session) menuList(r *types.ListFilter) ([]types.MenuRecord, error) {
	f := s.menuFilter()
	f.text("m.menu_name", r.MenuName)
	f.equal("m.status", r.Status)
	return many[types.MenuRecord](s, "SELECT "+menuColumns+" FROM sys_menu m"+f.where()+" ORDER BY m.parent_id,m.order_num,m.menu_id", f.args...)
}

func (s *Session) ListMenu(r *types.ListFilter) (*types.MenuList, error) {
	rows, e := s.menuList(r)
	return &types.MenuList{StatusResponse: ok(), Data: rows}, e
}

func (s *Session) GetMenu(r *types.RecordId) (*types.MenuResult, error) {
	if e := s.checkMenu(r.Id); e != nil {
		return nil, e
	}
	v, e := one[types.MenuRecord](s, "SELECT "+menuColumns+" FROM sys_menu m WHERE menu_id=?", r.Id)
	return &types.MenuResult{StatusResponse: ok(), Data: v}, e
}

func (s *Session) MenuTree() (*types.TreeResult, error) {
	rows, e := s.menuList(&types.ListFilter{})
	if e != nil {
		return nil, e
	}
	nodes := []node{}
	for _, r := range rows {
		nodes = append(nodes, node{r.MenuId, r.ParentId, r.MenuName, false})
	}
	return &types.TreeResult{StatusResponse: ok(), Data: tree(nodes)}, nil
}

func (s *Session) RoleMenuTree(r *types.RecordId) (*types.MenuTreeResult, error) {
	if e := s.roleAllowed(r.Id, false); e != nil {
		return nil, e
	}
	t, e := s.MenuTree()
	if e != nil {
		return nil, e
	}
	ids, e := many[int64](s, `SELECT rm.menu_id FROM sys_role_menu rm JOIN sys_role r ON rm.role_id=r.role_id WHERE rm.role_id=?
 AND (r.menu_check_strictly=0 OR NOT EXISTS(SELECT 1 FROM sys_menu m JOIN sys_role_menu child ON child.menu_id=m.menu_id WHERE child.role_id=r.role_id AND m.parent_id=rm.menu_id))`, r.Id)
	return &types.MenuTreeResult{StatusResponse: ok(), Menus: t.Data, CheckedKeys: ids}, e
}

func (s *Session) writeMenu(r *types.MenuRecord, create bool) (*types.StatusResponse, error) {
	if strings.TrimSpace(r.MenuName) == "" || !validStatus(r.Status) || (r.MenuType != "M" && r.MenuType != "C" && r.MenuType != "F") || r.OrderNum < 0 {
		return nil, apierror.Business("菜单名称、类型、排序或状态无效")
	}
	if r.MenuType != "F" && r.IsFrame == "0" {
		u, e := url.Parse(r.Path)
		if e != nil || u.Host == "" || (u.Scheme != "https" && u.Scheme != "http") {
			return nil, apierror.Business("外链必须是有效的 HTTP 或 HTTPS 地址")
		}
	}
	id := r.MenuId
	if create {
		id = 0
	} else if e := s.checkMenu(id); e != nil {
		return nil, e
	}
	if s.user != 1 {
		for _, key := range strings.Split(r.Perms, ",") {
			k := strings.TrimSpace(key)
			if k != "" && !s.grants[k] && !s.grants["*:*:*"] {
				return nil, apierror.Forbidden
			}
		}
	}
	if r.ParentId != 0 {
		if e := s.checkMenu(r.ParentId); e != nil {
			return nil, e
		}
		parent := r.ParentId
		seen := map[int64]bool{}
		for parent > 0 {
			if parent == id || seen[parent] {
				return nil, apierror.Business("上级菜单不能是自身或子菜单")
			}
			seen[parent] = true
			next, e := one[int64](s, "SELECT parent_id FROM sys_menu WHERE menu_id=?", parent)
			if e != nil {
				return nil, e
			}
			parent = next
		}
	}
	if e := s.rejectIf("SELECT COUNT(*) FROM sys_menu WHERE parent_id=? AND menu_name=? AND menu_id<>?", "同级菜单名称已存在", r.ParentId, r.MenuName, id); e != nil {
		return nil, e
	}
	args := []any{r.MenuName, r.ParentId, r.OrderNum, r.Path, r.Component, r.Query, r.RouteName, r.IsFrame, r.IsCache, r.MenuType, r.Visible, r.Status, r.Perms, r.Icon, s.name}
	if create {
		return saved(s.exec("INSERT INTO sys_menu(menu_name,parent_id,order_num,path,component,`query`,route_name,is_frame,is_cache,menu_type,visible,status,perms,icon,create_by,create_time) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,NOW())", args...))
	}
	return saved(s.exec("UPDATE sys_menu SET menu_name=?,parent_id=?,order_num=?,path=?,component=?,`query`=?,route_name=?,is_frame=?,is_cache=?,menu_type=?,visible=?,status=?,perms=?,icon=?,update_by=?,update_time=NOW() WHERE menu_id=?", append(args, id)...))
}

func (s *Session) CreateMenu(r *types.MenuRecord) (*types.StatusResponse, error) {
	return s.writeMenu(r, true)
}

func (s *Session) EditMenu(r *types.MenuRecord) (*types.StatusResponse, error) {
	return s.writeMenu(r, false)
}

func (s *Session) DeleteMenu(r *types.RecordIds) (*types.StatusResponse, error) {
	ids, e := parseIDs(r.Ids)
	if e != nil {
		return nil, e
	}
	for _, id := range ids {
		if e := s.checkMenu(id); e != nil {
			return nil, e
		}
		if e := s.rejectIf("SELECT COUNT(*) FROM sys_menu WHERE parent_id=?", "存在子菜单", id); e != nil {
			return nil, e
		}
		if e := s.rejectIf("SELECT COUNT(*) FROM sys_role_menu WHERE menu_id=?", "菜单已分配，不能删除", id); e != nil {
			return nil, e
		}
		if e := s.exec("DELETE FROM sys_menu WHERE menu_id=?", id); e != nil {
			return nil, e
		}
	}
	return saved(nil)
}
