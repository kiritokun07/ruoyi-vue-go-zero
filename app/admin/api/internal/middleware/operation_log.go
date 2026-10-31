package middleware

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type auditResponse struct {
	http.ResponseWriter
	status, code int
	msg          string
}

func (w *auditResponse) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
func (w *auditResponse) Write(p []byte) (int, error) {
	if strings.HasPrefix(w.Header().Get("Content-Type"), "application/json") {
		var v struct {
			Code int
			Msg  string
		}
		if json.Unmarshal(p, &v) == nil {
			w.code = v.Code
			w.msg = v.Msg
		}
	}
	return w.ResponseWriter.Write(p)
}
func (w *auditResponse) Unwrap() http.ResponseWriter { return w.ResponseWriter }
func OperationLog(db sqlx.SqlConn, secret string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			parts := strings.Split(strings.Trim(path, "/"), "/")
			if len(parts) < 2 || (parts[0] != "system" && parts[0] != "monitor" && parts[0] != "tool") || r.Method == "GET" && !strings.HasPrefix(path, "/tool/gen/genCode/") && !strings.Contains(path, "/unlock/") {
				next(w, r)
				return
			}
			claims := ParseBearer(r.Header.Get("Authorization"), []byte(secret))
			if claims == nil {
				next(w, r)
				return
			}
			start := time.Now()
			rec := &auditResponse{ResponseWriter: w}
			next(rec, r)
			business := map[string]int{"POST": 1, "PUT": 2, "DELETE": 3}[r.Method]
			if strings.HasSuffix(path, "/export") {
				business = 5
			} else if strings.Contains(path, "/import") {
				business = 6
			} else if strings.Contains(path, "/genCode") {
				business = 7
			} else if strings.HasSuffix(path, "/clean") {
				business = 9
			}
			status := 0
			if rec.status >= 400 || rec.code >= 400 {
				status = 1
			}
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			ctx, cancel := context.WithTimeout(context.WithoutCancel(r.Context()), 2*time.Second)
			defer cancel()
			titles := map[string]string{"user": "用户管理", "role": "角色管理", "dept": "部门管理", "post": "岗位管理", "menu": "菜单管理", "dict": "字典管理", "config": "参数管理", "notice": "通知公告", "job": "定时任务", "jobLog": "任务日志", "gen": "代码生成", "online": "在线用户", "operlog": "操作日志", "logininfor": "登录日志", "cache": "缓存管理"}
			title := titles[parts[1]]
			if title == "" {
				title = parts[1]
			}
			errorMessage := ""
			if status != 0 {
				errorMessage = rec.msg
			}
			result, _ := json.Marshal(map[string]any{"code": rec.code, "msg": rec.msg})
			params := map[string]string{}
			for _, key := range []string{"id", "roleId", "userId", "deptId", "jobId", "pageNum", "pageSize", "updateSupport"} {
				if v := r.URL.Query().Get(key); v != "" {
					params[key] = v
				}
			}
			safeParams, _ := json.Marshal(params)
			// Persist operation metadata only. Payloads may contain passwords, access tokens or uploaded files.
			_, e := db.ExecCtx(ctx, `INSERT INTO sys_oper_log(title,business_type,method,request_method,operator_type,oper_name,dept_name,oper_url,oper_ip,oper_param,json_result,status,error_msg,oper_time,cost_time)
   SELECT ?,?,?,?,1,u.user_name,IFNULL(d.dept_name,''),LEFT(?,255),?,LEFT(?,2000),LEFT(?,2000),?,LEFT(?,2000),NOW(),? FROM sys_user u LEFT JOIN sys_dept d ON d.dept_id=u.dept_id WHERE u.user_id=?`, title, business, r.Method+" "+parts[0]+"/"+parts[1], r.Method, path, ip, string(safeParams), string(result), status, errorMessage, time.Since(start).Milliseconds(), claims.UserID)
			if e != nil {
				logx.WithContext(ctx).Error(e)
			}
		}
	}
}
