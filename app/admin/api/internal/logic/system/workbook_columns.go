package system

// Column order, labels and conversions follow the Java @Excel declarations.
type exportColumn struct {
	Field, Label string
	Numeric      bool
	Values       map[string]string
	Suffix       string
}

var exportRolesColumns = []exportColumn{
	{Field: "RoleId", Label: "角色序号", Numeric: true},
	{Field: "RoleName", Label: "角色名称"},
	{Field: "RoleKey", Label: "角色权限"},
	{Field: "RoleSort", Label: "角色排序"},
	{Field: "DataScope", Label: "数据范围", Values: map[string]string{"1": "所有数据权限", "2": "自定义数据权限", "3": "本部门数据权限", "4": "本部门及以下数据权限", "5": "仅本人数据权限"}},
	{Field: "Status", Label: "角色状态", Values: map[string]string{"0": "正常", "1": "停用"}},
}

var exportPostsColumns = []exportColumn{
	{Field: "PostId", Label: "岗位序号", Numeric: true},
	{Field: "PostCode", Label: "岗位编码"},
	{Field: "PostName", Label: "岗位名称"},
	{Field: "PostSort", Label: "岗位排序"},
	{Field: "Status", Label: "状态", Values: map[string]string{"0": "正常", "1": "停用"}},
}

var exportParameterColumns = []exportColumn{
	{Field: "ConfigId", Label: "参数主键", Numeric: true},
	{Field: "ConfigName", Label: "参数名称"},
	{Field: "ConfigKey", Label: "参数键名"},
	{Field: "ConfigValue", Label: "参数键值"},
	{Field: "ConfigType", Label: "系统内置", Values: map[string]string{"Y": "是", "N": "否"}},
}

var exportDictTypeColumns = []exportColumn{
	{Field: "DictId", Label: "字典主键", Numeric: true},
	{Field: "DictName", Label: "字典名称"},
	{Field: "DictType", Label: "字典类型"},
	{Field: "Status", Label: "状态", Values: map[string]string{"0": "正常", "1": "停用"}},
}

var exportDictEntryColumns = []exportColumn{
	{Field: "DictCode", Label: "字典编码", Numeric: true},
	{Field: "DictSort", Label: "字典排序", Numeric: true},
	{Field: "DictLabel", Label: "字典标签"},
	{Field: "DictValue", Label: "字典键值"},
	{Field: "DictType", Label: "字典类型"},
	{Field: "IsDefault", Label: "是否默认", Values: map[string]string{"Y": "是", "N": "否"}},
	{Field: "Status", Label: "状态", Values: map[string]string{"0": "正常", "1": "停用"}},
}

var exportScheduleColumns = []exportColumn{
	{Field: "JobId", Label: "任务序号", Numeric: true},
	{Field: "JobName", Label: "任务名称"},
	{Field: "JobGroup", Label: "任务组名"},
	{Field: "InvokeTarget", Label: "调用目标字符串"},
	{Field: "CronExpression", Label: "执行表达式 "},
	{Field: "MisfirePolicy", Label: "计划策略 ", Values: map[string]string{"0": "默认", "1": "立即触发执行", "2": "触发一次执行", "3": "不触发立即执行"}},
	{Field: "Concurrent", Label: "并发执行", Values: map[string]string{"0": "允许", "1": "禁止"}},
	{Field: "Status", Label: "任务状态", Values: map[string]string{"0": "正常", "1": "暂停"}},
}

var exportJobLogColumns = []exportColumn{
	{Field: "JobLogId", Label: "日志序号"},
	{Field: "JobName", Label: "任务名称"},
	{Field: "JobGroup", Label: "任务组名"},
	{Field: "InvokeTarget", Label: "调用目标字符串"},
	{Field: "JobMessage", Label: "日志信息"},
	{Field: "Status", Label: "执行状态", Values: map[string]string{"0": "正常", "1": "失败"}},
	{Field: "ExceptionInfo", Label: "异常信息"},
}

var exportLoginEventColumns = []exportColumn{
	{Field: "InfoId", Label: "序号", Numeric: true},
	{Field: "UserName", Label: "用户账号"},
	{Field: "Status", Label: "登录状态", Values: map[string]string{"0": "成功", "1": "失败"}},
	{Field: "Ipaddr", Label: "登录地址"},
	{Field: "LoginLocation", Label: "登录地点"},
	{Field: "Browser", Label: "浏览器"},
	{Field: "Os", Label: "操作系统"},
	{Field: "Msg", Label: "提示消息"},
	{Field: "LoginTime", Label: "访问时间"},
}

var exportOperationColumns = []exportColumn{
	{Field: "OperId", Label: "操作序号", Numeric: true},
	{Field: "Title", Label: "操作模块"},
	{Field: "BusinessType", Label: "业务类型", Values: map[string]string{"0": "其它", "1": "新增", "2": "修改", "3": "删除", "4": "授权", "5": "导出", "6": "导入", "7": "强退", "8": "生成代码", "9": "清空数据"}},
	{Field: "Method", Label: "请求方法"},
	{Field: "RequestMethod", Label: "请求方式"},
	{Field: "OperatorType", Label: "操作类别", Values: map[string]string{"0": "其它", "1": "后台用户", "2": "手机端用户"}},
	{Field: "OperName", Label: "操作人员"},
	{Field: "DeptName", Label: "部门名称"},
	{Field: "OperUrl", Label: "请求地址"},
	{Field: "OperIp", Label: "操作地址"},
	{Field: "OperLocation", Label: "操作地点"},
	{Field: "OperParam", Label: "请求参数"},
	{Field: "JsonResult", Label: "返回参数"},
	{Field: "Status", Label: "状态", Values: map[string]string{"0": "正常", "1": "异常"}},
	{Field: "ErrorMsg", Label: "错误消息"},
	{Field: "OperTime", Label: "操作时间"},
	{Field: "CostTime", Label: "消耗时间", Suffix: "毫秒"},
}
