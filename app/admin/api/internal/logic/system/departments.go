package system

import (
	"slices"
	"strconv"
	"strings"

	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"ruoyi-vue-go-zero/app/admin/api/internal/model"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func (s *Session) checkDept(id int64) error {
	if !s.deptAllowed(id) {
		return apierror.Forbidden
	}
	return s.require("SELECT COUNT(*) FROM sys_dept WHERE dept_id=? AND del_flag='0'", id)
}

func (s *Session) departmentList(r *types.ListFilter, exclude int64) ([]types.DeptRecord, error) {
	return s.models().DepartmentRows(s.ctx, r, exclude, s.all, s.departments)
}

func (s *Session) ListDept(r *types.ListFilter) (*types.DeptList, error) {
	rows, e := s.departmentList(r, 0)
	return &types.DeptList{StatusResponse: ok(), Data: rows}, e
}

func (s *Session) ExcludeDepartment(r *types.RecordId) (*types.DeptList, error) {
	rows, e := s.departmentList(&types.ListFilter{}, r.Id)
	return &types.DeptList{StatusResponse: ok(), Data: rows}, e
}

func (s *Session) GetDept(r *types.RecordId) (*types.DeptResult, error) {
	if err := s.checkDept(r.Id); err != nil {
		return nil, err
	}
	d, e := one[types.DeptRecord](s, "SELECT "+deptColumns+" FROM sys_dept WHERE dept_id=?", r.Id)
	return &types.DeptResult{StatusResponse: ok(), Data: d}, e
}

type node struct {
	id, parent int64
	label      string
	disabled   bool
}

func tree(rows []node) []types.TreeNode {
	ids := map[int64]bool{}
	for _, r := range rows {
		ids[r.id] = true
	}
	groups := map[int64][]node{}
	for _, r := range rows {
		parent := r.parent
		if !ids[parent] {
			parent = 0
		}
		groups[parent] = append(groups[parent], r)
	}
	var children func(int64) []types.TreeNode
	children = func(p int64) []types.TreeNode {
		out := []types.TreeNode{}
		for _, r := range groups[p] {
			out = append(out, types.TreeNode{Id: r.id, Label: r.label, Disabled: r.disabled, Children: children(r.id)})
		}
		return out
	}
	return children(0)
}

func (s *Session) departmentTree() ([]types.TreeNode, error) {
	rows, e := s.departmentList(&types.ListFilter{}, 0)
	if e != nil {
		return nil, e
	}
	nodes := []node{}
	for _, r := range rows {
		nodes = append(nodes, node{r.DeptId, r.ParentId, r.DeptName, r.Status != "0"})
	}
	return tree(nodes), nil
}

func (s *Session) UserDepartments() (*types.TreeResult, error) {
	v, e := s.departmentTree()
	return &types.TreeResult{StatusResponse: ok(), Data: v}, e
}

func (s *Session) writeDept(r *types.DeptRecord, create bool) (*types.StatusResponse, error) {
	if strings.TrimSpace(r.DeptName) == "" || !validStatus(r.Status) || r.OrderNum < 0 {
		return nil, apierror.Business("部门名称、状态或排序无效")
	}
	id := r.DeptId
	if create {
		id = 0
	} else if err := s.checkDept(id); err != nil {
		return nil, err
	}
	ancestors := "0"
	if r.ParentId != 0 {
		if err := s.checkDept(r.ParentId); err != nil {
			return nil, err
		}
		p, err := one[types.DeptRecord](s, "SELECT "+deptColumns+" FROM sys_dept WHERE dept_id=?", r.ParentId)
		if err != nil {
			return nil, err
		}
		if p.Status != "0" {
			return nil, apierror.Business("上级部门已停用")
		}
		ancestors = p.Ancestors + "," + strconv.FormatInt(p.DeptId, 10)
		if !create && slices.Contains(strings.Split(ancestors, ","), strconv.FormatInt(id, 10)) {
			return nil, apierror.Business("上级部门不能是自身或子部门")
		}
	} else if !s.all {
		return nil, apierror.Forbidden
	}
	if err := s.rejectIf("SELECT COUNT(*) FROM sys_dept WHERE del_flag='0' AND parent_id=? AND dept_name=? AND dept_id<>?", "同级部门名称已存在", r.ParentId, r.DeptName, id); err != nil {
		return nil, err
	}
	if create {
		return saved(s.exec("INSERT INTO sys_dept(parent_id,ancestors,dept_name,order_num,leader,phone,email,status,create_by,create_time) VALUES(?,?,?,?,?,?,?,?,?,NOW())", r.ParentId, ancestors, r.DeptName, r.OrderNum, r.Leader, r.Phone, r.Email, r.Status, s.name))
	}
	if r.Status == "1" {
		if err := s.rejectIf("SELECT COUNT(*) FROM sys_dept WHERE del_flag='0' AND status='0' AND FIND_IN_SET(?,ancestors)", "存在启用的下级部门", id); err != nil {
			return nil, err
		}
	}
	old, err := one[string](s, "SELECT ancestors FROM sys_dept WHERE dept_id=?", id)
	if err != nil {
		return nil, err
	}
	if old != ancestors {
		if !s.all {
			if err := s.rejectIf("SELECT COUNT(*) FROM sys_dept WHERE del_flag='0' AND FIND_IN_SET(?,ancestors) AND dept_id NOT IN ("+numbers(s.departments)+")", "下级部门超出数据权限", id); err != nil {
				return nil, err
			}
		}
		if err := s.exec("UPDATE sys_dept SET ancestors=CONCAT(?,SUBSTRING(ancestors,CHAR_LENGTH(?)+1)) WHERE FIND_IN_SET(?,ancestors)", ancestors, old, id); err != nil {
			return nil, err
		}
	}
	return saved(s.exec("UPDATE sys_dept SET parent_id=?,ancestors=?,dept_name=?,order_num=?,leader=?,phone=?,email=?,status=?,update_by=?,update_time=NOW() WHERE dept_id=?", r.ParentId, ancestors, r.DeptName, r.OrderNum, r.Leader, r.Phone, r.Email, r.Status, s.name, id))
}

func (s *Session) CreateDept(r *types.DeptRecord) (*types.StatusResponse, error) {
	return s.writeDept(r, true)
}

func (s *Session) EditDept(r *types.DeptRecord) (*types.StatusResponse, error) {
	return s.writeDept(r, false)
}

func (s *Session) DeleteDept(r *types.RecordIds) (*types.StatusResponse, error) {
	ids, e := parseIDs(r.Ids)
	if e != nil {
		return nil, e
	}
	for _, id := range ids {
		if err := s.checkDept(id); err != nil {
			return nil, err
		}
		if err := s.rejectIf("SELECT COUNT(*) FROM sys_dept WHERE parent_id=? AND del_flag='0'", "存在下级部门", id); err != nil {
			return nil, err
		}
		if err := s.rejectIf("SELECT COUNT(*) FROM sys_user WHERE dept_id=? AND del_flag='0'", "部门存在用户", id); err != nil {
			return nil, err
		}
		if err := s.exec("UPDATE sys_dept SET del_flag='2',update_by=?,update_time=NOW() WHERE dept_id=?", s.name, id); err != nil {
			return nil, err
		}
	}
	return saved(nil)
}

func (s *Session) sortRecords(raw, orders string, dept bool) (*types.StatusResponse, error) {
	ids, e := parseIDs(raw)
	if e != nil {
		return nil, e
	}
	sorts := strings.Split(orders, ",")
	if len(ids) != len(sorts) {
		return nil, apierror.Business("排序参数无效")
	}
	for i, id := range ids {
		n, e := strconv.Atoi(sorts[i])
		if e != nil || n < 0 {
			return nil, apierror.Business("排序值无效")
		}
		if dept {
			e = s.checkDept(id)
		} else {
			e = s.checkMenu(id)
		}
		if e != nil {
			return nil, e
		}
		q := "UPDATE sys_menu SET order_num=? WHERE menu_id=?"
		if dept {
			q = "UPDATE sys_dept SET order_num=? WHERE dept_id=?"
		}
		if e = s.exec(q, n, id); e != nil {
			return nil, e
		}
	}
	return saved(nil)
}

func (s *Session) SortDepartments(r *types.SortInput) (*types.StatusResponse, error) {
	return s.sortRecords(r.DeptIds, r.OrderNums, true)
}

func (s *Session) SortMenus(r *types.SortInput) (*types.StatusResponse, error) {
	return s.sortRecords(r.MenuIds, r.OrderNums, false)
}

const deptColumns = model.DepartmentColumns
