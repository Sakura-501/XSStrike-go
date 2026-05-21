#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
PY_DIR="${XSSTRIKE_PYTHON_DIR:-$ROOT_DIR/../XSStrike}"
RESULT_ROOT="$ROOT_DIR/benchmarks/results"
RUN_ID="${1:-python-compare-$(date +%Y%m%d-%H%M%S)}"
RUN_DIR="$RESULT_ROOT/$RUN_ID"
LAB_ADDR="${XSSTRIKE_COMPARE_LAB_ADDR:-127.0.0.1:18081}"
LAB_BASE="http://$LAB_ADDR"
THREADS="${XSSTRIKE_GO_THREADS:-10}"
VENV_DIR="${XSSTRIKE_PYTHON_VENV:-$ROOT_DIR/benchmarks/.venv-xsstrike-python}"
PYTHON_BIN="$VENV_DIR/bin/python"
CORPUS="$ROOT_DIR/benchmarks/corpus/public/merged.txt"

mkdir -p "$RUN_DIR"

now_ms() {
  python3 -c 'import time; print(time.time_ns() // 1000000)'
}

if [[ ! -f "$CORPUS" ]]; then
  bash "$ROOT_DIR/benchmarks/scripts/fetch_public_corpus.sh"
fi

COMPARE_CORPUS="$RUN_DIR/corpus_utf8.txt"
python3 - "$CORPUS" "$COMPARE_CORPUS" <<'PY'
import sys
from pathlib import Path

src = Path(sys.argv[1])
dst = Path(sys.argv[2])
text = src.read_bytes().decode("utf-8", errors="ignore")
lines = [line.strip() for line in text.splitlines() if line.strip()]
dst.write_text("\n".join(lines) + "\n", encoding="utf-8")
PY

if [[ ! -x "$PYTHON_BIN" ]]; then
  python3 -m venv "$VENV_DIR"
  "$PYTHON_BIN" -m pip install -q -r "$PY_DIR/requirements.txt"
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

TARGET="$LAB_BASE/reflect/html?q=seed"

start="$(now_ms)"
"$GO_BIN" \
  --url "$TARGET" \
  --file "$COMPARE_CORPUS" \
  --threads "$THREADS" \
  --output "$RUN_DIR/go_bruteforce.json" \
  >"$RUN_DIR/go_bruteforce.log" 2>&1
go_end="$(now_ms)"
go_ms="$((go_end - start))"

start="$(now_ms)"
(
  cd "$PY_DIR"
  "$PYTHON_BIN" xsstrike.py \
    --url "$TARGET" \
    --file "$COMPARE_CORPUS" \
    --timeout 10 \
    --console-log-level ERROR \
    >"$RUN_DIR/python_bruteforce.log" 2>&1
)
py_end="$(now_ms)"
python_ms="$((py_end - start))"

python_version="$(grep -Eo 'v[0-9]+[.][0-9]+[.][0-9]+' "$PY_DIR/xsstrike.py" | head -1 || true)"
go_version="$("$GO_BIN" --version | awk '{print $2}')"

python3 - "$RUN_DIR" "$RUN_ID" "$TARGET" "$COMPARE_CORPUS" "$THREADS" "$go_version" "$python_version" "$go_ms" "$python_ms" "$CORPUS" <<'PY'
import json
import sys
from pathlib import Path

run_dir = Path(sys.argv[1])
corpus = Path(sys.argv[4])
go_report = json.loads((run_dir / "go_bruteforce.json").read_text(encoding="utf-8"))
payload_count = sum(1 for line in corpus.read_text(encoding="utf-8", errors="ignore").splitlines() if line.strip())
go_ms = int(sys.argv[8])
python_ms = int(sys.argv[9])
speedup = python_ms / go_ms if go_ms else 0.0
metadata = {
    "run_id": sys.argv[2],
    "target": sys.argv[3],
    "corpus": str(corpus),
    "source_corpus": sys.argv[10],
    "corpus_size": payload_count,
    "go": {
        "version": sys.argv[6],
        "threads": int(sys.argv[5]),
        "duration_ms": go_ms,
        "tested": go_report.get("tested", 0),
        "hits": len(go_report.get("hits", [])),
        "output": "go_bruteforce.json",
    },
    "python_original": {
        "version": sys.argv[7],
        "threads_for_bruteforce": 1,
        "duration_ms": python_ms,
        "tested_expected": payload_count,
        "hits_expected_on_reflect_lab": payload_count,
        "output": "python_bruteforce.log",
    },
    "speedup_go_vs_python": round(speedup, 4),
    "comparison_notes": [
        "Original XSStrike bruteforcer iterates payloads sequentially and logs hits to console.",
        "XSStrike-go bruteforce uses --threads worker-pool execution and writes machine-readable JSON.",
        "Both tools used the same local deterministic reflection endpoint and same public payload corpus.",
        "A UTF-8-clean copy of the public corpus is used because original XSStrike reads payload files as strict UTF-8.",
    ],
}
(run_dir / "comparison.json").write_text(json.dumps(metadata, indent=2, ensure_ascii=False, sort_keys=True), encoding="utf-8")

summary = [
    "# XSStrike-go vs XSStrike Python Bruteforce Comparison",
    "",
    f"- Run directory: `{run_dir}`",
    f"- Target: `{metadata['target']}`",
    f"- Corpus size: `{payload_count}` payloads",
    f"- Source corpus: `{metadata['source_corpus']}`",
    f"- Compared corpus: `{metadata['corpus']}`",
    "",
    "## What Was Optimized",
    "",
    "- Original XSStrike `modes/bruteforcer.py` sends payload checks sequentially and uses console output as the primary result channel.",
    "- XSStrike-go `internal/bruteforce` now uses `--threads` worker-pool execution, preserves payload order in the JSON report, and records every reflected payload per parameter.",
    "",
    "## Metrics",
    "",
    "| Tool | Threads | Duration | Tested | Hits | Output |",
    "| --- | ---: | ---: | ---: | ---: | --- |",
    f"| Original XSStrike Python {metadata['python_original']['version']} | 1 | {python_ms} ms | {payload_count} | {payload_count} expected on reflect lab | console log |",
    f"| XSStrike-go {metadata['go']['version']} | {metadata['go']['threads']} | {go_ms} ms | {metadata['go']['tested']} | {metadata['go']['hits']} | JSON report |",
    "",
    f"Measured speedup: `{speedup:.2f}x` faster for XSStrike-go on this local bruteforce workload.",
    "",
    "## Test Boundary",
    "",
    "- The comparison uses only `127.0.0.1` local lab traffic.",
    "- It does not scan external or production targets.",
    "- The same UTF-8-clean corpus copy is used for both tools because original XSStrike cannot load non-UTF-8 payload files.",
    "- The local lab reflects every payload by design, so the original Python hit count is expected to equal the compared corpus size after successful completion.",
    "",
]
(run_dir / "summary.md").write_text("\n".join(summary), encoding="utf-8")
print(run_dir / "summary.md")
PY

echo "Comparison finished: $RUN_DIR"
