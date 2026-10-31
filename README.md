# RuoYi Vue Go Zero

A **go-zero** backend implementing the API conventions and MySQL data model of [RuoYi-Vue](https://github.com/yangzongzhuan/RuoYi-Vue). Includes authentication, system administration, monitoring, Excel import/export and code generation.

English | [简体中文](README.zh-CN.md)

**Backend only.** By default, deployment starts only the API and connects to existing MySQL/Redis services. Bundled containers are optional. To use the administration UI, connect a compatible [RuoYi-Vue3 frontend](https://github.com/yangzongzhuan/RuoYi-Vue3).

## Install with Docker (recommended)

Go is not required. Use existing MySQL/Redis services, or enable the bundled containers in step 3.

### 1. Install the prerequisites

- Install [Git](https://git-scm.com/downloads).
- Windows/macOS: install and start [Docker Desktop](https://docs.docker.com/desktop/). On Windows, use Linux containers.
- Linux server: install [Docker Engine](https://docs.docker.com/engine/install/) and the [Compose plugin](https://docs.docker.com/compose/install/linux/).

Run `docker version` and `docker compose version` in a terminal. Both the Docker client and server should be available; use Compose v2.20 or newer.

### 2. Download the project

```bash
git clone https://github.com/kiritokun07/ruoyi-vue-go-zero.git
cd ruoyi-vue-go-zero
cp .env.example .env
```

These commands also work in PowerShell. All following commands run in this project directory. AI submodules are optional and are not needed to build or run the service.

### 3. Check the configuration

Open `.env` in a text editor. With `COMPOSE_PROFILES` empty, **no MySQL or Redis container starts**. Choose a connection mode first:

- **Existing services (default)**: set `MYSQL_HOST`, `MYSQL_PORT`, `MYSQL_DATABASE`, `MYSQL_USER`, `MYSQL_PASSWORD`, `RUOYI_REDIS_HOST` and its password. The database must contain the project schema; initialize an empty database using the [source guide](docs/development.md#initialize-the-database). The API does not automatically modify external databases.
- **Bundled containers**: set these three entries in `.env`. The remaining sample database credentials are suitable for a local trial:

```dotenv
COMPOSE_PROFILES=mysql,redis
MYSQL_HOST=mysql
RUOYI_REDIS_HOST=redis:6379
```

Use `COMPOSE_PROFILES=mysql` or `COMPOSE_PROFILES=redis` to enable only one service and connect externally to the other. Bundled MySQL uses port `3306` and requires a non-root application user. To reach host services, use `host.docker.internal`; those services must allow connections from the container network.

Other settings:

| Setting | Purpose |
| --- | --- |
| `RUOYI_HTTP_PORT` | API port, default `8888`; change to `18888` if occupied |
| `RUOYI_BIND_ADDRESS` | Default `127.0.0.1`, accessible only from this computer |
| `MYSQL_PASSWORD` / `MYSQL_ROOT_PASSWORD` | Connection password / administrator password for bundled MySQL only |
| `RUOYI_JWT_SECRET` | Token signing secret, at least 32 characters; keep stable across restarts |
| `RUOYI_REDIS_PASSWORD` | Redis password; match the external service, or leave empty if no password is set |

Before a shared/server installation, replace the sample passwords and secret with random alphanumeric values. For example, `openssl rand -hex 32` generates a suitable value. Keep `.env` private; Git ignores it.

### 4. Start

```bash
docker compose up -d --build --wait
```

The first start downloads images and compiles Go, which can take several minutes. Only bundled MySQL initializes an empty data volume automatically (database name defaults to `ry-vue`). Bundled MySQL/Redis do not occupy host ports 3306/6379.

### 5. Check the result

```bash
docker compose ps
```

Only `api` appears by default; enabled bundled services also appear as `mysql`/`redis`. Enabled services should be **healthy**. The API health check verifies HTTP availability only; CAPTCHA and login additionally require working database and Redis connections. Open:

- API welcome page: <http://localhost:8888/>
- Swagger API documentation: <http://localhost:8888/swagger-ui/index.html>
- CAPTCHA response: <http://localhost:8888/captchaImage>

Use your configured port if you changed it. `/captchaImage` returns JSON with a Base64 image; a JSON response is expected.

Development accounts are **`admin / admin123`** and **`ry / admin123`**. These are application accounts, separate from the database credentials in `.env`. Change their passwords before allowing other people to access the service.

For the UI, run the RuoYi-Vue3 frontend separately and configure its proxy to forward `/dev-api/*` to `http://127.0.0.1:8888`, removing the `/dev-api` prefix. The API itself does not serve a login page. A frontend running inside another container must use a reachable backend address rather than its own `localhost`.

## Everyday commands

```bash
docker compose logs -f --tail=100 api    # API logs; Ctrl+C only stops viewing
docker compose logs --tail=100 mysql   # Bundled database logs (when enabled)
docker compose restart api            # Restart the backend
docker compose down                   # Stop services; retain data

git pull
docker compose up -d --build --wait    # Rebuild and start after an update
```

Named volumes preserve bundled MySQL/Redis data and uploaded/generated files. External database and Redis storage remain managed by your existing services. **Do not use `docker compose down -v` unless you intend to delete all of this deployment's data.** Back up the database and files before upgrading. Initialization scripts run only on an empty MySQL volume; later schema upgrades must follow the release instructions.

When switching from bundled to external services, run `docker compose down` first (without `-v`), clear `COMPOSE_PROFILES`, update the connection settings and start again. Clearing a profile alone does not stop already-running database containers.

## Troubleshooting

| Symptom | What to do |
| --- | --- |
| Cannot connect to Docker | Start Docker Desktop / the Docker daemon, then retry |
| Port already allocated | Change `RUOYI_HTTP_PORT` in `.env`, then run the start command again |
| Image pull or Go module download fails | Check connectivity and Docker's proxy settings. For Go module downloads only, configure a reachable `GOPROXY` in `.env` (see its example) |
| MySQL is unhealthy | Run `docker compose logs --tail=100 mysql`; initial setup may take several minutes |
| MySQL access denied after changing `.env` | An existing database keeps its old password. Restore the previous `.env` value or change the database account password explicitly |
| API URL works but there is no UI | This repository deploys only the backend; connect the frontend separately |
| Another computer cannot access the API | For server deployment, set `RUOYI_BIND_ADDRESS=0.0.0.0`, configure the firewall and an HTTPS reverse proxy; replace sample credentials first |

Bundled MySQL and Redis are not published to host ports. To inspect bundled MySQL, use `docker compose exec mysql mysql -u ruoyi -p ry-vue` to inspect the database. See the [Compose startup documentation](https://docs.docker.com/compose/how-tos/startup-order/) for dependency health checks.

## Development and compatibility

For installation without Docker, SQL initialization order and development commands, see [the source setup guide](docs/development.md).

- REST contracts: `app/admin/api/desc/`; generated Handler/Logic packages are grouped by business module.
- Redis owns session expiry/renewal and coordinates same-account registration, creation and import writes. The original database schema is unchanged.
- XLSX import/export streams rows; import commits each user separately and retains successful rows if others fail, matching Java behavior. Legacy XLS import is supported.
- Tasks are explicit Go functions registered in `internal/tasks`. Built-in `ryTask` samples are included; custom Java tasks must be ported to Go.
- Generated business code uses RuoYi Java/MyBatis/SQL/Vue templates.

Issues and pull requests are welcome. Describe the reproduction steps and expected RuoYi behavior, add relevant tests, and update both README translations. Do not include credentials or private data in reports.

## Acknowledgments & license

Based on [RuoYi-Vue](https://github.com/yangzongzhuan/RuoYi-Vue) by [yangzongzhuan](https://github.com/yangzongzhuan). Thanks to [go-zero](https://github.com/zeromicro/go-zero), [go-zero-looklook](https://github.com/Mikaelemmmm/go-zero-looklook) and [Excelize](https://github.com/qax-os/excelize).

[MIT License](LICENSE). Generator templates retain their [upstream license notice](app/admin/api/internal/generator/LICENSE.ruoyi).
