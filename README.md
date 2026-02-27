# XSStrike-go

`XSStrike-go` is a Go rewrite project of the original Python repository
[`s0md3v/XSStrike`](https://github.com/s0md3v/XSStrike).

The goal is to migrate the core workflow step by step while keeping behavior
predictable, testable, and easy to rollback through small git commits.

## Source Project Feature Overview

The original `XSStrike` focuses on XSS testing and includes:

- Reflected and DOM XSS scanning
- Multi-threaded crawling
- Context-aware payload generation
- Fuzzing mode
- Hidden parameter discovery
- WAF detection and evasion-related payload workflows
- Blind XSS support
- HTTP header customization
- Payload encoding support

## Rewrite Goals

- Keep CLI usage familiar for existing XSStrike users
- Port core logic from Python to idiomatic Go packages
- Add tests for each migrated unit before expanding scope
- Keep every small feature in an independent commit for safe rollback

## Rewrite Plan

### Phase 1: Foundation

- [x] Initialize repository documentation and migration roadmap
- [ ] Create Go module, CLI entry, version, and banner
- [ ] Add shared config package for defaults (timeout, threads, headers)

### Phase 2: Utility Migration

- [ ] Port request/header parsing utilities (`extractHeaders`, `getUrl`)
- [ ] Port parameter parsing (`getParams`) for GET/POST/JSON modes
- [ ] Port common helpers used by payload engine

### Phase 3: Payload and Fuzz Core

- [ ] Port random case transform and payload composition primitives
- [ ] Port minimal payload generator based on context inputs
- [ ] Port basic fuzz payload execution path

### Phase 4: Request and Scan Flow

- [ ] Implement HTTP requester abstraction with timeout/proxy/header support
- [ ] Implement baseline reflected scan workflow
- [ ] Add initial DOM-related analysis scaffolding

### Phase 5: Crawl and Extended Modes

- [ ] Port crawl mode and seed handling
- [ ] Add file-based payload bruteforce mode
- [ ] Add WAF detection integration hooks

### Phase 6: Hardening

- [ ] Increase test coverage for parser and payload generator behavior
- [ ] Add compatibility checks against selected XSStrike Python outputs
- [ ] Document usage and migration differences

## Commit Strategy

For this repository, each completed feature point is handled as:

1. Implement one small migration step.
2. Run verification (`go test`, targeted run checks).
3. Commit with a focused message.
4. Push immediately to GitHub.

This keeps the rewrite timeline auditable and rollback-friendly.
