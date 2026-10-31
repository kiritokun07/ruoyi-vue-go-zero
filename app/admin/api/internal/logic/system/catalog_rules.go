package system

import (
	"regexp"
	"strings"

	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"ruoyi-vue-go-zero/app/admin/api/internal/model"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
	"github.com/microcosm-cc/bluemonday"
)

func result(value any, err error) (*types.Result, error) {
	return &types.Result{StatusResponse: ok(), Data: value}, err
}

var dictKeyPattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

func (s *Session) ValidateParameter(r *types.ParameterEntry, create bool) error {
	if strings.TrimSpace(r.ConfigName) == "" || strings.TrimSpace(r.ConfigKey) == "" || r.ConfigValue == "" {
		return apierror.Business("参数名称、键名和键值不能为空")
	}
	if r.ConfigType == "" {
		r.ConfigType = "N"
	}
	if r.ConfigType != "Y" && r.ConfigType != "N" {
		return apierror.Business("参数类型错误")
	}
	return nil
}

func (s *Session) ValidateDictType(r *types.DictTypeEntry, create bool) error {
	if r.DictName == "" || !dictKeyPattern.MatchString(r.DictType) || !validStatus(r.Status) {
		return apierror.Business("字典名称、类型或状态无效")
	}
	if !create {
		return s.exec("UPDATE sys_dict_data SET dict_type=? WHERE dict_type=(SELECT dict_type FROM sys_dict_type WHERE dict_id=?)", r.DictType, r.DictId)
	}
	return nil
}

func (s *Session) ValidateDictEntry(r *types.DictEntryEntry, create bool) error {
	if r.DictLabel == "" || r.DictValue == "" || !validStatus(r.Status) {
		return apierror.Business("字典标签、键值或状态无效")
	}
	if r.IsDefault == "" {
		r.IsDefault = "N"
	}
	return s.require("SELECT COUNT(*) FROM sys_dict_type WHERE dict_type=?", r.DictType)
}

var richText = bluemonday.UGCPolicy()

func (s *Session) ValidateBulletin(r *types.BulletinEntry, create bool) error {
	if r.NoticeTitle == "" || !validStatus(r.Status) || (r.NoticeType != "1" && r.NoticeType != "2") {
		return apierror.Business("公告标题、类型或状态无效")
	}
	r.NoticeContent = richText.Sanitize(r.NoticeContent)
	return nil
}

func (s *Session) BeforeDelete(table string, ids []int64) error {
	switch table {
	case "sys_dict_type":
		return s.rejectIf("SELECT COUNT(*) FROM sys_dict_data d JOIN sys_dict_type t ON t.dict_type=d.dict_type WHERE t.dict_id IN ("+numbers(ids)+")", "字典类型已分配数据")
	case "sys_config":
		return s.rejectIf("SELECT COUNT(*) FROM sys_config WHERE config_type='Y' AND config_id IN ("+numbers(ids)+")", "系统内置参数不能删除")
	case "sys_notice":
		return s.exec("DELETE FROM sys_notice_read WHERE notice_id IN (" + numbers(ids) + ")")
	}
	return nil
}

func (s *Session) DictOptions() (*types.Result, error) {
	v, e := many[types.DictTypeEntry](s, "SELECT "+model.DictTypeColumns+" FROM sys_dict_type ORDER BY dict_id")
	return result(v, e)
}

func (s *Session) PostOptions() (*types.Result, error) {
	v, e := many[types.PostRecord](s, "SELECT "+postColumns+" FROM sys_post ORDER BY post_sort")
	return result(v, e)
}

func (s *Session) RoleOptions() (*types.Result, error) {
	f := s.roleFilter()
	v, e := many[types.RoleRecord](s, "SELECT "+roleColumns+" FROM sys_role r"+f.where()+" ORDER BY role_sort", f.args...)
	return result(v, e)
}

type reader struct {
	UserID      int64  `db:"user_id" json:"userId"`
	UserName    string `db:"user_name" json:"userName"`
	NickName    string `db:"nick_name" json:"nickName"`
	DeptName    string `db:"dept_name" json:"deptName"`
	Phonenumber string `db:"phonenumber" json:"phonenumber"`
	ReadTime    string `db:"read_time" json:"readTime"`
}

func (s *Session) ReadUsers(r *types.AdminFilter) (*types.PageResult, error) {
	f := filter{}
	f.add("nr.notice_id=?", r.NoticeId)
	f.add(s.userScope())
	if r.SearchValue != "" {
		v := "%" + r.SearchValue + "%"
		f.add("(u.user_name LIKE ? OR u.nick_name LIKE ? OR u.phonenumber LIKE ?)", v, v, v)
	}
	from := " FROM sys_notice_read nr JOIN sys_user u ON u.user_id=nr.user_id LEFT JOIN sys_dept d ON d.dept_id=u.dept_id" + f.where()
	n, e := one[int64](s, "SELECT COUNT(*)"+from, f.args...)
	if e != nil {
		return nil, e
	}
	rows, e := many[reader](s, "SELECT u.user_id,u.user_name,u.nick_name,IFNULL(d.dept_name,'') dept_name,IFNULL(u.phonenumber,'') phonenumber,DATE_FORMAT(nr.read_time,'%Y-%m-%d %H:%i:%s') read_time"+from+" ORDER BY nr.read_time DESC"+page(&types.ListFilter{PageNum: r.PageNum, PageSize: r.PageSize}), f.args...)
	return &types.PageResult{StatusResponse: ok(), Rows: rows, Total: n}, e
}
