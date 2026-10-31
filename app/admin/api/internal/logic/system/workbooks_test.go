package system

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/xuri/excelize/v2"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func TestWorkbookPreservesText(t *testing.T) {
	values := [][]string{{"00123", "中文用户", "=1+1", "13800138000", "", "末列"}}
	file, err := workbook("用户", []string{"编号", "姓名", "备注", "手机", "空列", "末列"}, values)
	if err != nil {
		t.Fatal(err)
	}
	file.Bytes = readWorkbookTest(t, file)
	book, err := excelize.OpenReader(bytes.NewReader(file.Bytes))
	if err != nil {
		t.Fatal(err)
	}
	defer book.Close()
	if file.Name != "用户.xlsx" || len(book.GetSheetList()) != 1 {
		t.Fatal("unexpected workbook metadata")
	}
	formula, err := book.GetCellFormula("用户", "C2")
	if err != nil || formula != "" {
		t.Fatalf("text became formula: %q, %v", formula, err)
	}
	rows, err := userWorkbookRows(file.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || !reflect.DeepEqual(rows[1], values[0]) {
		t.Fatalf("unexpected rows: %#v", rows)
	}
}

func TestUserWorkbookLimits(t *testing.T) {
	for _, count := range []int{0, 10000, 10001} {
		t.Run(strconv.Itoa(count), func(t *testing.T) {
			rows := make([][]string, count)
			for i := range rows {
				rows[i] = []string{"用户"}
			}
			file, err := workbook("用户", importHeaders, rows)
			if err != nil {
				t.Fatal(err)
			}
			got, err := userWorkbookRows(readWorkbookTest(t, file))
			if count == 10000 {
				if err != nil || len(got) != 10001 {
					t.Fatalf("valid limit: %d rows, %v", len(got), err)
				}
			} else if err == nil {
				t.Fatal("expected empty/oversized workbook rejection")
			}
		})
	}
	if _, err := userWorkbookRows([]byte("not a workbook")); err == nil {
		t.Fatal("expected invalid workbook rejection")
	}
}

func TestUserWorkbookFirstSheetAndSparseCells(t *testing.T) {
	book := excelize.NewFile()
	defer book.Close()
	for cell, value := range map[string]string{"A1": "登录名称", "A2": "first", "C2": "third", "A4": "fourth"} {
		if err := book.SetCellStr("Sheet1", cell, value); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := book.NewSheet("第二张"); err != nil {
		t.Fatal(err)
	}
	buf, err := book.WriteToBuffer()
	if err != nil {
		t.Fatal(err)
	}
	rows, err := userWorkbookRows(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 4 || !reflect.DeepEqual(rows[1], []string{"first", "", "third"}) || len(rows[2]) != 0 {
		t.Fatalf("unexpected sparse rows: %#v", rows)
	}
}

func readWorkbookTest(t *testing.T, file *File) []byte {
	t.Helper()
	defer file.Close()
	data, err := io.ReadAll(file.Body)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
func userWorkbookRows(data []byte) (rows [][]string, err error) {
	file, err := os.CreateTemp("", "workbook-test-*")
	if err != nil {
		return nil, err
	}
	defer os.Remove(file.Name())
	if _, err = file.Write(data); err != nil {
		file.Close()
		return nil, err
	}
	file.Close()
	err = walkUserWorkbook(file.Name(), func(_ int, cells []string) error { rows = append(rows, cells); return nil })
	return rows, err
}

func TestStreamingExportBatchesAndCleanup(t *testing.T) {
	calls := 0
	file, err := streamWorkbook("批次", []string{"值"}, func(yield func([]string) error) error {
		return exportBatches(context.Background(), func(page int) ([]int, error) {
			calls++
			n := exportBatchSize
			if page == 3 {
				n = 1
			}
			rows := make([]int, n)
			for i := range rows {
				rows[i] = (page-1)*exportBatchSize + i
			}
			return rows, nil
		}, func(v int) error { return yield([]string{strconv.Itoa(v)}) })
	})
	if err != nil {
		t.Fatal(err)
	}
	path := file.Body.(*os.File).Name()
	recorder := httptest.NewRecorder()
	file.Send(recorder)
	if _, err = os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("temporary export was not removed")
	}
	rows, err := userWorkbookRows(recorder.Body.Bytes())
	if err != nil || len(rows) != 1002 || calls != 3 || rows[1001][0] != "1000" {
		t.Fatalf("rows=%d calls=%d error=%v", len(rows), calls, err)
	}
}

func TestPrepareUsersBeforeTransaction(t *testing.T) {
	file, err := workbook("导入", importHeaders, [][]string{{"103", "stream_user", "测试", "", "", "男", "正常"}})
	if err != nil {
		t.Fatal(err)
	}
	prepared, err := prepareUsers(context.Background(), bytes.NewReader(readWorkbookTest(t, file)))
	if err != nil {
		t.Fatal(err)
	}
	path := prepared.Name()
	prepared.Close()
	if _, err = os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("temporary import was not removed")
	}
	invalid, err := workbook("导入", importHeaders, [][]string{{"invalid", "user", "测试"}})
	if err != nil {
		t.Fatal(err)
	}
	prepared, err = prepareUsers(context.Background(), bytes.NewReader(readWorkbookTest(t, invalid)))
	if err != nil {
		t.Fatal(err)
	}
	defer prepared.Close()
	var record importedUser
	if err = json.NewDecoder(prepared).Decode(&record); err != nil || record.Error != "部门编号无效" {
		t.Fatalf("missing row error: %+v %v", record, err)
	}
}

func TestExportCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := exportBatches(ctx, func(int) ([]int, error) { t.Fatal("query ran after cancellation"); return nil, nil }, func(int) error { return nil })
	if !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}

type failingResponse struct{ header http.Header }

func (w *failingResponse) Header() http.Header       { return w.header }
func (w *failingResponse) WriteHeader(int)           {}
func (w *failingResponse) Write([]byte) (int, error) { return 0, io.ErrClosedPipe }
func TestExportCleanupOnDisconnect(t *testing.T) {
	file, err := workbook("测试", []string{"值"}, [][]string{{"内容"}})
	if err != nil {
		t.Fatal(err)
	}
	path := file.Body.(*os.File).Name()
	file.Send(&failingResponse{header: http.Header{}})
	if _, err = os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("export leaked after disconnect")
	}
}
func TestFailedPreparationCleansTemporaryFiles(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TMPDIR", dir)
	if _, err := prepareUsers(context.Background(), strings.NewReader("invalid")); err == nil {
		t.Fatal("expected format error")
	}
	files, err := os.ReadDir(dir)
	if err != nil || len(files) != 0 {
		t.Fatalf("temporary files left: %v %v", files, err)
	}
}

func TestDictionaryValidationHasNoPersistenceSideEffects(t *testing.T) {
	// No DB is provided: validation must not perform the old dictionary UPDATE.
	session := &Session{}
	if err := session.ValidateDictType(&types.DictTypeEntry{DictId: 1, DictName: "测试", DictType: "test_type", Status: "0"}, false); err != nil {
		t.Fatal(err)
	}
}

func TestRoleExportMatchesJavaColumns(t *testing.T) {
	file, err := exportRecords("角色", exportRolesColumns, func(yield func(types.RoleRecord) error) error {
		return yield(types.RoleRecord{RoleId: 1, RoleName: "管理员", RoleKey: "admin", DataScope: "1", Status: "0"})
	})
	if err != nil {
		t.Fatal(err)
	}
	book, err := excelize.OpenReader(bytes.NewReader(readWorkbookTest(t, file)))
	if err != nil {
		t.Fatal(err)
	}
	defer book.Close()
	rows, err := book.GetRows("角色")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(rows[0], []string{"角色序号", "角色名称", "角色权限", "角色排序", "数据范围", "角色状态"}) {
		t.Fatal(rows[0])
	}
	if rows[1][4] != "所有数据权限" || rows[1][5] != "正常" {
		t.Fatal(rows[1])
	}
	kind, err := book.GetCellType("角色", "A2")
	if err != nil || (kind != excelize.CellTypeNumber && kind != excelize.CellTypeUnset) {
		t.Fatalf("numeric ID: %v %v", kind, err)
	}
}
