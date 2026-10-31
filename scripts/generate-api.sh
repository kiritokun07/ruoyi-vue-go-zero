#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
api=app/admin/api/desc/admin.api
goctl api validate -api "$api"
goctl api go -api "$api" -dir app/admin/api --style go_zero
goctl api swagger --api "$api" --dir app/admin/api/internal/handler/console --filename swagger
