package system

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/xwb1989/sqlparser"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

var tableComment = regexp.MustCompile("(?i)\\bcomment\\s*=?\\s*'([^'\\\\]*)'")
var ident = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
var pkgName = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*(\.[A-Za-z_][A-Za-z0-9_]*)*$`)

func capName(v string) string {
	if v == "" {
		return ""
	}
	return strings.ToUpper(v[:1]) + v[1:]
}

func uncap(v string) string {
	if v == "" {
		return ""
	}
	return strings.ToLower(v[:1]) + v[1:]
}

func className(v string) string {
	out := ""
	for _, part := range strings.Split(v, "_") {
		out += capName(part)
	}
	return out
}

func tableNames(v string) ([]string, error) {
	names := strings.Split(v, ",")
	for _, n := range names {
		if !ident.MatchString(n) || strings.HasPrefix(n, "gen_") || strings.HasPrefix(n, "go_job_") {
			return nil, apierror.Business("不支持的表名")
		}
	}
	return names, nil
}

const tableColumns = `table_id,table_name,table_comment,class_name,tpl_category,COALESCE(tpl_web_type,'') tpl_web_type,
 COALESCE(sub_table_name,'') sub_table_name,COALESCE(sub_table_fk_name,'') sub_table_fk_name,package_name,module_name,business_name,function_name,function_author,gen_type,gen_path,
 COALESCE(options,'{}') options,COALESCE(remark,'') remark,IFNULL(DATE_FORMAT(create_time,'%Y-%m-%d %H:%i:%s'),'') create_time,IFNULL(DATE_FORMAT(update_time,'%Y-%m-%d %H:%i:%s'),'') update_time`
const columnFields = `column_id,table_id,column_name,column_comment,column_type,java_type,java_field,is_pk,is_increment,is_required,is_insert,is_edit,is_list,is_query,query_type,html_type,dict_type,sort`

func (s *Session) genPage(r *types.AdminFilter, physical bool) (*types.PageResult, error) {
	f := filter{}
	cols, table, order := tableColumns, "gen_table", "table_id DESC"
	if physical {
		cols = "table_name AS table_name,table_comment AS table_comment,IFNULL(DATE_FORMAT(create_time,'%Y-%m-%d %H:%i:%s'),'') create_time"
		table = "information_schema.tables"
		order = "table_name"
		f.add("table_schema=DATABASE() AND table_type='BASE TABLE' AND LEFT(table_name,4)<>'gen_' AND LEFT(table_name,7)<>'go_job_' AND table_name NOT IN (SELECT table_name FROM gen_table)")
	}
	f.text("table_name", r.TableName)
	f.text("table_comment", r.TableComment)
	f.dates("create_time", &types.ListFilter{BeginTime: r.BeginTime, EndTime: r.EndTime})
	from := " FROM " + table + f.where()
	n, e := one[int64](s, "SELECT COUNT(*)"+from, f.args...)
	if e != nil {
		return nil, e
	}
	rows, e := many[types.GenTable](s, "SELECT "+cols+from+" ORDER BY "+order+page(&types.ListFilter{PageNum: r.PageNum, PageSize: r.PageSize}), f.args...)
	return &types.PageResult{StatusResponse: ok(), Rows: rows, Total: n}, e
}

func (s *Session) GenList(r *types.AdminFilter) (*types.PageResult, error) {
	return s.genPage(r, false)
}

func (s *Session) GenDatabase(r *types.AdminFilter) (*types.PageResult, error) {
	return s.genPage(r, true)
}

func (s *Session) genColumns(id int64) ([]types.GenColumn, error) {
	return many[types.GenColumn](s, "SELECT "+columnFields+" FROM gen_table_column WHERE table_id=? ORDER BY sort,column_id", id)
}

func (s *Session) GenColumns(r *types.RecordId) (*types.PageResult, error) {
	rows, e := s.genColumns(r.Id)
	return &types.PageResult{StatusResponse: ok(), Rows: rows, Total: int64(len(rows))}, e
}

func (s *Session) genTable(id int64) (types.GenTable, error) {
	t, e := one[types.GenTable](s, "SELECT "+tableColumns+" FROM gen_table WHERE table_id=?", id)
	if e != nil {
		return t, e
	}
	t.Columns, e = s.genColumns(id)
	if e != nil {
		return t, e
	}
	if t.Options != "" {
		if e = json.Unmarshal([]byte(t.Options), &t.Params); e != nil {
			return t, e
		}
	}
	t.TreeCode, _ = t.Params["treeCode"].(string)
	t.TreeParentCode, _ = t.Params["treeParentCode"].(string)
	t.TreeName, _ = t.Params["treeName"].(string)
	t.ParentMenuId = t.Params["parentMenuId"]
	return t, nil
}

func (s *Session) GenInfo(r *types.RecordId) (*types.Result, error) {
	t, e := s.genTable(r.Id)
	if e != nil {
		return nil, e
	}
	tables, e := many[types.GenTable](s, "SELECT "+tableColumns+" FROM gen_table ORDER BY table_id")
	if e != nil {
		return nil, e
	}
	for i := range tables {
		tables[i].Columns, e = s.genColumns(tables[i].TableId)
		if e != nil {
			return nil, e
		}
	}
	return result(map[string]any{"info": t, "rows": t.Columns, "tables": tables}, nil)
}

func (s *Session) syncColumns(id int64, table string) error {
	cols, e := many[types.GenColumn](s, `SELECT column_name AS column_name,column_comment AS column_comment,column_type AS column_type,IF(column_key='PRI','1','0') is_pk,IF(extra LIKE '%auto_increment%','1','0') is_increment,IF(is_nullable='NO','1','0') is_required,ordinal_position sort FROM information_schema.columns WHERE table_schema=DATABASE() AND table_name=? ORDER BY ordinal_position`, table)
	if e != nil {
		return e
	}
	if len(cols) == 0 {
		return apierror.Business("数据表不存在")
	}
	prev, e := s.genColumns(id)
	if e != nil {
		return e
	}
	old := map[string]types.GenColumn{}
	for _, c := range prev {
		old[c.ColumnName] = c
	}
	for _, c := range cols {
		if v, ok := old[c.ColumnName]; ok {
			e = s.exec("UPDATE gen_table_column SET column_type=?,column_comment=?,is_pk=?,is_increment=? WHERE column_id=?", c.ColumnType, c.ColumnComment, c.IsPk, c.IsIncrement, v.ColumnId)
			if e != nil {
				return e
			}
			delete(old, c.ColumnName)
			continue
		}
		java, html := "String", "input"
		dt, _, _ := strings.Cut(c.ColumnType, "(")
		dt, _, _ = strings.Cut(dt, " ")
		if dt == "" || c.ColumnName == "" {
			return apierror.Business("数据库列定义为空")
		}
		switch dt {
		case "bigint":
			java = "Long"
		case "int", "smallint", "tinyint", "mediumint":
			java = "Integer"
		case "decimal", "numeric":
			java = "BigDecimal"
		case "float", "double":
			java = "Double"
		case "timestamp", "date", "datetime":
			java, html = "Date", "datetime"
		case "text", "mediumtext", "longtext":
			html = "textarea"
		}
		if strings.HasSuffix(c.ColumnName, "status") {
			html = "radio"
		}
		editable := "1"
		if c.IsPk == "1" || slices.Contains([]string{"create_time", "update_time", "create_by", "update_by"}, c.ColumnName) {
			editable = "0"
		}
		comment := c.ColumnComment
		if comment == "" {
			comment = c.ColumnName
		}
		e = s.exec("INSERT INTO gen_table_column(table_id,column_name,column_comment,column_type,java_type,java_field,is_pk,is_increment,is_required,is_insert,is_edit,is_list,is_query,query_type,html_type,dict_type,sort,create_by,create_time) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,'EQ',?,'',?,?,NOW())", id, c.ColumnName, comment, c.ColumnType, java, uncap(className(c.ColumnName)), c.IsPk, c.IsIncrement, c.IsRequired, editable, editable, editable, editable, html, c.Sort, s.name)
		if e != nil {
			return e
		}
	}
	for _, v := range old {
		if e = s.exec("DELETE FROM gen_table_column WHERE column_id=?", v.ColumnId); e != nil {
			return e
		}
	}
	return nil
}

func (s *Session) GenImport(r *types.TablesForm) (*types.StatusResponse, error) {
	names, e := tableNames(r.Tables)
	if e != nil {
		return nil, e
	}
	for _, name := range names {
		if e = s.rejectIf("SELECT COUNT(*) FROM gen_table WHERE table_name=?", "该表已导入", name); e != nil {
			return nil, e
		}
		comment, e := one[string](s, "SELECT table_comment FROM information_schema.tables WHERE table_schema=DATABASE() AND table_type='BASE TABLE' AND table_name=?", name)
		if e != nil {
			return nil, e
		}
		if comment == "" {
			comment = name
		}
		parts := strings.Split(name, "_")
		business := parts[len(parts)-1]
		res, e := s.db.ExecCtx(s.ctx, "INSERT INTO gen_table(table_name,table_comment,class_name,tpl_category,tpl_web_type,package_name,module_name,business_name,function_name,function_author,gen_type,gen_path,create_by,create_time) VALUES(?,?,?,'crud','element-plus','com.ruoyi.system','system',?,?,?,'0','/',?,NOW())", name, comment, className(name), business, comment, s.name, s.name)
		if e != nil {
			return nil, e
		}
		id, e := res.LastInsertId()
		if e != nil {
			return nil, e
		}
		if e = s.syncColumns(id, name); e != nil {
			return nil, e
		}
	}
	return saved(nil)
}

func (s *Session) GenSync(r *types.LookupKey) (*types.StatusResponse, error) {
	id, e := one[int64](s, "SELECT table_id FROM gen_table WHERE table_name=?", r.Key)
	if e != nil {
		return nil, e
	}
	return saved(s.syncColumns(id, r.Key))
}

func (s *Session) GenEdit(r *types.GenTable) (*types.StatusResponse, error) {
	if !ident.MatchString(r.ClassName) || !ident.MatchString(r.BusinessName) || !ident.MatchString(r.ModuleName) || !pkgName.MatchString(r.PackageName) || !slices.Contains([]string{"crud", "tree", "sub"}, r.TplCategory) {
		return nil, apierror.Business("生成配置无效")
	}
	t, e := s.genTable(r.TableId)
	if e != nil {
		return nil, e
	}
	available := map[string]bool{}
	ids := map[int64]bool{}
	for _, c := range t.Columns {
		available[c.ColumnName] = true
		ids[c.ColumnId] = true
	}
	if r.TplCategory == "tree" {
		for _, k := range []string{"treeCode", "treeParentCode", "treeName"} {
			v, _ := r.Params[k].(string)
			if !available[v] {
				return nil, apierror.Business("树表字段无效")
			}
		}
	}
	if r.TplCategory == "sub" {
		if r.SubTableName == t.TableName {
			return nil, apierror.Business("子表不能等于主表")
		}
		if e := s.require("SELECT COUNT(*) FROM gen_table_column c JOIN gen_table t ON t.table_id=c.table_id WHERE t.table_name=? AND c.column_name=?", r.SubTableName, r.SubTableFkName); e != nil {
			return nil, e
		}
	}
	opts, e := json.Marshal(r.Params)
	if e != nil {
		return nil, e
	}
	for i, c := range r.Columns {
		if !ids[c.ColumnId] || !ident.MatchString(c.JavaField) {
			return nil, apierror.Business("列配置无效")
		}
		delete(ids, c.ColumnId)
		e = s.exec("UPDATE gen_table_column SET column_comment=?,java_type=?,java_field=?,is_required=?,is_insert=?,is_edit=?,is_list=?,is_query=?,query_type=?,html_type=?,dict_type=?,sort=? WHERE column_id=? AND table_id=?", c.ColumnComment, c.JavaType, c.JavaField, c.IsRequired, c.IsInsert, c.IsEdit, c.IsList, c.IsQuery, c.QueryType, c.HtmlType, c.DictType, i+1, c.ColumnId, r.TableId)
		if e != nil {
			return nil, e
		}
	}
	return saved(s.exec("UPDATE gen_table SET table_comment=?,class_name=?,tpl_category=?,tpl_web_type=?,package_name=?,module_name=?,business_name=?,function_name=?,function_author=?,gen_type=?,gen_path=?,sub_table_name=?,sub_table_fk_name=?,options=?,remark=?,update_by=?,update_time=NOW() WHERE table_id=?", r.TableComment, r.ClassName, r.TplCategory, r.TplWebType, r.PackageName, r.ModuleName, r.BusinessName, r.FunctionName, r.FunctionAuthor, r.GenType, r.GenPath, r.SubTableName, r.SubTableFkName, string(opts), r.Remark, s.name, r.TableId))
}

func (s *Session) GenDelete(r *types.RecordIds) (*types.StatusResponse, error) {
	ids, e := parseIDs(r.Ids)
	if e != nil {
		return nil, e
	}
	if e = s.exec("DELETE FROM gen_table_column WHERE table_id IN (" + numbers(ids) + ")"); e != nil {
		return nil, e
	}
	return saved(s.exec("DELETE FROM gen_table WHERE table_id IN (" + numbers(ids) + ")"))
}

// Reconstruct ordinary InnoDB CREATE TABLE statements from their AST. MySQL table options
// such as DATA DIRECTORY and UNION must never be forwarded from the request.
func checkedDDL(input string) (string, string, error) {
	stmt, e := sqlparser.ParseStrictDDL(input)
	if e != nil {
		return "", "", apierror.Business("建表 SQL 格式不正确")
	}
	d, ok := stmt.(*sqlparser.DDL)
	if !ok || d.Action != sqlparser.CreateStr || d.TableSpec == nil || len(d.TableSpec.Columns) == 0 || !d.NewName.Qualifier.IsEmpty() {
		return "", "", apierror.Business("仅支持当前库的 CREATE TABLE")
	}
	name := d.NewName.Name.String()
	if _, e := tableNames(name); e != nil {
		return "", "", e
	}
	comment := tableComment.FindStringSubmatch(d.TableSpec.Options)
	d.TableSpec.Options = " ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"
	if len(comment) == 2 {
		d.TableSpec.Options += " COMMENT=" + sqlparser.String(sqlparser.NewStrVal([]byte(comment[1])))
	}

	return name, sqlparser.String(d), nil
}

func (s *Session) GenCreate(r *types.DDLForm) (*types.StatusResponse, error) {
	if s.user != 1 {
		return nil, apierror.Forbidden
	}
	pieces, e := sqlparser.SplitStatementToPieces(r.Sql)
	if e != nil {
		return nil, e
	}
	type ddl struct{ name, sql string }
	var plans []ddl
	for _, v := range pieces {
		if strings.TrimSpace(v) == "" {
			continue
		}
		n, q, e := checkedDDL(v)
		if e != nil {
			return nil, e
		}
		plans = append(plans, ddl{n, q})
	}
	if len(plans) == 0 {
		return nil, apierror.Business("请填写建表 SQL")
	}
	for _, plan := range plans {
		if e := s.exec(plan.sql); e != nil {
			return nil, e
		}
		e = s.services.DB.TransactCtx(s.ctx, func(ctx context.Context, tx sqlx.Session) error {
			copy := *s
			copy.db = tx
			_, e := copy.GenImport(&types.TablesForm{Tables: plan.name})
			return e
		})
		if e != nil {
			return nil, fmt.Errorf("表 %s 已创建，导入元数据失败: %w", plan.name, e)
		}
	}
	return saved(nil)
}
