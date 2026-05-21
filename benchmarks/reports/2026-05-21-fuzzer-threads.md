# XSStrike-go Fuzzer Thread Benchmark (2026-05-21)

## Run Context

- Tool version: `v1.2.1`
- Local lab: `cmd/xsstrike-lab` at `http://127.0.0.1:18080`
- Corpus: `benchmarks/corpus/public/merged.txt`
- Corpus unique payloads: `2184`
- Scope: active fuzzer request execution with identical scan/crawl/WAF/DOM checks

## Reproduction Commands

```bash
XSSTRIKE_GO_THREADS=1 bash benchmarks/scripts/evaluate_tool.sh 20260521-threads-1
XSSTRIKE_GO_THREADS=10 bash benchmarks/scripts/evaluate_tool.sh 20260521-threads-10
```

## Metrics Snapshot

| Metric | 1 thread | 10 threads |
| --- | ---: | ---: |
| Elapsed wall time | 3s | 1s |
| Speedup | 1.00x | 3.00x |
| Fuzzer vulnerable endpoint | 2184/2184 (100.00%) | 2184/2184 (100.00%) |
| Fuzzer sanitized endpoint | 64/2184 (2.93%) | 64/2184 (2.93%) |
| Bruteforce tested / hits | 1 / 1 | 1 / 1 |
| Scan generated candidates | html=3072, attr=3096, script=3077 | html=3072, attr=3096, script=3077 |
| WAF detected | True, CloudFlare Web Application Firewall (CloudFlare) | True, CloudFlare Web Application Firewall (CloudFlare) |
| DOM findings | 3 | 3 |

## Result

The active fuzzer worker pool reduced local benchmark wall time from 3 seconds to 1 second on the same public payload corpus, while preserving the vulnerable-endpoint hit rate, sanitized-endpoint rate, generated candidate counts, WAF detection, and DOM findings. This confirms that the `--threads` implementation improves execution speed without changing detection output in the controlled benchmark.

## Limits

- The wall-time value uses second-level timing from `benchmarks/scripts/evaluate_tool.sh`, so very small local runs should be interpreted as coarse smoke data.
- The benchmark uses a local deterministic lab and does not prove behavior on external targets.
- Memory and CPU counters are not part of this report; the measured improvement is request-throughput wall time.
