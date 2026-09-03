#!/bin/sh

set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)"
API_DIR="$ROOT_DIR/apps/api"
GO_CACHE="$ROOT_DIR/.cache/go-build"

pids=""

start_process() {
  name="$1"
  shift
  (
    cd "$API_DIR"
    GOCACHE="$GO_CACHE" "$@"
  ) &
  pid=$!
  pids="$pids $pid"
  echo "started $name (pid $pid)"
}

stop_processes() {
  for pid in $pids; do
    kill "$pid" 2>/dev/null || true
  done
}

trap 'stop_processes' INT TERM EXIT

mkdir -p "$GO_CACHE"

start_process api go run .
start_process notifications go run ./workers/notifications
start_process stories go run ./workers/stories
start_process media go run ./workers/media

wait
