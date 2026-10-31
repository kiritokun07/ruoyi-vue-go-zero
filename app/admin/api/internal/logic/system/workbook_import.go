package system

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"html"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/extrame/xls"
	"github.com/xuri/excelize/v2"
	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type importedUser struct {
	Line  int
	Error string
	User  types.UserInput
}
type preparedUsers struct{ *os.File }

func (p *preparedUsers) Close() { p.File.Close(); os.Remove(p.Name()) }

type contextReader struct {
	ctx    context.Context
	reader io.Reader
}

func (r contextReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.reader.Read(p)
}

// Parse before opening any row transactions. Normalized records and row errors
// are spooled to disk, preserving row numbers without keeping the sheet in RAM.
func prepareUsers(ctx context.Context, input io.Reader) (_ *preparedUsers, err error) {
	source, err := os.CreateTemp("", "ruoyi-import-*")
	if err != nil {
		return nil, err
	}
	defer func() { source.Close(); os.Remove(source.Name()) }()
	n, err := io.Copy(source, io.LimitReader(contextReader{ctx, input}, (16<<20)+1))
	if err != nil {
		return nil, err
	}
	if n > 16<<20 {
		return nil, apierror.Business("导入文件过大")
	}
	if err = source.Close(); err != nil {
		return nil, err
	}
	spool, err := os.CreateTemp("", "ruoyi-users-*.json")
	if err != nil {
		return nil, err
	}
	prepared := &preparedUsers{spool}
	keep := false
	defer func() {
		if !keep {
			prepared.Close()
		}
	}()
	encoder := json.NewEncoder(spool)
	columns := map[string]int{}
	err = walkUserWorkbook(source.Name(), func(line int, row []string) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		if line == 1 {
			for i, v := range row {
				columns[strings.TrimSpace(v)] = i
			}
			for _, h := range []string{"登录名称", "用户名称"} {
				if _, ok := columns[h]; !ok {
					return apierror.Business("模板缺少列：" + h)
				}
			}
			return nil
		}
		if len(row) == 0 {
			return nil
		}
		get := func(key string) string {
			index, ok := columns[key]
			if !ok || index >= len(row) {
				return ""
			}
			return strings.TrimSpace(row[index])
		}
		user := types.UserInput{UserName: get("登录名称"), NickName: get("用户名称"), Phonenumber: get("手机号码"), Email: get("用户邮箱"), Status: "0", Sex: "2"}
		if value := get("部门编号"); value != "" {
			id, e := strconv.ParseInt(value, 10, 64)
			if e != nil {
				return encoder.Encode(importedUser{Line: line, User: user, Error: "部门编号无效"})
			}
			user.DeptId = &id
		}
		switch get("用户性别") {
		case "男", "0":
			user.Sex = "0"
		case "女", "1":
			user.Sex = "1"
		}
		if v := get("账号状态"); v == "停用" || v == "1" {
			user.Status = "1"
		}
		if err := validateProfile(types.ProfileInput{NickName: user.NickName, Email: user.Email, Phonenumber: user.Phonenumber}); err != nil {
			return encoder.Encode(importedUser{Line: line, User: user, Error: err.Error()})
		}
		return encoder.Encode(importedUser{Line: line, User: user})
	})
	if err != nil {
		return nil, err
	}
	if _, err = spool.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	keep = true
	return prepared, nil
}

func ImportUserFile(ctx context.Context, sc *svc.ServiceContext, file io.Reader, update bool) (*types.Result, error) {
	// Reject unauthorized users before doing expensive parsing; recheck inside
	// each row transaction so permissions cannot become stale during parsing.
	if _, err := Read(ctx, sc, "system:user:import", func(s *Session) (*types.StatusResponse, error) { r := ok(); return &r, nil }); err != nil {
		return nil, err
	}
	prepared, err := prepareUsers(ctx, file)
	if err != nil {
		return nil, err
	}
	defer prepared.Close()
	password, err := sc.Admin.InitialPassword(ctx)
	if err != nil {
		return nil, err
	}
	return importUsers(ctx, prepared, func(record importedUser) error {
		_, err := withUserWrite(ctx, sc, record.User.UserName, func(ctx context.Context) (*types.StatusResponse, error) {
			return withSession(ctx, sc, "system:user:import", true, func(s *Session) (*types.StatusResponse, error) { return nil, s.importUser(record, password, update) })
		})
		return err
	})
}

// Each row has its own transaction, matching Java's partial-success behavior.
func importUsers(ctx context.Context, input io.Reader, apply func(importedUser) error) (*types.Result, error) {
	decoder := json.NewDecoder(input)
	success, failure := 0, 0
	var successes, failures strings.Builder
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		var record importedUser
		err := decoder.Decode(&record)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		if record.Error != "" {
			err = apierror.Business(record.Error)
		} else {
			err = apply(record)
		}
		name := html.EscapeString(record.User.UserName)
		if err != nil {
			failure++
			message := "系统异常，请稍后重试"
			var business apierror.Business
			if errors.As(err, &business) {
				message = business.Error()
			} else if errors.Is(err, apierror.Forbidden) {
				message = apierror.Forbidden.Error()
			} else {
				logx.WithContext(ctx).Error(err)
			}
			fmt.Fprintf(&failures, "<br/>%d、账号 %s 第%d行导入失败：%s", failure, name, record.Line, html.EscapeString(message))
		} else {
			success++
			fmt.Fprintf(&successes, "<br/>%d、账号 %s 导入成功", success, name)
		}
	}
	if failure > 0 {
		return nil, apierror.Business(fmt.Sprintf("很抱歉，导入失败！共 %d 条数据格式不正确，错误如下：%s", failure, failures.String()))
	}
	if success == 0 {
		return nil, apierror.Business("导入用户数据不能为空！")
	}
	return &types.Result{StatusResponse: types.StatusResponse{Code: 200, Msg: fmt.Sprintf("恭喜您，数据已全部导入成功！共 %d 条，数据如下：%s", success, successes.String())}}, nil
}

func (s *Session) importPreparedUsers(input io.Reader, update bool) (*types.Result, error) {
	password, err := s.models().InitialPassword(s.ctx)
	if err != nil {
		return nil, err
	}
	return importUsers(s.ctx, input, func(record importedUser) error { return s.importUser(record, password, update) })
}

func (s *Session) importUser(record importedUser, password string, update bool) error {
	user := record.User
	user.Password = password
	ids, err := s.models().UserIDsByName(s.ctx, user.UserName)
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		_, err = s.CreateUser(&user)
		return err
	}
	if !update {
		return apierror.Business("用户已存在")
	}
	id := ids[0]
	if err = s.userAllowed(id, true); err != nil {
		return err
	}
	if user.DeptId != nil {
		if err = s.checkDept(*user.DeptId); err != nil {
			return err
		}
	}
	if err = s.checkProfile(id, types.ProfileInput{NickName: user.NickName, Email: user.Email, Phonenumber: user.Phonenumber}); err != nil {
		return err
	}
	return s.models().UpdateImportedUser(s.ctx, id, &user, s.name)
}

// Iterate the first worksheet. XLS compatibility remains on the legacy parser;
// XLSX uses Excelize's row iterator and bounded decompression.
func walkUserWorkbook(path string, yield func(int, []string) error) (err error) {
	defer func() {
		if recover() != nil {
			err = apierror.Business("工作簿格式损坏")
		}
	}()
	source, err := os.Open(path)
	if err != nil {
		return err
	}
	defer source.Close()
	magic := make([]byte, 8)
	_, _ = source.Read(magic)
	if _, err = source.Seek(0, io.SeekStart); err != nil {
		return err
	}
	if bytes.Equal(magic, []byte{0xd0, 0xcf, 0x11, 0xe0, 0xa1, 0xb1, 0x1a, 0xe1}) {
		book, e := xls.OpenReader(source, "utf-8")
		if e != nil || book.NumSheets() == 0 {
			return apierror.Business("无效的 XLS 工作簿")
		}
		sheet := book.GetSheet(0)
		if sheet == nil || sheet.MaxRow > 10000 {
			return apierror.Business("工作表无效或超过一万行")
		}
		if sheet.MaxRow < 1 {
			return apierror.Business("没有可导入的用户")
		}
		for i := 0; i <= int(sheet.MaxRow); i++ {
			row := sheet.Row(i)
			var cells []string
			if row != nil {
				cells = make([]string, min(row.LastCol(), 256))
				for j := range cells {
					cells[j] = row.Col(j)
				}
			}
			if err = yield(i+1, cells); err != nil {
				return err
			}
		}
		return nil
	}
	archive, e := zip.OpenReader(path)
	if e != nil {
		return apierror.Business("请选择 XLSX 格式的用户工作簿")
	}
	var size uint64
	for _, entry := range archive.File {
		if entry.UncompressedSize64 > 64<<20 || size+entry.UncompressedSize64 > 64<<20 {
			archive.Close()
			return apierror.Business("工作簿解压后过大")
		}
		size += entry.UncompressedSize64
	}
	if err = archive.Close(); err != nil {
		return err
	}
	book, e := excelize.OpenFile(path, excelize.Options{UnzipSizeLimit: 64 << 20, UnzipXMLSizeLimit: 16 << 20})
	if e != nil {
		return apierror.Business("无效的 XLSX 工作簿")
	}
	defer book.Close()
	sheets := book.GetSheetList()
	if len(sheets) == 0 {
		return apierror.Business("没有可导入的用户")
	}
	iter, e := book.Rows(sheets[0])
	if e != nil {
		return apierror.Business("工作簿格式损坏")
	}
	defer iter.Close()
	line := 0
	for iter.Next() {
		line++
		if line > 10001 {
			return apierror.Business("工作表超过一万行")
		}
		cells, e := iter.Columns()
		if e != nil {
			return apierror.Business("工作簿格式损坏")
		}
		if err = yield(line, cells[:min(len(cells), 256)]); err != nil {
			return err
		}
	}
	if iter.Error() != nil {
		return apierror.Business("工作簿格式损坏")
	}
	if line < 2 {
		return apierror.Business("没有可导入的用户")
	}
	return nil
}
