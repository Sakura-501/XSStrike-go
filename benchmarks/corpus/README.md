# Public XSS Corpus

该目录用于从公开仓库拉取 XSS payload 语料，用于回归测试和评测对比。

## 已接入来源

- `RenwaX23/XSS-Payloads`
- `pgaijin66/XSS-Payloads`
- `swisskyrepo/PayloadsAllTheThings`（`XSS Injection/Intruders`）

## 使用方式

在仓库根目录执行：

```bash
bash benchmarks/scripts/fetch_public_corpus.sh
```

脚本会输出：

- `benchmarks/corpus/public/renwax23.txt`
- `benchmarks/corpus/public/pgaijin66.txt`
- `benchmarks/corpus/public/payloadsallthethings.txt`
- `benchmarks/corpus/public/merged.txt`
- `benchmarks/corpus/public/stats.txt`

## 说明

- `public/_cache/` 仅用于本地缓存 clone，不纳入版本控制。
- 合并结果会做去空行、去注释、去重，保留原始 payload 内容。
