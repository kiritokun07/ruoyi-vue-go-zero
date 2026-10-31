# RuoYi Vue Go Zero

基于 **go-zero** 的若依 Go 后端。

[English](README.md) | 简体中文

沿用 [yangzongzhuan](https://github.com/yangzongzhuan) 的原版项目 [RuoYi-Vue](https://github.com/yangzongzhuan/RuoYi-Vue) 的接口约定和 MySQL 数据模型，独立实现 Go 后端。管理界面可搭配兼容的 [RuoYi-Vue3 前端](https://github.com/yangzongzhuan/RuoYi-Vue3)。

## 功能

- 验证码、JWT 认证、角色权限与数据权限。
- 用户、角色、部门、岗位、菜单、字典、参数和公告管理。
- 在线会话、审计日志、定时任务、Redis 和主机监控。
- 文件上传下载、Excel 导入导出、Swagger 接口文档。
- 基于若依 Java、MyBatis、SQL 和 Vue 模板的代码生成。

采用单体 REST 服务，依赖 MySQL 和 Redis。定时任务当前支持内置 `ryTask` 示例目标，新增目标需要实现对应 Go 代码。

## 快速开始

需要 **Go 1.27.3**、MySQL 和 **Redis 6.2+**。MySQL 本地集成验证版本为 9.4.0。核心依赖为 go-zero 1.10.3 和 Excelize v2，完整版本见 [go.mod](go.mod)。

```bash
git clone https://github.com/kiritokun07/ruoyi-vue-go-zero.git
cd ruoyi-vue-go-zero
go mod download
```

在**空的开发数据库**中，按顺序执行以下脚本一次。不要在已有数据库上直接运行整套脚本。

```bash
mysql -u root -p -e 'CREATE DATABASE `ry-vue` CHARACTER SET utf8mb4'

mysql -u root -p ry-vue < sql/auth_schema.sql
mysql -u root -p ry-vue < sql/dev_seed.sql
mysql -u root -p ry-vue < sql/access_schema.sql
mysql -u root -p ry-vue < sql/access_dev_seed.sql
mysql -u root -p ry-vue < sql/system_schema.sql
mysql -u root -p ry-vue < sql/header_notice.sql
mysql -u root -p ry-vue < sql/system_dev_seed.sql
mysql -u root -p ry-vue < sql/operations_schema.sql
mysql -u root -p ry-vue < sql/scheduler_schema.sql
```

替换示例数据库凭据后启动：

```bash
export RUOYI_MYSQL_DSN='dev:password@tcp(127.0.0.1:3306)/ry-vue?parseTime=true&loc=Local'
export RUOYI_REDIS_HOST='127.0.0.1:6379'
export RUOYI_REDIS_PASSWORD=''
export RUOYI_JWT_SECRET="$(openssl rand -hex 32)"

go run ./app/admin/api -f app/admin/api/etc/admin-api.yaml
```

- 服务地址：`http://127.0.0.1:8888`
- Swagger：`/swagger-ui/index.html`
- 开发账号：`admin / admin123`、`ry / admin123`，对外提供服务前请修改密码。
- 前端代理：将 `/dev-api/*` 转发到后端，并移除 `/dev-api` 前缀。
- 部署时保持 JWT 密钥稳定，并保留 `var/` 下的运行时文件。

详细配置见[配置文件](app/admin/api/etc/admin-api.yaml)和[配置默认值](app/admin/api/internal/config/config.go)，接口文档通过 Swagger 查看。

## 开发

```text
app/admin/api/
├── desc/       # 接口契约
├── etc/        # 配置
├── internal/   # Handler, Logic, Model, ServiceContext
└── admin.go
scripts/        # 代码生成
sql/            # 表结构与开发种子数据
```

[desc/admin.api](app/admin/api/desc/admin.api) 是唯一生成入口。契约及 Handler/Logic 按业务模块分组，Model 负责持久化，ServiceContext 初始化依赖。

```bash
go install github.com/zeromicro/go-zero/tools/goctl@v1.10.2
./scripts/generate-api.sh
go build ./...
go test ./...
go vet ./...
```

MySQL 集成测试需连接已初始化的开发数据库：

```bash
RUOYI_SYSTEM_TEST_DSN="$RUOYI_MYSQL_DSN" go test ./...
```

欢迎提交 Issue 和 Pull Request。请保持若依接口行为兼容、补充相关测试，并同步更新中英文 README。开发指引见 [AGENTS.md](AGENTS.md)，可选 AI 子模块通过 `git submodule update --init --recursive` 获取。

## 致谢与许可证

感谢 [RuoYi-Vue](https://github.com/yangzongzhuan/RuoYi-Vue)、[go-zero](https://github.com/zeromicro/go-zero)、[go-zero-looklook](https://github.com/Mikaelemmmm/go-zero-looklook) 和 [Excelize](https://github.com/qax-os/excelize)。

本项目采用 [MIT 许可证](LICENSE)，若依代码生成模板保留[上游许可声明](app/admin/api/internal/generator/LICENSE.ruoyi)。
