package system

import (
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
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
	return nil
}

func (s *Session) ValidateDictEntry(r *types.DictEntryEntry, create bool) error {
	if r.DictLabel == "" || r.DictValue == "" || !validStatus(r.Status) {
		return apierror.Business("字典标签、键值或状态无效")
	}
	if r.IsDefault == "" {
		r.IsDefault = "N"
	}
	return requireCount(s.models().DictTypeCount(s.ctx, r.DictType))
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
		n, err := s.models().DictTypeReferences(s.ctx, ids)
		if err != nil {
			return err
		}
		if n > 0 {
			return apierror.Business("字典类型已分配数据")
		}
		return nil
	case "sys_config":
		n, err := s.models().BuiltinParameterCount(s.ctx, ids)
		if err != nil {
			return err
		}
		if n > 0 {
			return apierror.Business("系统内置参数不能删除")
		}
		return nil
	}
	return nil
}

func (s *Session) DictOptions() (*types.Result, error) {
	v, e := s.models().DictOptions(s.ctx)
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
