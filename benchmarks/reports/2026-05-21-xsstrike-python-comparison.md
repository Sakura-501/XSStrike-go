# XSStrike-go vs Original XSStrike Comparison (2026-05-21)

## Compared Optimization

This report compares the original Python XSStrike bruteforce flow with the rewritten XSStrike-go bruteforce flow on the same local lab target and same payload corpus.

### Original XSStrike Behavior

- Code path: `XSStrike/modes/bruteforcer.py`
- Execution model: sequential payload requests per parameter
- Result channel: console log
- CLI `--threads`: used by crawl flow, not by bruteforce

### XSStrike-go Optimization

- Code path: `XSStrike-go/internal/bruteforce`
- Execution model: ordered worker pool controlled by `--threads`
- Result channel: structured JSON report
- Parity improvement: records every reflected payload per parameter, matching the original Python behavior instead of stopping at the first hit

## Reproduction

```bash
XSSTRIKE_GO_THREADS=10 bash benchmarks/scripts/compare_with_python.sh 20260521-v131-python-compare-clean2
```

The script starts `cmd/xsstrike-lab` on `127.0.0.1`, builds the Go binaries once, prepares a local Python virtual environment for the original project, and runs both tools against the deterministic reflection endpoint.

Original XSStrike reads payload files as strict UTF-8. The comparison therefore creates a UTF-8-clean copy of the public corpus and uses that same cleaned file for both tools.

## Metrics

| Metric | Original XSStrike Python v3.1.5 | XSStrike-go v1.3.1 |
| --- | ---: | ---: |
| Bruteforce threads | 1 | 10 |
| Payload tests | 2185 | 2185 |
| Reflected hits | 2185 expected on reflect lab | 2185 |
| Duration | 1241 ms | 500 ms |
| Speedup | 1.00x | 2.48x |
| Output format | console log | JSON report |

## Go Regression Tests

```bash
go test ./...
go test -race ./internal/bruteforce ./internal/requester
```

The added bruteforce test verifies that worker-pool execution runs requests in parallel and still returns deterministic payload order in the report.

## Safety Boundary

- All benchmark traffic is local to `127.0.0.1`.
- The benchmark does not scan external or production targets.
- No destructive write, persistence, credential collection, or hidden network behavior is introduced.
