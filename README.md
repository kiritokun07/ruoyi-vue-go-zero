# RuoYi Vue Go Zero

A Go backend for RuoYi, built with **go-zero**.

English | [简体中文](README.zh-CN.md)

Based on the API conventions and MySQL data model of [RuoYi-Vue](https://github.com/yangzongzhuan/RuoYi-Vue) by [yangzongzhuan](https://github.com/yangzongzhuan). This is an independent backend implementation; use a compatible [RuoYi-Vue3 frontend](https://github.com/yangzongzhuan/RuoYi-Vue3) for the UI.

## Features

- CAPTCHA, JWT authentication, role permissions and data scopes.
- Users, roles, departments, posts, menus, dictionaries, configuration and notices.
- Online sessions, audit logs, scheduled jobs, Redis and host monitoring.
- File uploads/downloads, Excel import/export and Swagger documentation.
- Code generation using RuoYi Java, MyBatis, SQL and Vue templates.

Runs as a single REST service with MySQL and Redis. Scheduled jobs currently support the built-in `ryTask` sample targets; additional targets require Go implementations.

## Quick start

Requires **Go 1.27.3**, MySQL and **Redis 6.2+**. Local MySQL integration was verified with 9.4.0. Core dependencies are go-zero 1.10.3 and Excelize v2; see [go.mod](go.mod).

```bash
git clone https://github.com/kiritokun07/ruoyi-vue-go-zero.git
cd ruoyi-vue-go-zero
go mod download
```

Initialize an **empty development database**, running these scripts once in order. Do not apply the full sequence to an existing database.

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

Replace the sample database credentials, then start the service:

```bash
export RUOYI_MYSQL_DSN='dev:password@tcp(127.0.0.1:3306)/ry-vue?parseTime=true&loc=Local'
export RUOYI_REDIS_HOST='127.0.0.1:6379'
export RUOYI_REDIS_PASSWORD=''
export RUOYI_JWT_SECRET="$(openssl rand -hex 32)"

go run ./app/admin/api -f app/admin/api/etc/admin-api.yaml
```

- API: `http://127.0.0.1:8888`
- Swagger: `/swagger-ui/index.html`
- Development accounts: `admin / admin123`, `ry / admin123`. Change these passwords before exposing the service.
- Frontend proxy: forward `/dev-api/*` to the API and strip `/dev-api`.
- Keep the JWT secret stable across deployments and preserve runtime files under `var/`.

See [configuration](app/admin/api/etc/admin-api.yaml) and [configuration defaults](app/admin/api/internal/config/config.go). API documentation is available in Swagger.

## Development

```text
app/admin/api/
├── desc/       # API contracts
├── etc/        # Configuration
├── internal/   # Handler, Logic, Model, ServiceContext
└── admin.go
scripts/        # Code generation
sql/            # Schema and development seeds
```

[desc/admin.api](app/admin/api/desc/admin.api) is the single generation entry point. Contracts and Handler/Logic packages are grouped by business module; Models handle persistence, and ServiceContext initializes dependencies.

```bash
go install github.com/zeromicro/go-zero/tools/goctl@v1.10.2
./scripts/generate-api.sh
go build ./...
go test ./...
go vet ./...
```

For MySQL integration tests, use an initialized development database:

```bash
RUOYI_SYSTEM_TEST_DSN="$RUOYI_MYSQL_DSN" go test ./...
```

Issues and pull requests are welcome. Preserve RuoYi API behavior, add relevant tests, and update both README translations. Repository guidance is in [AGENTS.md](AGENTS.md); optional AI submodules can be fetched with `git submodule update --init --recursive`.

## Acknowledgments & license

Thanks to [RuoYi-Vue](https://github.com/yangzongzhuan/RuoYi-Vue), [go-zero](https://github.com/zeromicro/go-zero), [go-zero-looklook](https://github.com/Mikaelemmmm/go-zero-looklook) and [Excelize](https://github.com/qax-os/excelize).

Licensed under [MIT](LICENSE). RuoYi generator templates retain their [upstream license notice](app/admin/api/internal/generator/LICENSE.ruoyi).
