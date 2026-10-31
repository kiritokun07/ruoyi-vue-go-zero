# 源码启动指南

[返回 README](../README.zh-CN.md) | [English](development.md)

需要 **Go 1.27.3**、MySQL 8.4+、Redis 6.2+。本地代码生成工具为 goctl 1.10.2，核心依赖为 go-zero 1.10.3、Excelize v2，见 [go.mod](../go.mod)。

以下命令使用 Bash/zsh/Git Bash。先启动 MySQL 和 Redis，再下载项目并进入项目目录。

## 初始化数据库

按以下顺序在**空数据库**中执行一次，会导入开发账号和示例数据。不要对已有若依业务库重复执行整套脚本。

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

应用应使用有权访问 `ry-vue` 的数据库账号，将下方 `dev:password` 替换为实际账号和密码。使用代码生成器的建表功能时，该账号还需要相应的 DDL 权限。

## 启动后端

```bash
export RUOYI_MYSQL_DSN='dev:password@tcp(127.0.0.1:3306)/ry-vue?parseTime=true&loc=Local'
export RUOYI_REDIS_HOST='127.0.0.1:6379'
export RUOYI_REDIS_PASSWORD=''
export RUOYI_JWT_SECRET='replace-with-your-own-random-secret-at-least-32-characters'
go mod download
go run ./app/admin/api -f app/admin/api/etc/admin-api.yaml
```

请自行生成密钥（例如 `openssl rand -hex 32`），重启后保持不变。根目录的 `.env` 由 Docker Compose 读取；`go run` 需要像上面一样导出环境变量。PowerShell 可使用 `$env:NAME='value'` 设置。

打开 <http://127.0.0.1:8888/swagger-ui/index.html>。开发账号为 `admin / admin123`、`ry / admin123`。前端需要另外安装。本地配置默认监听 `127.0.0.1`，容器配置监听 `0.0.0.0`。

配置见 [admin-api.yaml](../app/admin/api/etc/admin-api.yaml) 和[默认值](../app/admin/api/internal/config/config.go)。保留 `var/` 文件和 Redis 数据；删除会话数据后，用户需要重新登录。

## 开发与验证

```text
app/admin/api/desc/       API 契约，入口 admin.api
app/admin/api/internal/  Handler、Logic、Model、ServiceContext
scripts/                 代码生成脚本
deploy/                  容器配置
sql/                     表结构和开发数据
```

```bash
go install github.com/zeromicro/go-zero/tools/goctl@v1.10.2
./scripts/generate-api.sh
go build ./...
go test ./...
go vet ./...
```

未设置 `RUOYI_SYSTEM_TEST_DSN` 时，MySQL 集成测试会跳过。请使用按上述脚本初始化的**独立、可丢弃测试库**，不要连接生产库：

```bash
RUOYI_SYSTEM_TEST_DSN='dev:password@tcp(127.0.0.1:3306)/ruoyi_test?parseTime=true&loc=Local' go test ./...
```

此示例假定将前面命令中的 `ry-vue` 替换成了 `ruoyi_test`，并完成建库与初始化。测试会临时写入测试数据。

开发约定见 [AGENTS.md](../AGENTS.md)。可执行 `git submodule update --init --recursive` 下载可选 AI 工具，它们不是运行依赖。

Excel 导入导出接口在 `.api` 分组中显式设置 `timeout: 1m`，普通接口使用 go-zero 默认值。大批量操作可能超时；已成功导入的记录会保留。
