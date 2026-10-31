# RuoYi Vue Go Zero

基于 **go-zero** 的若依后端，沿用原版 [RuoYi-Vue](https://github.com/yangzongzhuan/RuoYi-Vue) 的接口约定和 MySQL 数据模型，提供登录认证、系统管理、监控、Excel 导入导出与代码生成。

[English](README.md) | 简体中文

**本仓库只包含后端。** 默认仅启动后端 API，连接已有 MySQL 和 Redis；也可通过配置启用配套容器。需要管理界面时，请单独接入兼容的 [RuoYi-Vue3 前端](https://github.com/yangzongzhuan/RuoYi-Vue3)。

## Docker 安装（推荐）

无需安装 Go。默认连接已有 MySQL 和 Redis；没有数据库或缓存时，可按第 3 步启用配套容器。

### 1. 安装工具

- 安装 [Git](https://git-scm.com/downloads)。
- Windows/macOS：安装并启动 [Docker Desktop](https://docs.docker.com/desktop/)。Windows 使用 Linux 容器模式。
- Linux 服务器：安装 [Docker Engine](https://docs.docker.com/engine/install/) 和 [Compose 插件](https://docs.docker.com/compose/install/linux/)。

打开终端，执行 `docker version` 和 `docker compose version`。前者应同时显示客户端和服务端信息；Compose 需要 v2.20 或更新版本。

### 2. 下载项目

```bash
git clone https://github.com/kiritokun07/ruoyi-vue-go-zero.git
cd ruoyi-vue-go-zero
cp .env.example .env
```

以上命令也可以在 PowerShell 中执行。后续命令都在这个项目目录执行。AI 子模块是可选开发工具，安装和运行服务不需要下载它们。

### 3. 检查配置

用文本编辑器打开 `.env`。默认 `COMPOSE_PROFILES` 留空，**不会启动 MySQL、Redis 容器**。请先选择连接方式：

- **使用已有服务（默认）**：填写实际 `MYSQL_HOST`、`MYSQL_PORT`、`MYSQL_DATABASE`、`MYSQL_USER`、`MYSQL_PASSWORD`、`RUOYI_REDIS_HOST` 和密码。已有数据库需包含项目表结构；空库请按[源码指南](docs/development.zh-CN.md#初始化数据库)初始化，后端不会自动修改外部数据库。
- **启用配套容器**：将下面三项写入 `.env`，其余数据库账号、密码沿用示例即可用于本机体验：

```dotenv
COMPOSE_PROFILES=mysql,redis
MYSQL_HOST=mysql
RUOYI_REDIS_HOST=redis:6379
```

也可以只设置 `COMPOSE_PROFILES=mysql` 或 `COMPOSE_PROFILES=redis`，另一项继续连接外部服务。内置 MySQL 端口为 `3306`，数据库普通账号不能使用 `root`。容器访问宿主机请使用 `host.docker.internal`；宿主机服务必须允许容器网络连接。

其他配置：

| 配置项 | 用途 |
| --- | --- |
| `RUOYI_HTTP_PORT` | 后端端口，默认 `8888`；被占用时可改为 `18888` |
| `RUOYI_BIND_ADDRESS` | 默认 `127.0.0.1`，仅允许本机访问 |
| `MYSQL_PASSWORD` / `MYSQL_ROOT_PASSWORD` | 连接数据库的密码 / 仅内置 MySQL 使用的管理员密码 |
| `RUOYI_JWT_SECRET` | token 签名密钥，至少 32 个字符；重启时保持不变 |
| `RUOYI_REDIS_PASSWORD` | Redis 密码；外部服务需填写实际密码，无密码时留空 |

多人使用或部署到服务器前，请将示例密码和密钥改成随机字母数字串。例如，`openssl rand -hex 32` 可生成一个合适的随机值。不要公开 `.env`，Git 已忽略该文件。

### 4. 启动

```bash
docker compose up -d --build --wait
```

第一次启动会下载镜像并编译 Go，可能需要几分钟。只有启用内置 MySQL 时才会初始化空数据卷中的数据库（默认 `ry-vue`）；配套 MySQL、Redis 不占用宿主机的 3306、6379 端口。

### 5. 确认安装成功

```bash
docker compose ps
```

默认只看到 `api`；启用配套容器后还会看到 `mysql`、`redis`。已启用的服务应显示 **healthy**。后端健康检查只确认 HTTP 服务可访问，验证码和登录接口还需数据库、Redis 连接正常。在浏览器中打开：

- API 欢迎页：<http://localhost:8888/>
- Swagger 接口文档：<http://localhost:8888/swagger-ui/index.html>
- 验证码接口：<http://localhost:8888/captchaImage>

修改过端口时，请替换地址中的 `8888`。验证码接口返回带 Base64 图片的 JSON，看到 JSON 是正常现象。

开发账号为 **`admin / admin123`**、**`ry / admin123`**。它们是若依登录账号，与 `.env` 中的数据库密码不同。开放给其他人访问前，请修改这些账号的密码。

需要管理界面时，单独运行 RuoYi-Vue3 前端，将其 `/dev-api/*` 代理到 `http://127.0.0.1:8888`，并去掉 `/dev-api` 前缀。后端本身不提供登录网页。前端若运行在另一个容器内，应使用可访问的后端地址，不能指向它自己的 `localhost`。

## 常用操作

```bash
docker compose logs -f --tail=100 api    # 查看后端日志，Ctrl+C 只退出查看
docker compose logs --tail=100 mysql   # 查看内置数据库日志（启用后）
docker compose restart api            # 重启后端
docker compose down                   # 停止服务，保留数据

git pull
docker compose up -d --build --wait    # 更新代码后重新构建并启动
```

启用内置服务时，数据库和 Redis 数据保存在 Docker 数据卷中；上传和生成的文件始终保存在后端数据卷。外部数据库、Redis 数据由你原有的服务管理。**不要执行 `docker compose down -v`，除非确定要删除本次部署的全部数据。** 升级前请备份数据库和文件。初始化 SQL 只会在 MySQL 数据卷为空时执行，后续表结构升级应按版本说明操作。

从内置容器切换为外部服务时，先执行 `docker compose down`（不要加 `-v`），再清空 `COMPOSE_PROFILES`、修改连接地址并重新启动。仅清空 profile 不会自动停止已运行的数据库容器。

## 常见问题

| 现象 | 处理方法 |
| --- | --- |
| 无法连接 Docker | 启动 Docker Desktop 或 Docker 服务后重试 |
| 提示端口被占用 | 修改 `.env` 中的 `RUOYI_HTTP_PORT`，重新执行启动命令 |
| 下载镜像或 Go 依赖失败 | 检查网络和 Docker 的代理配置；仅 Go 依赖下载失败时，可在 `.env` 添加 `GOPROXY=https://goproxy.cn,direct` 后重新构建 |
| MySQL 显示 unhealthy | 执行 `docker compose logs --tail=100 mysql`；首次初始化可能需要几分钟 |
| 修改 `.env` 后数据库提示密码错误 | 已存在的数据库仍使用旧密码；恢复原配置，或显式修改数据库账号密码 |
| 能打开 API，却没有管理界面 | 本仓库只部署后端，需要另外接入前端 |
| 其他电脑无法访问 | 服务器部署时将 `RUOYI_BIND_ADDRESS` 改为 `0.0.0.0`，配置防火墙和 HTTPS 反向代理；先替换示例密码 |

启用配套容器后，MySQL 和 Redis 仍不开放宿主机端口。查看内置数据库可执行 `docker compose exec mysql mysql -u ruoyi -p ry-vue`。服务启动顺序的原理见 [Compose 官方说明](https://docs.docker.com/compose/how-tos/startup-order/)。

## 开发与兼容说明

不使用 Docker 的安装方式、SQL 初始化顺序和开发命令，见[源码启动指南](docs/development.zh-CN.md)。

- REST 契约位于 `app/admin/api/desc/`，Handler、Logic 按业务模块分组。
- Redis 管理会话有效期和续期，并按账号控制注册、新增和导入的并发写入，保持原数据库结构。
- XLSX 逐行导入、分批导出；导入每条用户记录独立提交，失败时保留成功行，与 Java 一致。支持旧版 XLS 导入。
- 定时任务通过 `internal/tasks` 显式注册 Go 函数，内置 `ryTask` 示例；自定义 Java 任务需移植为 Go 实现。
- 代码生成使用若依 Java、MyBatis、SQL、Vue 模板。

欢迎提交 Issue 和 PR。请说明复现步骤和预期的若依行为，补充相关测试，并同步更新中英文说明。反馈中不要包含密码或私有数据。

## 致谢与许可证

基于 [yangzongzhuan](https://github.com/yangzongzhuan) 的原版项目 [RuoYi-Vue](https://github.com/yangzongzhuan/RuoYi-Vue)。感谢 [go-zero](https://github.com/zeromicro/go-zero)、[go-zero-looklook](https://github.com/Mikaelemmmm/go-zero-looklook) 和 [Excelize](https://github.com/qax-os/excelize)。

采用 [MIT 许可证](LICENSE)。代码生成模板保留[上游许可证声明](app/admin/api/internal/generator/LICENSE.ruoyi)。
