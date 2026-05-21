#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
RESULT_ROOT="$ROOT_DIR/benchmarks/results"
RUN_ID="${1:-$(date +%Y%m%d-%H%M%S)}"
RUN_DIR="$RESULT_ROOT/$RUN_ID"
LAB_ADDR="127.0.0.1:18080"
LAB_BASE="http://$LAB_ADDR"
THREADS="${XSSTRIKE_GO_THREADS:-10}"

mkdir -p "$RUN_DIR"

now_ms() {
  python3 -c 'import time; print(time.time_ns() // 1000000)'
}

START_MS="$(now_ms)"

if [[ ! -f "$ROOT_DIR/benchmarks/corpus/public/merged.txt" ]]; then
  bash "$ROOT_DIR/benchmarks/scripts/fetch_public_corpus.sh"
fi

GO_BIN="$RUN_DIR/xsstrike-go"
LAB_BIN="$RUN_DIR/xsstrike-lab"
go build -o "$GO_BIN" ./cmd/xsstrike-go
go build -o "$LAB_BIN" ./cmd/xsstrike-lab

"$LAB_BIN" --addr "$LAB_ADDR" >"$RUN_DIR/lab.log" 2>&1 &
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
  local start end
  start="$(now_ms)"
  "$GO_BIN" --url "$url" --output "$RUN_DIR/$name.json" >"$RUN_DIR/$name.log" 2>&1
  end="$(now_ms)"
  echo "$((end - start))" >"$RUN_DIR/$name.duration_ms"
}

run_scan "scan_html" "$LAB_BASE/reflect/html?q=seed"
run_scan "scan_attr" "$LAB_BASE/reflect/attr?q=seed"
run_scan "scan_script" "$LAB_BASE/reflect/script?q=seed"
run_scan "scan_dom" "$LAB_BASE/dom?q=seed"
run_scan "scan_waf" "$LAB_BASE/waf?q=seed"

start="$(now_ms)"
"$GO_BIN" \
  --fuzzer \
  --url "$LAB_BASE/reflect/html?q=seed" \
  --file "$ROOT_DIR/benchmarks/corpus/public/merged.txt" \
  --threads "$THREADS" \
  --output "$RUN_DIR/fuzzer_vuln.json" \
  >"$RUN_DIR/fuzzer_vuln.log" 2>&1
end="$(now_ms)"
echo "$((end - start))" >"$RUN_DIR/fuzzer_vuln.duration_ms"

start="$(now_ms)"
"$GO_BIN" \
  --fuzzer \
  --url "$LAB_BASE/reflect/sanitized?q=seed" \
  --file "$ROOT_DIR/benchmarks/corpus/public/merged.txt" \
  --threads "$THREADS" \
  --output "$RUN_DIR/fuzzer_sanitized.json" \
  >"$RUN_DIR/fuzzer_sanitized.log" 2>&1
end="$(now_ms)"
echo "$((end - start))" >"$RUN_DIR/fuzzer_sanitized.duration_ms"

start="$(now_ms)"
"$GO_BIN" \
  --url "$LAB_BASE/reflect/html?q=seed" \
  --file "$ROOT_DIR/benchmarks/corpus/public/merged.txt" \
  --threads "$THREADS" \
  --output "$RUN_DIR/bruteforce_public.json" \
  >"$RUN_DIR/bruteforce_public.log" 2>&1
end="$(now_ms)"
echo "$((end - start))" >"$RUN_DIR/bruteforce_public.duration_ms"

END_MS="$(now_ms)"
TOOL_VERSION="$("$GO_BIN" --version | awk '{print $2}')"
python3 - "$RUN_DIR" "$RUN_ID" "$TOOL_VERSION" "$THREADS" "$LAB_BASE" "$((END_MS - START_MS))" <<'PY'
import json
import sys
from pathlib import Path

run_dir = Path(sys.argv[1])
durations = {}
for path in run_dir.glob("*.duration_ms"):
    durations[path.name.removesuffix(".duration_ms")] = int(path.read_text().strip())

metadata = {
    "run_id": sys.argv[2],
    "tool_version": sys.argv[3],
    "threads": int(sys.argv[4]),
    "lab_base": sys.argv[5],
    "elapsed_ms": int(sys.argv[6]),
    "durations_ms": durations,
}
(run_dir / "metadata.json").write_text(json.dumps(metadata, indent=2, sort_keys=True), encoding="utf-8")
PY

python3 "$ROOT_DIR/benchmarks/scripts/summarize_results.py" "$RUN_DIR"

echo "Benchmark finished: $RUN_DIR"
