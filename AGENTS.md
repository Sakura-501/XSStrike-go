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
