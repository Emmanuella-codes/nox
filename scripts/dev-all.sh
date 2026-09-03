#!/bin/sh

set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)"

pids=""

start_process() {
  name="$1"
  shift
  (
    cd "$ROOT_DIR"
    "$@"
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

start_process web yarn workspace web dev
start_process backend sh ./scripts/dev-backend.sh

wait
