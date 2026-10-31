package system

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/xuri/excelize/v2"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

const exportBatchSize = 500

// XLSX rows spill to Excelize temporary files. The final ZIP is written directly
// to disk so neither all records nor the complete download accumulate in memory.
func streamCells(name string, headers []string, produce func(func([]any) error) error) (_ *File, err error) {
	book := excelize.NewFile()
	defer book.Close()
	if err = book.SetSheetName("Sheet1", name); err != nil {
		return nil, err
	}
	writer, err := book.NewStreamWriter(name)
	if err != nil {
		return nil, err
	}
	index := 0
	yield := func(values []any) error {
		index++
		cell, e := excelize.CoordinatesToCellName(1, index)
		if e != nil {
			return e
		}
		cells := make([]interface{}, len(values))
		for i, v := range values {
			cells[i] = v
		}
		return writer.SetRow(cell, cells)
	}
	headerCells := make([]any, len(headers))
	for i, h := range headers {
		headerCells[i] = h
	}
	if err = yield(headerCells); err != nil {
		return nil, err
	}
	if err = produce(yield); err != nil {
		return nil, err
	}
	if err = writer.Flush(); err != nil {
		return nil, err
	}
	tmp, err := os.CreateTemp("", "ruoyi-export-*.xlsx")
	if err != nil {
		return nil, err
	}
	keep := false
	defer func() {
		if !keep {
			tmp.Close()
			os.Remove(tmp.Name())
		}
	}()
	// Excelize's default Write buffers the ZIP. Its supported ZipWriter hook lets
	// us direct archive entries to the file instead. No encrypted ZIP is used.
	book.SetZipWriter(func(io.Writer) excelize.ZipWriter { return &workbookZIP{zip.NewWriter(tmp)} })
	if err = book.Write(io.Discard); err != nil {
		return nil, err
	}
	if _, err = tmp.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	keep = true
	return &File{Name: name + ".xlsx", Type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", Body: tmp, Cleanup: func() { os.Remove(tmp.Name()) }}, nil
}

func streamWorkbook(name string, headers []string, produce func(func([]string) error) error) (*File, error) {
	return streamCells(name, headers, func(yield func([]any) error) error {
		return produce(func(row []string) error {
			cells := make([]any, len(row))
			for i, v := range row {
				cells[i] = v
			}
			return yield(cells)
		})
	})
}

// Advertise ZIP64-capable local headers up front, since disk output cannot use
// Excelize's post-processing of an in-memory ZIP for entries larger than 4 GiB.
type workbookZIP struct{ *zip.Writer }

func (z *workbookZIP) Create(name string) (io.Writer, error) {
	return z.CreateHeader(&zip.FileHeader{Name: name, Method: zip.Deflate, ReaderVersion: 45})
}

func workbook(name string, headers []string, rows [][]string) (*File, error) {
	return streamWorkbook(name, headers, func(yield func([]string) error) error {
		for _, row := range rows {
			if err := yield(row); err != nil {
				return err
			}
		}
		return nil
	})
}
func exportBatches[T any](ctx context.Context, load func(int) ([]T, error), yield func(T) error) error {
	for page := 1; ; page++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		rows, err := load(page)
		if err != nil {
			return err
		}
		for _, row := range rows {
			if err = ctx.Err(); err != nil {
				return err
			}
			if err = yield(row); err != nil {
				return err
			}
		}
		if len(rows) < exportBatchSize {
			return nil
		}
	}
}
func exportRecords[T any](title string, columns []exportColumn, produce func(func(T) error) error) (*File, error) {
	labels := make([]string, len(columns))
	for i, c := range columns {
		labels[i] = c.Label
	}
	return streamCells(title, labels, func(yield func([]any) error) error {
		return produce(func(row T) error {
			value := reflect.ValueOf(row)
			cells := make([]any, len(columns))
			for i, c := range columns {
				field := value.FieldByName(c.Field)
				if !field.IsValid() {
					return fmt.Errorf("unknown export field %s", c.Field)
				}
				if field.Kind() == reflect.Pointer {
					field = field.Elem()
				}
				cells[i] = ""
				if !field.IsValid() {
					continue
				}
				text := fmt.Sprint(field.Interface())
				if c.Values != nil {
					cells[i] = c.Values[text]
				} else if c.Numeric {
					cells[i] = field.Interface()
				} else {
					cells[i] = text + c.Suffix
				}
			}
			return yield(cells)
		})
	})
}

var importHeaders = []string{"部门编号", "登录名称", "用户名称", "用户邮箱", "手机号码", "用户性别", "账号状态"}

func (s *Session) UserTemplate() (*File, error) {
	return workbook("用户导入模板", importHeaders, nil)
}

func (s *Session) ExportUsers(r *types.ListFilter) (*File, error) {
	f := *r
	return streamCells("用户列表", strings.Fields("用户序号 登录名称 用户名称 用户邮箱 手机号码 用户性别 账号状态 最后登录IP 最后登录时间 部门名称 部门负责人"), func(yield func([]any) error) error {
		return exportBatches(s.ctx, func(page int) ([]types.UserProfile, error) {
			f.PageNum = page
			f.PageSize = exportBatchSize
			return s.models().UserExportBatch(s.ctx, &f, s.userScope())
		}, func(u types.UserProfile) error {
			dept, leader, date := "", "", ""
			if u.Dept != nil {
				dept, leader = u.Dept.DeptName, u.Dept.Leader
			}
			if u.LoginDate != nil {
				date = *u.LoginDate
			}
			return yield([]any{u.UserId, u.UserName, u.NickName, u.Email, u.Phonenumber, map[string]string{"0": "男", "1": "女", "2": "未知"}[u.Sex], map[string]string{"0": "正常", "1": "停用"}[u.Status], u.LoginIp, date, dept, leader})
		})
	})
}
func (s *Session) ExportRoles(r *types.ListFilter) (*File, error) {
	f := *r
	var args []any
	if !s.all {
		args = []any{s.name}
	}
	return exportRecords("角色列表", exportRolesColumns, func(yield func(types.RoleRecord) error) error {
		return exportBatches(s.ctx, func(page int) ([]types.RoleRecord, error) {
			f.PageNum = page
			f.PageSize = exportBatchSize
			rows, _, err := s.models().RoleRows(s.ctx, &f, false, s.roleScope(), args)
			return rows, err
		}, yield)
	})
}
func (s *Session) ExportPosts(r *types.ListFilter) (*File, error) {
	f := *r
	return exportRecords("岗位列表", exportPostsColumns, func(yield func(types.PostRecord) error) error {
		return exportBatches(s.ctx, func(page int) ([]types.PostRecord, error) {
			f.PageNum = page
			f.PageSize = exportBatchSize
			rows, _, err := s.models().PostRows(s.ctx, &f, false)
			return rows, err
		}, yield)
	})
}

func (s *Session) ExportParameter(r *types.AdminFilter) (*File, error) {
	f := *r
	return exportRecords("参数配置", exportParameterColumns, func(yield func(types.ParameterEntry) error) error {
		return exportBatches(s.ctx, func(page int) ([]types.ParameterEntry, error) {
			f.PageNum = page
			f.PageSize = exportBatchSize
			return s.models().ParameterExportBatch(s.ctx, &f)
		}, yield)
	})
}

func (s *Session) ExportDictType(r *types.AdminFilter) (*File, error) {
	f := *r
	return exportRecords("字典类型", exportDictTypeColumns, func(yield func(types.DictTypeEntry) error) error {
		return exportBatches(s.ctx, func(page int) ([]types.DictTypeEntry, error) {
			f.PageNum = page
			f.PageSize = exportBatchSize
			return s.models().DictTypeExportBatch(s.ctx, &f)
		}, yield)
	})
}

func (s *Session) ExportDictEntry(r *types.AdminFilter) (*File, error) {
	f := *r
	return exportRecords("字典数据", exportDictEntryColumns, func(yield func(types.DictEntryEntry) error) error {
		return exportBatches(s.ctx, func(page int) ([]types.DictEntryEntry, error) {
			f.PageNum = page
			f.PageSize = exportBatchSize
			return s.models().DictEntryExportBatch(s.ctx, &f)
		}, yield)
	})
}

func (s *Session) ExportSchedule(r *types.AdminFilter) (*File, error) {
	f := *r
	return exportRecords("定时任务", exportScheduleColumns, func(yield func(types.ScheduleEntry) error) error {
		return exportBatches(s.ctx, func(page int) ([]types.ScheduleEntry, error) {
			f.PageNum = page
			f.PageSize = exportBatchSize
			return s.models().ScheduleExportBatch(s.ctx, &f)
		}, yield)
	})
}

func (s *Session) ExportJobLog(r *types.AdminFilter) (*File, error) {
	f := *r
	return exportRecords("调度日志", exportJobLogColumns, func(yield func(types.JobLogEntry) error) error {
		return exportBatches(s.ctx, func(page int) ([]types.JobLogEntry, error) {
			f.PageNum = page
			f.PageSize = exportBatchSize
			return s.models().JobLogExportBatch(s.ctx, &f)
		}, yield)
	})
}

func (s *Session) ExportLoginEvent(r *types.AdminFilter) (*File, error) {
	f := *r
	return exportRecords("登录日志", exportLoginEventColumns, func(yield func(types.LoginEventEntry) error) error {
		return exportBatches(s.ctx, func(page int) ([]types.LoginEventEntry, error) {
			f.PageNum = page
			f.PageSize = exportBatchSize
			return s.models().LoginEventExportBatch(s.ctx, &f)
		}, yield)
	})
}

func (s *Session) ExportOperation(r *types.AdminFilter) (*File, error) {
	f := *r
	return exportRecords("操作日志", exportOperationColumns, func(yield func(types.OperationEntry) error) error {
		return exportBatches(s.ctx, func(page int) ([]types.OperationEntry, error) {
			f.PageNum = page
			f.PageSize = exportBatchSize
			return s.models().OperationExportBatch(s.ctx, &f)
		}, yield)
	})
}
