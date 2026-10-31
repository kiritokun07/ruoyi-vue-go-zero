package system

import (
	"bytes"
	"reflect"
	"strconv"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestWorkbookPreservesText(t *testing.T) {
	values := [][]string{{"00123", "中文用户", "=1+1", "13800138000", "", "末列"}}
	file, err := workbook("用户", []string{"编号", "姓名", "备注", "手机", "空列", "末列"}, values)
	if err != nil {
		t.Fatal(err)
	}
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
			got, err := userWorkbookRows(file.Bytes)
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
