# XSStrike-go

`XSStrike-go` 是对原始 Python 项目 [`s0md3v/XSStrike`](https://github.com/s0md3v/XSStrike)
的 Go 语言重写实现。

项目目标是按小步提交迁移核心能力，确保每个功能点都可独立回滚、可测试、可持续推进。

## 项目定位

- 目标：逐步复刻 XSStrike 的 XSS 扫描与分析工作流
- 方式：一个小功能点一个 commit，并立即 push
- 原则：先迁移基础能力，再迁移复杂扫描链路

## 功能作用（按模块）

下面是本项目要覆盖的核心能力及其作用：

1. 参数提取与请求预处理
- 作用：把 URL 查询参数、POST 表单、JSON body 统一解析为键值对。
- 价值：后续注入测试、payload 变换、扫描流程都依赖同一份参数输入。

2. 头部解析与可定制请求配置
- 作用：支持把字符串形式的 header 转为结构化 map，并合并默认请求头。
- 价值：便于复现真实浏览器请求，减少误判。

3. Payload 生成器（上下文向量）
- 作用：根据标签、事件、填充符等组合生成候选 XSS payload 向量。
- 价值：比固定 payload 字典更灵活，为上下文分析做准备。

4. Fuzzer 模式
- 作用：输出 fuzz payload 和生成向量，用于快速验证目标输入点行为。
- 价值：在完整扫描器完成前，可先用于回归验证和目标探测。

5. 编码链路（如 base64）
- 作用：对 payload 做可选编码，模拟部分绕过场景。
- 价值：为后续 WAF 对抗和自定义注入策略提供基础能力。

6. 后续目标（进行中）
- HTTP requester、扫描主流程、DOM 相关分析、爬虫、WAF 检测等。

## 当前实现状态

已实现：
- Go CLI 入口（banner/version/help）
- 默认配置与参数解析
- header/params 工具函数
- converter 与运行时变量存储工具
- 字符串/脚本/URL 上下文辅助函数
- HTTP requester（GET/POST、JSON/Form、proxy、timeout、delay）
- reflected 参数扫描最小闭环
- DOM source/sink 分析骨架
- payload 生成器（扩展规则集）
- `--fuzzer`（向量模式 + 文件 payload 模式）
- `--encode base64` 编码链路
- 扫描报告 JSON 输出

未实现（后续迁移）：
- 爬虫与种子递归模式
- WAF 检测与更完整上下文分析与利用链评分
- blind XSS 与多模式统一工作流增强

## 安装与运行

### 环境要求

- Go `>= 1.25`

### 克隆

```bash
git clone https://github.com/Sakura-501/XSStrike-go.git
cd XSStrike-go
```

### 直接运行

```bash
go run ./cmd/xsstrike-go --version
```

### 编译

```bash
go build -o xsstrike-go ./cmd/xsstrike-go
./xsstrike-go --version
```

## 使用教程

### 1) 查看帮助

```bash
go run ./cmd/xsstrike-go --help
```

### 2) 基础目标输入

```bash
go run ./cmd/xsstrike-go \
  --url "https://example.com/search?q=test"
```

作用：
- 显示当前运行默认参数（timeout/threads/delay）
- 尝试解析 URL 参数

### 3) 带 POST 数据

```bash
go run ./cmd/xsstrike-go \
  --url "https://example.com/login" \
  --data "username=alice&password=123"
```

作用：
- 从 `--data` 解析参数
- 为后续注入点遍历做准备

### 4) JSON 数据模式

```bash
go run ./cmd/xsstrike-go \
  --url "https://example.com/api" \
  --json \
  --data '{"name":"alice","age":10}'
```

作用：
- 把 JSON body 转为统一键值参数

### 5) 自定义 Header

```bash
go run ./cmd/xsstrike-go \
  --url "https://example.com" \
  --headers "User-Agent: XSStrike-go\nX-Test: 1"
```

作用：
- 解析并统计自定义头部数量

### 6) Fuzzer 模式

```bash
go run ./cmd/xsstrike-go --fuzzer --limit 10
```

作用：
- 输出 fuzz payload 规模
- 生成并展示上下文向量样例

### 7) 文件 payload 模式（兼容 `-f/--file`）

```bash
go run ./cmd/xsstrike-go --fuzzer --file default --limit 10
```

作用：
- 载入默认 payload 集或自定义 payload 文件
- 支持与 `--encode base64` 叠加

### 8) 扫描报告落盘（JSON）

```bash
go run ./cmd/xsstrike-go \
  --url "https://example.com/?q=1" \
  --output report.json
```

作用：
- 输出参数反射统计
- 输出 DOM source/sink 分析结果

## 参数说明（当前版本）

- `-u, --url`：目标 URL
- `--data`：POST/请求数据
- `--json`：按 JSON body 解析 `--data`
- `--path`：路径注入模式（预留）
- `--fuzzer`：进入 fuzz 输出模式
- `--encode`：payload 编码（已开始接入）
- `--headers`：自定义 header 字符串
- `--timeout`：请求超时（秒）
- `--threads`：并发线程数
- `--delay`：请求间隔（秒）
- `--limit`：fuzzer 模式输出条数
- `--proxy`：代理地址（如 `http://127.0.0.1:8080`）
- `-f, --file`：fuzzer 载入 payload 文件（`default` 使用内置集）
- `--output, --output-json`：将扫描报告写入 JSON 文件
- `-v, --version`：显示版本

## 开发与测试

```bash
go test ./...
```

建议每次改动执行：

1. `gofmt -w ./cmd ./internal`
2. `go test ./...`
3. `git commit`
4. `git push`

## 重写计划（Roadmap）

### Phase 1: Foundation

- [x] 初始化文档
- [x] CLI 入口、版本与 banner
- [x] 默认配置与参数解析

### Phase 2: Utility Migration

- [x] `extractHeaders` / `getUrl` / `getParams` 迁移
- [x] `converter`、运行时变量工具迁移
- [x] 更多字符串与上下文辅助函数迁移

### Phase 3: Payload & Fuzz

- [x] `randomUpper` + 生成器核心循环迁移
- [x] `--fuzzer` 最小闭环
- [x] 编码链路与更多 payload 规则迁移

### Phase 4: Scan Core

- [x] HTTP requester（timeout/proxy/headers）
- [x] reflected scan 最小闭环
- [x] DOM 分析骨架

### Phase 5: Crawl & Extended

- [ ] crawl 模式
- [ ] payload 文件 bruteforce
- [ ] WAF 检测与策略扩展

## 许可证

本项目沿用 GPLv3 开源协议（见 `LICENSE`）。
