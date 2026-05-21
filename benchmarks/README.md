# Benchmarks

该目录提供可重复的评测流程，用于对当前 `XSStrike-go` 版本做能力对比。

## 评测内容

- 本地漏洞靶场（反射 HTML/属性/脚本上下文、DOM sink、WAF 响应）
- 公开 payload 语料（去重后）
- 扫描、主动 fuzz、bruteforce 的指标采集与汇总

## 一键执行

```bash
bash benchmarks/scripts/evaluate_tool.sh
```

可用环境变量固定主动 fuzzer 并发数，便于做性能对比：

```bash
XSSTRIKE_GO_THREADS=1 bash benchmarks/scripts/evaluate_tool.sh threads-1
XSSTRIKE_GO_THREADS=10 bash benchmarks/scripts/evaluate_tool.sh threads-10
```

与原 Python XSStrike 的本地对比：

```bash
XSSTRIKE_GO_THREADS=10 bash benchmarks/scripts/compare_with_python.sh python-compare
```

该脚本会使用同一个本地 `xsstrike-lab` 反射端点和同一份公开 payload 语料，比较原 Python `modes/bruteforcer.py` 串行执行与 XSStrike-go `--threads` 并发执行的耗时、测试数和命中数。

执行后会生成：

- `benchmarks/results/<run-id>/scan_*.json`
- `benchmarks/results/<run-id>/fuzzer_*.json`
- `benchmarks/results/<run-id>/bruteforce_public.json`
- `benchmarks/results/<run-id>/metadata.json`
- `benchmarks/results/<run-id>/summary.md`

可追踪的基线报告可提交到 `benchmarks/reports/`（例如 `2026-02-27-baseline-v1.md`）。
本轮并发优化对比报告见 `benchmarks/reports/2026-05-21-fuzzer-threads.md`。
原 Python XSStrike 与 XSStrike-go 的 bruteforce 对比报告见 `benchmarks/reports/2026-05-21-xsstrike-python-comparison.md`。

## 单独启动靶场

```bash
go run ./cmd/xsstrike-lab --addr 127.0.0.1:18080
```
