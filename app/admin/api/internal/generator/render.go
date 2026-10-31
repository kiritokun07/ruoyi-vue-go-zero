// Package generator renders the upstream RuoYi Java, MyBatis and Vue templates in Go.
package generator

import (
	"bytes"
	"embed"
	"fmt"
	"path"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	vtl "github.com/iron-s/go-vtl"

	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

//go:embed vm
var assets embed.FS

type column struct {
	types.GenColumn
	Pk, Increment, Required, Insert, Edit, List, Query, SuperColumn, UsableColumn bool
	CapJavaField, ShortComment                                                    string
}

func (c column) ReadConverterExp() string {
	_, tail, ok := strings.Cut(c.ColumnComment, "（")
	if !ok {
		return ""
	}
	text, _, _ := strings.Cut(tail, "）")
	var parts []string
	for _, part := range strings.Fields(text) {
		_, n := utf8.DecodeRuneInString(part)
		parts = append(parts, part[:n]+"="+part[n:])
	}
	return strings.Join(parts, ",")
}

type table struct {
	types.GenTable
	Crud, Tree, Sub bool
	Columns         []column
	SubTable        *table
}

func (t table) IsSuperColumn(s string) bool {
	return slices.Contains([]string{"createBy", "createTime", "updateBy", "updateTime", "remark"}, s) || (t.Tree && slices.Contains([]string{"parentName", "parentId", "orderNum", "ancestors"}, s))
}
func cap(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
func uncap(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}
func makeTable(g types.GenTable) *table {
	t := &table{GenTable: g, Crud: g.TplCategory == "crud", Tree: g.TplCategory == "tree", Sub: g.TplCategory == "sub"}
	for _, v := range g.Columns {
		c := column{GenColumn: v, Pk: v.IsPk == "1", Increment: v.IsIncrement == "1", Required: v.IsRequired == "1", Insert: v.IsInsert == "1", Edit: v.IsEdit == "1", List: v.IsList == "1", Query: v.IsQuery == "1", CapJavaField: cap(v.JavaField)}
		c.ShortComment, _, _ = strings.Cut(v.ColumnComment, "（")
		if len(v.JavaField) > 2 && v.JavaField[1] >= 'A' && v.JavaField[1] <= 'Z' {
			c.CapJavaField = v.JavaField
		}
		c.SuperColumn = t.IsSuperColumn(v.JavaField) || slices.Contains([]string{"parentName", "parentId", "orderNum", "ancestors"}, v.JavaField)
		c.UsableColumn = slices.Contains([]string{"parentId", "orderNum", "remark"}, v.JavaField)
		t.Columns = append(t.Columns, c)
	}
	return t
}
func imports(t *table) []string {
	keys := map[string]bool{}
	if t.Sub {
		keys["java.util.List"] = true
	}
	for _, c := range t.Columns {
		switch c.JavaType {
		case "Date":
			keys["java.util.Date"] = true
			keys["com.fasterxml.jackson.annotation.JsonFormat"] = true
		case "BigDecimal":
			keys["java.math.BigDecimal"] = true
		}
	}
	out := []string{}
	for k := range keys {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

type Generated struct{ Template, Path, Text string }

var literal = regexp.MustCompile(`(?s)#\[\[(.*?)\]\]#`)

func Render(g types.GenTable, sub *types.GenTable) ([]Generated, error) {
	t := makeTable(g)
	if len(t.Columns) == 0 {
		return nil, fmt.Errorf("表 %s 没有可生成的列", g.TableName)
	}
	pk := t.Columns[0]
	dict := map[string]bool{}
	for _, c := range t.Columns {
		if c.Pk {
			pk = c
		}
		if c.DictType != "" {
			dict["'"+c.DictType+"'"] = true
		}
	}
	dicts := []string{}
	for k := range dict {
		dicts = append(dicts, k)
	}
	sort.Strings(dicts)
	parent := "3"
	if g.ParentMenuId != nil {
		parent = fmt.Sprint(g.ParentMenuId)
	}
	ctx := map[string]any{"tplCategory": g.TplCategory, "ClassName": g.ClassName, "className": uncap(g.ClassName), "functionName": g.FunctionName, "packageName": g.PackageName, "basePackage": path.Dir(strings.ReplaceAll(g.PackageName, ".", "/")), "moduleName": g.ModuleName, "businessName": g.BusinessName, "BusinessName": cap(g.BusinessName), "tableName": g.TableName, "author": g.FunctionAuthor, "datetime": time.Now().Format(time.DateOnly), "pkColumn": pk, "importList": imports(t), "permissionPrefix": g.ModuleName + ":" + g.BusinessName, "columns": t.Columns, "table": t, "dicts": strings.Join(dicts, ", "), "parentMenuId": parent, "expandColumn": 1, "tree_parent_code": g.TreeParentCode, "tree_name": g.TreeName}
	ctx["Entity"] = "BaseEntity"
	if t.Tree {
		ctx["Entity"] = "TreeEntity"
	}
	ctx["basePackage"] = strings.ReplaceAll(ctx["basePackage"].(string), "/", ".")
	for _, c := range t.Columns {
		for k, v := range map[string]string{"treeCode": g.TreeCode, "treeName": g.TreeName, "treeParentCode": g.TreeParentCode} {
			if v == c.ColumnName {
				ctx[k] = c.JavaField
			}
		}
	}
	if t.Sub {
		if sub == nil {
			return nil, fmt.Errorf("子表未导入")
		}
		st := makeTable(*sub)
		t.SubTable = st
		ctx["subTable"] = st
		ctx["subClassName"] = st.ClassName
		ctx["subclassName"] = uncap(st.ClassName)
		ctx["subTableName"] = st.TableName
		ctx["subTableFkName"] = g.SubTableFkName
		ctx["subImportList"] = imports(st)
		for _, c := range st.Columns {
			if c.ColumnName == g.SubTableFkName {
				ctx["subTableFkClassName"] = cap(c.JavaField)
				ctx["subTableFkclassName"] = c.JavaField
			}
		}
	}
	j := "main/java/" + strings.ReplaceAll(g.PackageName, ".", "/")
	c, b, m := g.ClassName, g.BusinessName, g.ModuleName
	paths := map[string]string{"vm/java/domain.java.vm": j + "/domain/" + c + ".java", "vm/java/mapper.java.vm": j + "/mapper/" + c + "Mapper.java", "vm/java/service.java.vm": j + "/service/I" + c + "Service.java", "vm/java/serviceImpl.java.vm": j + "/service/impl/" + c + "ServiceImpl.java", "vm/java/controller.java.vm": j + "/controller/" + c + "Controller.java", "vm/xml/mapper.xml.vm": "main/resources/mapper/" + m + "/" + c + "Mapper.xml", "vm/sql/sql.vm": b + "Menu.sql", "vm/js/api.js.vm": "vue/api/" + m + "/" + b + ".js"}
	vue := "vm/vue/"
	if g.TplWebType == "element-plus" {
		vue += "v3/"
	}
	view := "index.vue.vm"
	if t.Tree {
		view = "index-tree.vue.vm"
	}
	paths[vue+view] = "vue/views/" + m + "/" + b + "/index.vue"
	if t.Sub {
		paths["vm/java/sub-domain.java.vm"] = j + "/domain/" + sub.ClassName + ".java"
	}
	var out []Generated
	for key, name := range paths {
		raw, e := assets.ReadFile(key)
		if e != nil {
			return nil, e
		}
		src := string(raw)
		i := 0
		src = literal.ReplaceAllStringFunc(src, func(v string) string {
			k := fmt.Sprintf("literal%d", i)
			ctx[k] = v[3 : len(v)-3]
			i++
			return "${" + k + "}"
		})
		tmpl, e := vtl.Parse(src, "", "")
		if e != nil {
			return nil, fmt.Errorf("%s: %w", key, e)
		}
		var buf bytes.Buffer
		if e = tmpl.WithMaxIterations(100000).Execute(&buf, ctx); e != nil {
			return nil, fmt.Errorf("%s: %w", key, e)
		}
		out = append(out, Generated{key, name, buf.String()})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out, nil
}
