# AGENTS.md

## 项目目标

`XSStrike-go` 是对原始 Python 项目 `s0md3v/XSStrike` 的 Go 重写，当前策略是：

- 小步快跑：一个功能点一个 commit
- 每次改动都优先补测试
- 保持 CLI、扫描链路、评测脚本都可独立演进

## 快速上手

```bash
go test ./...
go build ./...
go run ./cmd/xsstrike-go --help
```

常用开发命令：

```bash
go test ./internal/options ./internal/requester ./internal/crawl
go run ./cmd/xsstrike-go --url "https://example.com/?q=1"
go run ./cmd/xsstrike-go --crawl --url "https://example.com" --level 2
bash benchmarks/scripts/evaluate_tool.sh
```

## 目录速览

- `cmd/xsstrike-go`: 主 CLI 入口
- `cmd/xsstrike-lab`: 本地 benchmark/lab 靶场
- `internal/options`: 参数解析与校验
- `internal/requester`: HTTP 请求发送
- `internal/scan`: 单目标扫描主链路
- `internal/crawl`: 爬虫、表单提取、表单扫描
- `internal/fuzz`: fuzzer 主动发包模式
- `internal/bruteforce`: payload 文件爆破模式
- `internal/reflection`: 反射检测、过滤与候选生成
- `internal/dom`: DOM source/sink 分析
- `internal/report`: JSON 报告输出
- `internal/retirejs`: JS 组件漏洞识别
- `db/`: WAF 与 retirejs 数据库
- `benchmarks/`: 公开语料、评测脚本、基线报告

## 当前重要行为

- `--output/--output-json` 会自动创建父目录
- crawl seeds 会自动 `TrimSpace` 并去重
- POST 请求不会再修改调用方传入的 headers map
- CLI 会校验明显非法输入，例如无效编码、空 blind payload、非法 proxy

## 开发约定

1. 修改逻辑时优先补/改对应单测。
2. 变更后至少运行 `go test ./...`。
3. 涉及 CLI 行为变动时，同步更新 `README.md`。
4. 涉及发版时，同步查看 `RELEASE_POLICY.md` 与 `CHANGELOG.md`。
5. 保持提交粒度小，便于回滚与 benchmark 对比。

## 推荐阅读顺序

1. `README.md`
2. `cmd/xsstrike-go/main.go`
3. `internal/options/options.go`
4. `internal/scan/reflect.go`
5. `internal/crawl/run.go`
6. `benchmarks/README.md`

## 后续优化规划

按优先级建议继续推进：

### P0：稳定性与可维护性

1. 给 `scan/crawl/fuzz/bruteforce` 补更多异常路径测试：
   - 超时
   - 非法 URL
   - 空响应
   - 代理错误
2. 统一 CLI 错误输出与退出码，减少 `main.go` 中重复分支。
3. 梳理共享工具函数，减少 `cloneMap`、参数归一化这类重复逻辑。

### P1：扫描效果优化

1. 提升反射上下文分析质量：
   - 更细粒度区分 HTML/属性/JS/URL 上下文
   - 给 candidate 增加更稳定的置信度规则
2. 完善 path 模式与表单模式的 payload 注入覆盖。
3. 增加更多 WAF/RetireJS 数据校验测试，降低误报和回归风险。

### P2：性能与工程化

1. 真正使用 `--threads` 做 crawl/scan 并发控制。
2. 为 benchmark 增加更多固定 case，并输出更稳定的对比报告。
3. 增加 CI 中的 `go vet` / benchmark smoke test。

### P3：体验优化

1. 优化命令行帮助文本与模式说明。
2. 改进 JSON 报告字段可读性，方便后续接入别的工具链。
3. 逐步补齐与原 Python 版的行为差异说明。

## 执行策略

- 每次只做一个小点。
- 每个小点都带测试。
- 每完成一个点就单独 commit。
- README 只记录用户视角；规划和开发视角统一维护在本文件。
