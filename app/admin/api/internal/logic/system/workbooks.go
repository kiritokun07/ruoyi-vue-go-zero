package system

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/extrame/xls"
	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
	"github.com/xuri/excelize/v2"
)

func workbook(name string, headers []string, rows [][]string) (*File, error) {
	book := excelize.NewFile()
	defer book.Close()
	if e := book.SetSheetName("Sheet1", name); e != nil {
		return nil, e
	}
	writer, e := book.NewStreamWriter(name)
	if e != nil {
		return nil, e
	}
	writeRow := func(index int, values []string) error {
		cell, err := excelize.CoordinatesToCellName(1, index)
		if err != nil {
			return err
		}
		cells := make([]interface{}, len(values))
		for i, value := range values {
			cells[i] = value
		}
		return writer.SetRow(cell, cells)
	}
	if e = writeRow(1, headers); e != nil {
		return nil, e
	}
	for i, values := range rows {
		if e = writeRow(i+2, values); e != nil {
			return nil, e
		}
	}
	if e = writer.Flush(); e != nil {
		return nil, e
	}
	var buf bytes.Buffer
	if e = book.Write(&buf); e != nil {
		return nil, e
	}
	return &File{Name: name + ".xlsx", Type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", Bytes: buf.Bytes()}, nil
}

func exportRecords[T any](title, columns string, rows []T) (*File, error) {
	fields, labels := []string{}, []string{}
	for _, pair := range strings.Fields(columns) {
		f, l, _ := strings.Cut(pair, "=")
		fields = append(fields, f)
		labels = append(labels, l)
	}
	values := make([][]string, 0, len(rows))
	for _, r := range rows {
		v := reflect.ValueOf(r)
		row := []string{}
		for _, field := range fields {
			f := v.FieldByName(field)
			if f.Kind() == reflect.Pointer {
				f = f.Elem()
			}
			cell := ""
			if f.IsValid() {
				cell = fmt.Sprint(f.Interface())
			}
			row = append(row, cell)
		}
		values = append(values, row)
	}
	return workbook(title, labels, values)
}

var importHeaders = []string{"部门编号", "登录名称", "用户名称", "用户邮箱", "手机号码", "用户性别", "账号状态"}

func (s *Session) UserTemplate() (*File, error) {
	return workbook("用户导入模板", importHeaders, nil)
}

func (s *Session) ExportUsers(r *types.ListFilter) (*File, error) {
	rows := [][]string{}
	f := *r
	f.PageSize = 500
	for f.PageNum = 1; ; f.PageNum++ {
		p, e := s.ListUsers(&f)
		if e != nil {
			return nil, e
		}
		for _, u := range p.Rows {
			dept, leader, date := "", "", ""
			if u.Dept != nil {
				dept, leader = u.Dept.DeptName, u.Dept.Leader
			}
			if u.LoginDate != nil {
				date = *u.LoginDate
			}
			rows = append(rows, []string{strconv.FormatInt(u.UserId, 10), u.UserName, u.NickName, u.Email, u.Phonenumber, map[string]string{"0": "男", "1": "女", "2": "未知"}[u.Sex], map[string]string{"0": "正常", "1": "停用"}[u.Status], u.LoginIp, date, dept, leader})
		}
		if len(p.Rows) < f.PageSize {
			break
		}
	}
	return workbook("用户列表", strings.Fields("用户序号 登录名称 用户名称 用户邮箱 手机号码 用户性别 账号状态 最后登录IP 最后登录时间 部门名称 部门负责人"), rows)
}

func (s *Session) ExportRoles(r *types.ListFilter) (*File, error) {
	rows := []types.RoleRecord{}
	f := *r
	f.PageSize = 500
	for f.PageNum = 1; ; f.PageNum++ {
		p, e := s.ListRole(&f)
		if e != nil {
			return nil, e
		}
		rows = append(rows, p.Rows...)
		if len(p.Rows) < f.PageSize {
			break
		}
	}
	return exportRecords("角色列表", "RoleId=角色编号 RoleName=角色名称 RoleKey=权限字符 RoleSort=显示顺序 Status=角色状态", rows)
}

func (s *Session) ExportPosts(r *types.ListFilter) (*File, error) {
	rows := []types.PostRecord{}
	f := *r
	f.PageSize = 500
	for f.PageNum = 1; ; f.PageNum++ {
		p, e := s.ListPost(&f)
		if e != nil {
			return nil, e
		}
		rows = append(rows, p.Rows...)
		if len(p.Rows) < f.PageSize {
			break
		}
	}
	return exportRecords("岗位列表", "PostId=岗位编号 PostName=岗位名称 PostCode=岗位编码 PostSort=岗位排序 Status=岗位状态", rows)
}

func (s *Session) ImportUsers(input io.Reader, update bool) (*types.Result, error) {
	data, e := io.ReadAll(io.LimitReader(input, (16<<20)+1))
	if e != nil {
		return nil, e
	}
	if len(data) > 16<<20 {
		return nil, apierror.Business("导入文件过大")
	}
	rows, e := userWorkbookRows(data)
	if e != nil {
		return nil, e
	}
	if len(rows) < 2 {
		return nil, apierror.Business("没有可导入的用户")
	}
	columns := map[string]int{}
	for i, c := range rows[0] {
		columns[strings.TrimSpace(c)] = i
	}
	for _, h := range []string{"登录名称", "用户名称"} {
		if _, ok := columns[h]; !ok {
			return nil, apierror.Business("模板缺少列：" + h)
		}
	}
	pwd, e := one[string](s, "SELECT config_value FROM sys_config WHERE config_key='sys.user.initPassword' LIMIT 1")
	if e != nil {
		return nil, e
	}
	n := 0
	for i, r := range rows[1:] {
		get := func(k string) string {
			j, ok := columns[k]
			if !ok || j >= len(r) {
				return ""
			}
			return strings.TrimSpace(r[j])
		}
		if len(r) == 0 {
			continue
		}
		u := types.UserInput{UserName: get("登录名称"), NickName: get("用户名称"), Phonenumber: get("手机号码"), Email: get("用户邮箱"), Password: pwd, Status: "0", Sex: "2"}
		if d := get("部门编号"); d != "" {
			id, e := strconv.ParseInt(d, 10, 64)
			if e != nil {
				return nil, apierror.Business("部门编号无效")
			}
			u.DeptId = &id
		}
		switch get("用户性别") {
		case "男", "0":
			u.Sex = "0"
		case "女", "1":
			u.Sex = "1"
		}
		if v := get("账号状态"); v == "停用" || v == "1" {
			u.Status = "1"
		}
		ids, e := many[int64](s, "SELECT user_id FROM sys_user WHERE user_name=? AND del_flag='0'", u.UserName)
		if e != nil {
			return nil, e
		}
		if len(ids) == 0 {
			_, e = s.CreateUser(&u)
		} else if !update {
			e = apierror.Business("用户已存在")
		} else {
			id := ids[0]
			if e = s.userAllowed(id, true); e == nil && u.DeptId != nil {
				e = s.checkDept(*u.DeptId)
			}
			if e == nil {
				e = s.checkProfile(id, types.ProfileInput{NickName: u.NickName, Email: u.Email, Phonenumber: u.Phonenumber})
			}
			if e == nil {
				e = s.exec("UPDATE sys_user SET nick_name=?,email=?,phonenumber=?,sex=?,status=?,dept_id=COALESCE(?,dept_id),update_by=?,update_time=NOW() WHERE user_id=?", u.NickName, u.Email, u.Phonenumber, u.Sex, u.Status, u.DeptId, s.name, id)
			}
		}
		if e != nil {
			return nil, fmt.Errorf("第%d行: %w", i+2, e)
		}
		n++
	}
	return &types.Result{StatusResponse: types.StatusResponse{Code: 200, Msg: fmt.Sprintf("成功导入 %d 条用户记录", n)}}, nil
}

func (s *Session) ExportParameter(r *types.AdminFilter) (*File, error) {
	rows, _, e := s.ParameterRows(r, false)
	if e != nil {
		return nil, e
	}
	return exportRecords("参数配置", "ConfigId=参数主键 ConfigName=参数名称 ConfigKey=参数键名 ConfigValue=参数键值 ConfigType=系统内置 Remark=备注", rows)
}

func (s *Session) ExportDictType(r *types.AdminFilter) (*File, error) {
	rows, _, e := s.DictTypeRows(r, false)
	if e != nil {
		return nil, e
	}
	return exportRecords("字典类型", "DictId=字典编号 DictName=字典名称 DictType=字典类型 Status=状态 Remark=备注", rows)
}

func (s *Session) ExportDictEntry(r *types.AdminFilter) (*File, error) {
	rows, _, e := s.DictEntryRows(r, false)
	if e != nil {
		return nil, e
	}
	return exportRecords("字典数据", "DictCode=字典编码 DictSort=字典排序 DictLabel=字典标签 DictValue=字典键值 DictType=字典类型 Status=状态", rows)
}

func (s *Session) ExportSchedule(r *types.AdminFilter) (*File, error) {
	rows, _, e := s.ScheduleRows(r, false)
	if e != nil {
		return nil, e
	}
	return exportRecords("定时任务", "JobId=任务编号 JobName=任务名称 JobGroup=任务组名 InvokeTarget=调用目标字符串 CronExpression=执行表达式 Status=状态", rows)
}

func (s *Session) ExportJobLog(r *types.AdminFilter) (*File, error) {
	rows, _, e := s.JobLogRows(r, false)
	if e != nil {
		return nil, e
	}
	return exportRecords("调度日志", "JobLogId=日志编号 JobName=任务名称 JobGroup=任务组名 InvokeTarget=调用目标 JobMessage=日志信息 Status=执行状态 ExceptionInfo=异常信息 CreateTime=执行时间", rows)
}

func (s *Session) ExportLoginEvent(r *types.AdminFilter) (*File, error) {
	rows, _, e := s.LoginEventRows(r, false)
	if e != nil {
		return nil, e
	}
	return exportRecords("登录日志", "InfoId=访问编号 UserName=用户名称 Ipaddr=登录地址 LoginLocation=登录地点 Browser=浏览器 Os=操作系统 Status=登录状态 Msg=操作信息 LoginTime=登录时间", rows)
}

func (s *Session) ExportOperation(r *types.AdminFilter) (*File, error) {
	rows, _, e := s.OperationRows(r, false)
	if e != nil {
		return nil, e
	}
	return exportRecords("操作日志", "OperId=日志序号 Title=模块标题 BusinessType=业务类型 OperName=操作人员 OperIp=主机地址 OperUrl=请求地址 Status=操作状态 OperTime=操作时间 CostTime=消耗时间", rows)
}

func ImportUserFile(ctx context.Context, sc *svc.ServiceContext, file io.Reader, update bool) (*types.Result, error) {
	return withSession(ctx, sc, "system:user:import", true, func(s *Session) (*types.Result, error) { return s.ImportUsers(file, update) })
}

// Both workbook formats are normalized to first-sheet text. Legacy XLS readers may panic on malformed BIFF records.
func userWorkbookRows(data []byte) (rows [][]string, err error) {
	defer func() {
		if recover() != nil {
			rows = nil
			err = apierror.Business("工作簿格式损坏")
		}
	}()
	if bytes.HasPrefix(data, []byte{0xd0, 0xcf, 0x11, 0xe0, 0xa1, 0xb1, 0x1a, 0xe1}) {
		book, e := xls.OpenReader(bytes.NewReader(data), "utf-8")
		if e != nil || book.NumSheets() == 0 {
			return nil, apierror.Business("无效的 XLS 工作簿")
		}
		sheet := book.GetSheet(0)
		if sheet == nil || sheet.MaxRow > 10000 {
			return nil, apierror.Business("工作表无效或超过一万行")
		}
		for i := 0; i <= int(sheet.MaxRow); i++ {
			row := sheet.Row(i)
			cells := make([]string, min(row.LastCol(), 256))
			for j := range cells {
				cells[j] = row.Col(j)
			}
			rows = append(rows, cells)
		}
		return rows, nil
	}
	z, e := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if e != nil {
		return nil, apierror.Business("请选择 XLSX 格式的用户工作簿")
	}
	var size uint64
	for _, f := range z.File {
		if f.UncompressedSize64 > 64<<20 || size+f.UncompressedSize64 > 64<<20 {
			return nil, apierror.Business("工作簿解压后过大")
		}
		size += f.UncompressedSize64
	}
	book, e := excelize.OpenReader(bytes.NewReader(data), excelize.Options{
		UnzipSizeLimit: 64 << 20, UnzipXMLSizeLimit: 16 << 20,
	})
	if e != nil {
		return nil, apierror.Business("无效的 XLSX 工作簿")
	}
	defer book.Close()
	sheets := book.GetSheetList()
	if len(sheets) == 0 {
		return nil, apierror.Business("没有可导入的用户")
	}
	iter, e := book.Rows(sheets[0])
	if e != nil {
		return nil, apierror.Business("工作簿格式损坏")
	}
	defer iter.Close()
	for iter.Next() {
		if len(rows) >= 10001 {
			return nil, apierror.Business("工作表超过一万行")
		}
		cells, e := iter.Columns()
		if e != nil {
			return nil, apierror.Business("工作簿格式损坏")
		}
		rows = append(rows, cells[:min(len(cells), 256)])
	}
	if e = iter.Error(); e != nil {
		return nil, apierror.Business("工作簿格式损坏")
	}
	if len(rows) < 2 {
		return nil, apierror.Business("没有可导入的用户")
	}
	return rows, nil
}
