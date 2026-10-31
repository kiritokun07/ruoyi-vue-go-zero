package system

import (
	"net/mail"
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"ruoyi-vue-go-zero/app/admin/api/internal/model"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func (s *Session) userList(r *types.ListFilter, membership string) (*types.UserPage, error) {
	if membership != "" {
		if err := s.roleAllowed(r.RoleId, false); err != nil {
			return nil, err
		}
	}
	rows, total, err := s.models().UserRows(s.ctx, r, s.userScope(), membership)
	if err != nil {
		return nil, err
	}
	return &types.UserPage{StatusResponse: ok(), Rows: rows, Total: total}, nil
}

func (s *Session) ListUsers(r *types.ListFilter) (*types.UserPage, error) { return s.userList(r, "") }
func (s *Session) AllocatedUsers(r *types.ListFilter) (*types.UserPage, error) {
	return s.userList(r, "EXISTS")
}

func (s *Session) UnallocatedUsers(r *types.ListFilter) (*types.UserPage, error) {
	return s.userList(r, "NOT EXISTS")
}

func (s *Session) profile(id int64) (types.UserProfile, error) {
	return model.ReadUserProfile(s.ctx, s.db, id, false)
}

func (s *Session) Profile() (*types.ProfileResult, error) {
	u, err := s.profile(s.user)
	if err != nil {
		return nil, err
	}
	names := []string{}
	for _, role := range u.Roles {
		names = append(names, role.RoleName)
	}
	posts, err := many[string](s, "SELECT p.post_name FROM sys_post p JOIN sys_user_post up ON up.post_id=p.post_id WHERE up.user_id=? ORDER BY p.post_sort", s.user)
	return &types.ProfileResult{StatusResponse: ok(), Data: u, RoleGroup: strings.Join(names, ","), PostGroup: strings.Join(posts, ",")}, err
}

func (s *Session) userDetails(id int64) (*types.UserDetail, error) {
	result := &types.UserDetail{StatusResponse: ok(), RoleIds: []int64{}, PostIds: []int64{}}
	if id != 0 {
		if err := s.userAllowed(id, false); err != nil {
			return nil, err
		}
		u, err := s.profile(id)
		if err != nil {
			return nil, err
		}
		result.Data = &u
		for _, role := range u.Roles {
			result.RoleIds = append(result.RoleIds, role.RoleId)
		}
		result.PostIds, err = many[int64](s, "SELECT post_id FROM sys_user_post WHERE user_id=?", id)
		if err != nil {
			return nil, err
		}
	}
	var err error
	result.Roles, err = s.roleOptions(id)
	if err != nil {
		return nil, err
	}
	result.Posts, err = many[types.PostRecord](s, "SELECT "+postColumns+" FROM sys_post WHERE status='0' ORDER BY post_sort")
	return result, err
}

func (s *Session) UserOptions() (*types.UserDetail, error) { return s.userDetails(0) }
func (s *Session) UserDetail(r *types.RecordId) (*types.UserDetail, error) {
	return s.userDetails(r.Id)
}

func (s *Session) UserGrants(r *types.RecordId) (*types.UserGrants, error) {
	if err := s.userAllowed(r.Id, false); err != nil {
		return nil, err
	}
	u, err := s.profile(r.Id)
	if err != nil {
		return nil, err
	}
	roles, err := s.roleOptions(r.Id)
	return &types.UserGrants{StatusResponse: ok(), User: u, Roles: roles}, err
}

var phonePattern = regexp.MustCompile(`^[0-9]{11}$`)

func validateProfile(r types.ProfileInput) error {
	if strings.TrimSpace(r.NickName) == "" || utf8.RuneCountInString(r.NickName) > 30 {
		return apierror.Business("昵称不能为空且不能超过30个字符")
	}
	if r.Phonenumber != "" && !phonePattern.MatchString(r.Phonenumber) {
		return apierror.Business("手机号码格式不正确")
	}
	if r.Email != "" {
		addr, err := mail.ParseAddress(r.Email)
		if err != nil || addr.Address != r.Email || len(r.Email) > 50 {
			return apierror.Business("邮箱格式不正确")
		}
	}
	return nil
}
func (s *Session) checkProfile(id int64, r types.ProfileInput) error {
	if err := validateProfile(r); err != nil {
		return err
	}
	if r.Phonenumber != "" {
		n, err := s.models().ProfilePhoneConflicts(s.ctx, id, r.Phonenumber)
		if err != nil {
			return err
		}
		if n > 0 {
			return apierror.Business("手机号码已存在")
		}
	}
	if r.Email != "" {
		n, err := s.models().ProfileEmailConflicts(s.ctx, id, r.Email)
		if err != nil {
			return err
		}
		if n > 0 {
			return apierror.Business("邮箱已存在")
		}
	}
	return nil
}

func (s *Session) SaveProfile(r *types.ProfileInput) (*types.StatusResponse, error) {
	if err := s.checkProfile(s.user, *r); err != nil {
		return nil, err
	}
	return saved(s.exec("UPDATE sys_user SET nick_name=?,email=?,phonenumber=?,sex=?,update_by=?,update_time=NOW() WHERE user_id=?", r.NickName, r.Email, r.Phonenumber, r.Sex, s.name, s.user))
}

func hashPassword(p string) (string, error) {
	n := utf8.RuneCountInString(p)
	if n < 5 || n > 20 {
		return "", apierror.Business("密码长度应为5至20个字符")
	}
	h, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(h), err
}

func (s *Session) setPassword(id int64, p string) error {
	h, err := hashPassword(p)
	if err != nil {
		return err
	}
	return s.exec("UPDATE sys_user SET password=?,pwd_update_date=NOW(),update_by=?,update_time=NOW() WHERE user_id=?", h, s.name, id)
}

func (s *Session) ChangePassword(r *types.PasswordInput) (*types.StatusResponse, error) {
	h, err := one[string](s, "SELECT password FROM sys_user WHERE user_id=?", s.user)
	if err != nil {
		return nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(h), []byte(r.OldPassword)) != nil {
		return nil, apierror.Business("修改密码失败，旧密码错误")
	}
	if r.OldPassword == r.NewPassword {
		return nil, apierror.Business("新密码不能与旧密码相同")
	}
	return saved(s.setPassword(s.user, r.NewPassword))
}

func (s *Session) ResetPassword(r *types.UserPassword) (*types.StatusResponse, error) {
	if err := s.userAllowed(r.UserId, true); err != nil {
		return nil, err
	}
	return saved(s.setPassword(r.UserId, r.Password))
}

func validStatus(v string) bool { return v == "0" || v == "1" }
func (s *Session) ChangeUserStatus(r *types.UserStatus) (*types.StatusResponse, error) {
	if !validStatus(r.Status) {
		return nil, apierror.Business("状态无效")
	}
	if err := s.userAllowed(r.UserId, true); err != nil {
		return nil, err
	}
	return saved(s.exec("UPDATE sys_user SET status=?,update_by=?,update_time=NOW() WHERE user_id=?", r.Status, s.name, r.UserId))
}

func (s *Session) writeUser(r *types.UserInput, create bool) (*types.StatusResponse, error) {
	id := r.UserId
	if create {
		id = 0
	} else if err := s.userAllowed(id, true); err != nil {
		return nil, err
	}
	if r.DeptId != nil {
		if err := s.checkDept(*r.DeptId); err != nil {
			return nil, err
		}
	} else if !s.all {
		return nil, apierror.Forbidden
	}
	if !validStatus(r.Status) {
		return nil, apierror.Business("状态无效")
	}
	if err := s.checkProfile(id, types.ProfileInput{NickName: r.NickName, Email: r.Email, Phonenumber: r.Phonenumber}); err != nil {
		return nil, err
	}
	if err := s.checkAssignments(r.RoleIds); err != nil {
		return nil, err
	}
	for _, post := range r.PostIds {
		if err := s.require("SELECT COUNT(*) FROM sys_post WHERE status='0' AND post_id=?", post); err != nil {
			return nil, err
		}
	}
	if create {
		n := utf8.RuneCountInString(r.UserName)
		if n < 2 || n > 20 {
			return nil, apierror.Business("账号长度应为2至20个字符")
		}
		if err := s.rejectIf("SELECT COUNT(*) FROM sys_user WHERE user_name=? AND del_flag='0'", "用户名称已存在", r.UserName); err != nil {
			return nil, err
		}
		h, err := hashPassword(r.Password)
		if err != nil {
			return nil, err
		}
		res, err := s.db.ExecCtx(s.ctx, "INSERT INTO sys_user(user_name,nick_name,password,dept_id,email,phonenumber,sex,status,remark,create_by,create_time) VALUES(?,?,?,?,?,?,?,?,?,?,NOW())", r.UserName, r.NickName, h, r.DeptId, r.Email, r.Phonenumber, r.Sex, r.Status, r.Remark, s.name)
		if err != nil {
			return nil, err
		}
		id, err = res.LastInsertId()
		if err != nil {
			return nil, err
		}
	} else if err := s.exec("UPDATE sys_user SET nick_name=?,dept_id=?,email=?,phonenumber=?,sex=?,status=?,remark=?,update_by=?,update_time=NOW() WHERE user_id=?", r.NickName, r.DeptId, r.Email, r.Phonenumber, r.Sex, r.Status, r.Remark, s.name, id); err != nil {
		return nil, err
	}
	if err := s.replaceLinks("sys_user_role", "user_id", "role_id", id, r.RoleIds); err != nil {
		return nil, err
	}
	return saved(s.replaceLinks("sys_user_post", "user_id", "post_id", id, r.PostIds))
}

func (s *Session) CreateUser(r *types.UserInput) (*types.StatusResponse, error) {
	return s.writeUser(r, true)
}

func (s *Session) EditUser(r *types.UserInput) (*types.StatusResponse, error) {
	return s.writeUser(r, false)
}

func (s *Session) DeleteUsers(r *types.RecordIds) (*types.StatusResponse, error) {
	ids, err := parseIDs(r.Ids)
	if err != nil {
		return nil, err
	}
	for _, id := range ids {
		if id == s.user {
			return nil, apierror.Business("当前用户不能删除")
		}
		if err := s.userAllowed(id, true); err != nil {
			return nil, err
		}
		if err := s.exec("UPDATE sys_user SET del_flag='2',update_by=?,update_time=NOW() WHERE user_id=?", s.name, id); err != nil {
			return nil, err
		}
		for _, t := range []string{"sys_user_role", "sys_user_post"} {
			if err := s.exec("DELETE FROM "+t+" WHERE user_id=?", id); err != nil {
				return nil, err
			}
		}
	}
	return saved(nil)
}

func (s *Session) AssignRoles(r *types.RoleAssignment) (*types.StatusResponse, error) {
	if err := s.userAllowed(r.UserId, true); err != nil {
		return nil, err
	}
	ids, err := parseIDs(r.RoleIds)
	if err != nil {
		return nil, err
	}
	if err = s.checkAssignments(ids); err != nil {
		return nil, err
	}
	return saved(s.replaceLinks("sys_user_role", "user_id", "role_id", r.UserId, ids))
}
