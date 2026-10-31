package system

import (
	"slices"
	"strings"

	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

const roleColumns = `r.role_id,r.role_name,r.role_key,r.role_sort,r.data_scope,r.status,IFNULL(r.remark,'') remark,
 r.menu_check_strictly,r.dept_check_strictly,DATE_FORMAT(r.create_time,'%Y-%m-%d %H:%i:%s') create_time,r.role_id=1 admin`

func (s *Session) roleScope() string {
	if s.all {
		return "1=1"
	}
	return "(r.create_by=? OR EXISTS (SELECT 1 FROM sys_user_role ur JOIN sys_user u ON u.user_id=ur.user_id WHERE ur.role_id=r.role_id AND u.del_flag='0' AND " + s.userScope() + "))"
}

func (s *Session) roleFilter() filter {
	f := filter{}
	f.add("r.del_flag='0'")
	if !s.all {
		f.add(s.roleScope(), s.name)
	}
	return f
}

func (s *Session) ListRole(r *types.ListFilter) (*types.RoleList, error) {
	f := s.roleFilter()
	f.text("r.role_name", r.RoleName)
	f.text("r.role_key", r.RoleKey)
	f.equal("r.status", r.Status)
	f.dates("r.create_time", r)
	from := " FROM sys_role r" + f.where()
	n, e := one[int64](s, "SELECT COUNT(*)"+from, f.args...)
	if e != nil {
		return nil, e
	}
	rows, e := many[types.RoleRecord](s, "SELECT "+roleColumns+from+" ORDER BY r.role_sort,r.role_id"+page(r), f.args...)
	return &types.RoleList{StatusResponse: ok(), Rows: rows, Total: n}, e
}

func (s *Session) roleAllowed(id int64, write bool) error {
	if write && id == 1 {
		return apierror.Business("不允许操作超级管理员角色")
	}
	f := s.roleFilter()
	f.add("r.role_id=?", id)
	if err := s.require("SELECT COUNT(*) FROM sys_role r"+f.where(), f.args...); err != nil {
		return err
	}
	if write && !s.all {
		return s.rejectIf("SELECT COUNT(*) FROM sys_user_role ur JOIN sys_user u ON u.user_id=ur.user_id WHERE ur.role_id=? AND u.del_flag='0' AND NOT "+s.userScope(), "角色已分配给范围外用户", id)
	}
	return nil
}

func (s *Session) GetRole(r *types.RecordId) (*types.RoleResult, error) {
	if e := s.roleAllowed(r.Id, false); e != nil {
		return nil, e
	}
	v, e := one[types.RoleRecord](s, "SELECT "+roleColumns+" FROM sys_role r WHERE role_id=?", r.Id)
	return &types.RoleResult{StatusResponse: ok(), Data: v}, e
}

func (s *Session) roleOptions(user int64) ([]types.RoleRecord, error) {
	f := s.roleFilter()
	f.add("r.status='0'")
	if user != 1 {
		f.add("r.role_id<>1")
	}
	rows, e := many[types.RoleRecord](s, "SELECT "+roleColumns+",EXISTS(SELECT 1 FROM sys_user_role ur WHERE ur.role_id=r.role_id AND ur.user_id=?) flag FROM sys_role r"+f.where()+" ORDER BY role_sort", append([]any{user}, f.args...)...)
	return rows, e
}

func (s *Session) checkAssignments(ids []int64) error {
	for _, id := range ids {
		if id == 1 {
			return apierror.Forbidden
		}
		if e := s.roleAllowed(id, false); e != nil {
			return e
		}
		role, e := one[types.RoleRecord](s, "SELECT "+roleColumns+" FROM sys_role r WHERE r.role_id=?", id)
		if e != nil {
			return e
		}
		if role.Status != "0" {
			return apierror.Business("角色已停用")
		}
		if s.user != 1 {
			menus, e := many[int64](s, "SELECT menu_id FROM sys_role_menu WHERE role_id=?", id)
			if e != nil {
				return e
			}
			for _, m := range menus {
				if e := s.checkMenu(m); e != nil {
					return e
				}
			}
			if !s.all && !slices.Contains(s.roleIDs, id) && role.DataScope != "5" {
				if role.DataScope != "2" {
					return apierror.Forbidden
				}
				depts, e := many[int64](s, "SELECT dept_id FROM sys_role_dept WHERE role_id=?", id)
				if e != nil {
					return e
				}
				for _, d := range depts {
					if !s.deptAllowed(d) {
						return apierror.Forbidden
					}
				}
			}
		}
	}
	return nil
}

func (s *Session) writeRole(r *types.RoleRecord, create bool) (*types.StatusResponse, error) {
	if strings.TrimSpace(r.RoleName) == "" || strings.TrimSpace(r.RoleKey) == "" || !validStatus(r.Status) || r.RoleSort < 0 {
		return nil, apierror.Business("角色名称、权限字符、排序或状态无效")
	}
	for _, k := range strings.Split(r.RoleKey, ",") {
		if strings.TrimSpace(k) == "admin" {
			return nil, apierror.Forbidden
		}
	}
	id := r.RoleId
	if create {
		id = 0
	} else if e := s.roleAllowed(id, true); e != nil {
		return nil, e
	}
	if e := s.rejectIf("SELECT COUNT(*) FROM sys_role WHERE del_flag='0' AND (role_name=? OR role_key=?) AND role_id<>?", "角色名称或权限字符已存在", r.RoleName, r.RoleKey, id); e != nil {
		return nil, e
	}
	for _, m := range r.MenuIds {
		if e := s.checkMenu(m); e != nil {
			return nil, e
		}
	}
	if create {
		scope := "1"
		if !s.all {
			scope = "5"
		}
		res, e := s.db.ExecCtx(s.ctx, "INSERT INTO sys_role(role_name,role_key,role_sort,data_scope,menu_check_strictly,status,remark,create_by,create_time) VALUES(?,?,?,?,?,?,?,?,NOW())", r.RoleName, r.RoleKey, r.RoleSort, scope, r.MenuCheckStrictly, r.Status, r.Remark, s.name)
		if e != nil {
			return nil, e
		}
		id, e = res.LastInsertId()
		if e != nil {
			return nil, e
		}
	} else if e := s.exec("UPDATE sys_role SET role_name=?,role_key=?,role_sort=?,menu_check_strictly=?,status=?,remark=?,update_by=?,update_time=NOW() WHERE role_id=?", r.RoleName, r.RoleKey, r.RoleSort, r.MenuCheckStrictly, r.Status, r.Remark, s.name, id); e != nil {
		return nil, e
	}
	return saved(s.replaceLinks("sys_role_menu", "role_id", "menu_id", id, r.MenuIds))
}

func (s *Session) CreateRole(r *types.RoleRecord) (*types.StatusResponse, error) {
	return s.writeRole(r, true)
}

func (s *Session) EditRole(r *types.RoleRecord) (*types.StatusResponse, error) {
	return s.writeRole(r, false)
}

func (s *Session) RoleStatus(r *types.RoleRecord) (*types.StatusResponse, error) {
	if e := s.roleAllowed(r.RoleId, true); e != nil {
		return nil, e
	}
	if !validStatus(r.Status) {
		return nil, apierror.Business("角色状态无效")
	}
	return saved(s.exec("UPDATE sys_role SET status=?,update_by=?,update_time=NOW() WHERE role_id=?", r.Status, s.name, r.RoleId))
}

func (s *Session) RoleScope(r *types.RoleRecord) (*types.StatusResponse, error) {
	if e := s.roleAllowed(r.RoleId, true); e != nil {
		return nil, e
	}
	if !slices.Contains([]string{"1", "2", "3", "4", "5"}, r.DataScope) {
		return nil, apierror.Business("数据范围无效")
	}
	if !s.all && r.DataScope != "2" && r.DataScope != "5" {
		return nil, apierror.Forbidden
	}
	ids := []int64{}
	if r.DataScope == "2" {
		ids = r.DeptIds
		for _, d := range ids {
			if e := s.checkDept(d); e != nil {
				return nil, e
			}
		}
	}
	if e := s.exec("UPDATE sys_role SET data_scope=?,dept_check_strictly=?,update_by=?,update_time=NOW() WHERE role_id=?", r.DataScope, r.DeptCheckStrictly, s.name, r.RoleId); e != nil {
		return nil, e
	}
	return saved(s.replaceLinks("sys_role_dept", "role_id", "dept_id", r.RoleId, ids))
}

func (s *Session) RoleDepartments(r *types.RecordId) (*types.DeptTreeResult, error) {
	if e := s.roleAllowed(r.Id, false); e != nil {
		return nil, e
	}
	tree, e := s.departmentTree()
	if e != nil {
		return nil, e
	}
	ids, e := many[int64](s, `SELECT rd.dept_id FROM sys_role_dept rd JOIN sys_role r ON rd.role_id=r.role_id WHERE rd.role_id=?
 AND (r.dept_check_strictly=0 OR NOT EXISTS(SELECT 1 FROM sys_dept d JOIN sys_role_dept child ON child.dept_id=d.dept_id WHERE child.role_id=r.role_id AND d.parent_id=rd.dept_id))`, r.Id)
	return &types.DeptTreeResult{StatusResponse: ok(), Depts: tree, CheckedKeys: ids}, e
}

func (s *Session) DeleteRole(r *types.RecordIds) (*types.StatusResponse, error) {
	ids, e := parseIDs(r.Ids)
	if e != nil {
		return nil, e
	}
	for _, id := range ids {
		if e := s.roleAllowed(id, true); e != nil {
			return nil, e
		}
		if e := s.rejectIf("SELECT COUNT(*) FROM sys_user_role WHERE role_id=?", "角色已分配给用户", id); e != nil {
			return nil, e
		}
		for _, t := range []string{"sys_role_menu", "sys_role_dept"} {
			if e := s.exec("DELETE FROM "+t+" WHERE role_id=?", id); e != nil {
				return nil, e
			}
		}
		if e := s.exec("UPDATE sys_role SET del_flag='2',update_by=?,update_time=NOW() WHERE role_id=?", s.name, id); e != nil {
			return nil, e
		}
	}
	return saved(nil)
}

func (s *Session) members(role int64, ids []int64, add bool) (*types.StatusResponse, error) {
	if e := s.roleAllowed(role, true); e != nil {
		return nil, e
	}
	if add {
		if e := s.checkAssignments([]int64{role}); e != nil {
			return nil, e
		}
	}
	for _, id := range ids {
		if e := s.userAllowed(id, true); e != nil {
			return nil, e
		}
		q := "DELETE FROM sys_user_role WHERE user_id=? AND role_id=?"
		if add {
			q = "INSERT INTO sys_user_role(user_id,role_id) VALUES(?,?) ON DUPLICATE KEY UPDATE user_id=user_id"
		}
		if e := s.exec(q, id, role); e != nil {
			return nil, e
		}
	}
	return saved(nil)
}

func (s *Session) AddMembers(r *types.MemberAssignment) (*types.StatusResponse, error) {
	ids, e := parseIDs(r.UserIds)
	if e != nil {
		return nil, e
	}
	return s.members(r.RoleId, ids, true)
}

func (s *Session) RemoveMembers(r *types.MemberAssignment) (*types.StatusResponse, error) {
	ids, e := parseIDs(r.UserIds)
	if e != nil {
		return nil, e
	}
	return s.members(r.RoleId, ids, false)
}

func (s *Session) RemoveMember(r *types.MemberRemoval) (*types.StatusResponse, error) {
	return s.members(r.RoleId, []int64{r.UserId}, false)
}
