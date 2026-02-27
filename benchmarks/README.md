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

执行后会生成：

- `benchmarks/results/<run-id>/scan_*.json`
- `benchmarks/results/<run-id>/fuzzer_*.json`
- `benchmarks/results/<run-id>/bruteforce_public.json`
- `benchmarks/results/<run-id>/summary.md`

## 单独启动靶场

```bash
go run ./cmd/xsstrike-lab --addr 127.0.0.1:18080
```
