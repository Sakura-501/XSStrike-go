# XSStrike-go Release Policy

本项目使用语义化版本（`MAJOR.MINOR.PATCH`），并结合“重写阶段”的功能点推进节奏。

## 版本号规则

- `PATCH`：文档、测试、重构或小修复，不改变对外功能行为。
- `MINOR`：向后兼容的功能增强（新增扫描能力、检测逻辑优化、评测能力增强）。
- `MAJOR`：不兼容变更，或扫描核心架构有重大调整。

## 强制发版触发条件

满足任一条件必须发版（至少 `MINOR`）：

1. 自上个版本起，累计完成 `10` 个以上功能小点（已合并并推送）。
2. 扫描核心链路发生重大变化（如反射分析、候选生成、crawl/主动 fuzz 架构升级）。
3. 输出格式或 CLI 参数出现不兼容调整（需 `MAJOR`）。

## 标准发版流程

1. 更新 [CHANGELOG.md](/Users/w1nd/Desktop/开源计划重写 2026/XSStrike-go/CHANGELOG.md)。
2. 更新 [internal/version/version.go](/Users/w1nd/Desktop/开源计划重写 2026/XSStrike-go/internal/version/version.go)。
3. 执行 `go test ./...`，确保测试通过。
4. 提交 release commit 并推送。
5. 打版本 tag 并推送 tag（例如 `git tag v1.1.0 && git push origin v1.1.0`）。

## 分支与提交约定

- 每个功能小点单独 commit。
- 每个 commit 完成后立即 push。
- 对外发布时聚合为一个明确版本点，便于回滚和对比评测。
