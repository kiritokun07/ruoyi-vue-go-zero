#!/bin/sh
set -eu

project_dir=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$project_dir"

git submodule update --init --recursive
go version
goctl --version
mkdir -p .local/bin
cd .ai/mcp-zero
go build -o "$project_dir/.local/bin/mcp-zero" .
"$project_dir/.local/bin/mcp-zero" -version
