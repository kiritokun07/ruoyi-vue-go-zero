# Run from source

[Back to README](../README.md) | [简体中文](development.zh-CN.md)

Requires **Go 1.27.3**, MySQL 8.4+ and Redis 6.2+. The local toolchain uses goctl 1.10.2; core dependencies are go-zero 1.10.3 and Excelize v2. See [go.mod](../go.mod).

The commands below use Bash/zsh/Git Bash. Start MySQL and Redis first, then clone the repository and enter its directory.

## Initialize the database

Run the following once against an **empty database**, in order. This installs development accounts and sample data. Do not run the full sequence against an existing RuoYi database.

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

Use a database account with access to `ry-vue` for the application; replace `dev:password` below with your actual credentials. The account also needs DDL privileges if you use the code generator's create-table function.

## Start the backend

```bash
export RUOYI_MYSQL_DSN='dev:password@tcp(127.0.0.1:3306)/ry-vue?parseTime=true&loc=Local'
export RUOYI_REDIS_HOST='127.0.0.1:6379'
export RUOYI_REDIS_PASSWORD=''
export RUOYI_JWT_SECRET='replace-with-your-own-random-secret-at-least-32-characters'
go mod download
go run ./app/admin/api -f app/admin/api/etc/admin-api.yaml
```

Generate your own secret (for example with `openssl rand -hex 32`) and retain it across restarts. The root `.env` file is read by Docker Compose; `go run` requires exported environment variables as shown above. PowerShell users can set variables with `$env:NAME='value'`.

Open <http://127.0.0.1:8888/swagger-ui/index.html>. Use `admin / admin123` or `ry / admin123` for development. The frontend is installed separately. Local configuration binds to `127.0.0.1`; container configuration binds to `0.0.0.0`.

Configuration: [admin-api.yaml](../app/admin/api/etc/admin-api.yaml), [defaults](../app/admin/api/internal/config/config.go). Preserve files under `var/` and Redis data across restarts; deleting session data requires users to log in again.

## Develop and verify

```text
app/admin/api/desc/       API contracts, entry: admin.api
app/admin/api/internal/  Handler, Logic, Model, ServiceContext
scripts/                 Code generation
deploy/                  Container configuration
sql/                     Schema and development seeds
```

```bash
go install github.com/zeromicro/go-zero/tools/goctl@v1.10.2
./scripts/generate-api.sh
go build ./...
go test ./...
go vet ./...
```

Tests that require MySQL are skipped unless `RUOYI_SYSTEM_TEST_DSN` is set. Use a **separate disposable database initialized with the scripts above**, never a production database:

```bash
RUOYI_SYSTEM_TEST_DSN='dev:password@tcp(127.0.0.1:3306)/ruoyi_test?parseTime=true&loc=Local' go test ./...
```

The example assumes you created and initialized `ruoyi_test` instead of `ry-vue`. Tests temporarily write fixture data.

Follow [AGENTS.md](../AGENTS.md). Optional AI tooling can be downloaded with `git submodule update --init --recursive`; it is not a runtime dependency.
