#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
RESULT_ROOT="$ROOT_DIR/benchmarks/results"
RUN_ID="${1:-$(date +%Y%m%d-%H%M%S)}"
RUN_DIR="$RESULT_ROOT/$RUN_ID"
LAB_ADDR="127.0.0.1:18080"
LAB_BASE="http://$LAB_ADDR"

mkdir -p "$RUN_DIR"

if [[ ! -f "$ROOT_DIR/benchmarks/corpus/public/merged.txt" ]]; then
  bash "$ROOT_DIR/benchmarks/scripts/fetch_public_corpus.sh"
fi

go run ./cmd/xsstrike-lab --addr "$LAB_ADDR" >"$RUN_DIR/lab.log" 2>&1 &
LAB_PID=$!
cleanup() {
  kill "$LAB_PID" >/dev/null 2>&1 || true
}
trap cleanup EXIT

ready=0
for _ in $(seq 1 40); do
  if curl -fsS "$LAB_BASE/healthz" >/dev/null 2>&1; then
    ready=1
    break
  fi
  sleep 0.25
done
if [[ "$ready" != "1" ]]; then
  echo "Lab server not ready at $LAB_BASE" >&2
  exit 1
fi

run_scan() {
  local name="$1"
  local url="$2"
  go run ./cmd/xsstrike-go --url "$url" --output "$RUN_DIR/$name.json" >"$RUN_DIR/$name.log" 2>&1
}

run_scan "scan_html" "$LAB_BASE/reflect/html?q=seed"
run_scan "scan_attr" "$LAB_BASE/reflect/attr?q=seed"
run_scan "scan_script" "$LAB_BASE/reflect/script?q=seed"
run_scan "scan_dom" "$LAB_BASE/dom?q=seed"
run_scan "scan_waf" "$LAB_BASE/waf?q=seed"

go run ./cmd/xsstrike-go \
  --fuzzer \
  --url "$LAB_BASE/reflect/html?q=seed" \
  --file "$ROOT_DIR/benchmarks/corpus/public/merged.txt" \
  --output "$RUN_DIR/fuzzer_vuln.json" \
  >"$RUN_DIR/fuzzer_vuln.log" 2>&1

go run ./cmd/xsstrike-go \
  --fuzzer \
  --url "$LAB_BASE/reflect/sanitized?q=seed" \
  --file "$ROOT_DIR/benchmarks/corpus/public/merged.txt" \
  --output "$RUN_DIR/fuzzer_sanitized.json" \
  >"$RUN_DIR/fuzzer_sanitized.log" 2>&1

go run ./cmd/xsstrike-go \
  --url "$LAB_BASE/reflect/html?q=seed" \
  --file "$ROOT_DIR/benchmarks/corpus/public/merged.txt" \
  --output "$RUN_DIR/bruteforce_public.json" \
  >"$RUN_DIR/bruteforce_public.log" 2>&1

python3 "$ROOT_DIR/benchmarks/scripts/summarize_results.py" "$RUN_DIR"

echo "Benchmark finished: $RUN_DIR"
