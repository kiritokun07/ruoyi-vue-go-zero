#!/bin/sh
set -eu

project_dir=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
export PATH="/opt/homebrew/bin:$HOME/go/bin:$PATH"
GOCTL_PATH=${GOCTL_PATH:-$(command -v goctl)}
export GOCTL_PATH

if [ ! -x "$project_dir/.local/bin/mcp-zero" ]; then
    echo 'mcp-zero is not built. Run: sh scripts/setup-ai.sh' >&2
    exit 1
fi

cd "$project_dir"
exec "$project_dir/.local/bin/mcp-zero"
