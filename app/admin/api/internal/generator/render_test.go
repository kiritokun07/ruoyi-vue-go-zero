package generator

import (
	"strings"
	"testing"

	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func TestRender(t *testing.T) {
	g := types.GenTable{ClassName: "Widget", TableName: "biz_widget", ModuleName: "system", PackageName: "com.ruoyi.system", BusinessName: "widget", FunctionName: "样例", TplWebType: "element-plus", TreeCode: "id", TreeParentCode: "parent_id", TreeName: "title", SubTableName: "biz_child", SubTableFkName: "parent_id", Columns: []types.GenColumn{
		{ColumnName: "id", JavaField: "id", JavaType: "Long", IsPk: "1", IsIncrement: "1", ColumnComment: "编号"},
		{ColumnName: "title", JavaField: "title", JavaType: "String", HtmlType: "input", IsInsert: "1", IsEdit: "1", IsList: "1", IsQuery: "1", QueryType: "LIKE", ColumnComment: "类型（0正常 1停用）", DictType: "sys_normal_disable"},
		{ColumnName: "parent_id", JavaField: "parentId", JavaType: "Long", HtmlType: "input", ColumnComment: "父级"},
	}}
	sub := g
	sub.ClassName = "Child"
	sub.TableName = "biz_child"
	for _, web := range []string{"element-ui", "element-plus"} {
		g.TplWebType = web
		for _, typ := range []string{"crud", "tree", "sub"} {
			g.TplCategory = typ
			t.Run(web+"/"+typ, func(t *testing.T) {
				files, e := Render(g, &sub)
				if e != nil {
					t.Fatal(e)
				}
				if len(files) < 9 {
					t.Fatal("missing templates")
				}
				for _, f := range files {
					if strings.Contains(f.Text, "#foreach(") || f.Text == "" {
						t.Fatalf("unrendered %s", f.Path)
					}
				}
			})
		}
	}
}
